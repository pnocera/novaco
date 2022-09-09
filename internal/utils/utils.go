package utils

import (
	"net"
	"os"
	"strings"
	"text/template"

	cmd "github.com/ShinyTrinkets/overseer"
	"github.com/go-ini/ini"
)

func IP() string {
	ip, err := GetIP()
	if err != nil {
		logger.Error("Error getting IP address", err)
		return "127.0.0.1"
	}

	return ip.String()
}

func GetIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}

func Render(input string, output string, params interface{}) error {
	tmpl, err := template.ParseFiles(input)
	if err != nil {
		return err
	}

	f, err := os.Create(output)
	if err != nil {
		return err
	}
	tmpl.Execute(f, params)

	return nil
}

func IsTemporaryFile(name string) bool {
	return strings.HasSuffix(name, "~") || // vim
		strings.HasPrefix(name, ".#") || // emacs
		(strings.HasPrefix(name, "#") && strings.HasSuffix(name, "#")) // emacs
}

func WatchStatus(statusFeed chan *cmd.ProcessJSON) {

	go func() {
		for state := range statusFeed {
			if state.ID == "gitea" && state.State == "running" {
				// do relevant git initialization here
				gitini := Join(ConfigPath("gitea"), "gitea.ini")
				cfg, err := ini.Load(gitini)
				if err != nil {
					logger.Error("Error loading git ini file: %v", err)
				} else {
					url := cfg.Section("server").Key("ROOT_URL").String()

					err = WaitForUrl(url)
					if err != nil {
						logger.Error("Error waiting for git url: %v", err)
					} else {
						err = InitGitea()
						if err != nil {
							logger.Error("Error initializing gitea: %v", err)
						}
					}
				}
			}
		}
	}()
}

func StringsContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
