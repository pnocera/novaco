package utils

import (
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	cmd "github.com/ShinyTrinkets/overseer"
	"github.com/pnocera/novaco/internal/settings"
)

var sets = settings.GetSettings()

// Create self-signed certificate for Git TLS
func CreateGitSelfSignedKeyCert() (string, string, error) {

	key, cert, err := MakeCert()

	if err != nil {
		return "", "", err
	}
	outkey := Join(DataPath("git"), "localhost.key")
	outcert := Join(DataPath("git"), "localhost.crt")
	err = ioutil.WriteFile(outkey, []byte(key), 0644)
	if err != nil {
		return "", "", err
	}
	err = ioutil.WriteFile(outcert, []byte(cert), 0644)
	if err != nil {
		return "", "", err
	}

	return outkey, outcert, nil

}

func WatchStatus(statusFeed chan *cmd.ProcessJSON) {

	go func() {
		for state := range statusFeed {
			if state.ID == "gitea" && state.State == "running" {
				// do relevant git initialization here
				gitadd := sets.GetGitAddress()
				err := WaitForUrl(gitadd)

				if err != nil {
					sets.Logger.Error("Error waiting for git url: %v", err)
				} else {
					err = WaitForUrl(sets.GetConsulAddress())
					if err != nil {
						sets.Logger.Error("Error waiting for consul url: %v", err)
					} else {
						err = InitGitea()
						if err != nil {
							sets.Logger.Error("Error initializing gitea: %v", err)
						}
					}
				}
			}
		}
	}()
}

// wait for url to return a 200
func WaitForUrl(url string) error {
	for {
		resp, err := http.Get(url)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if resp.StatusCode == 200 {
			return nil
		}
	}
}

// initalize gitea
func InitGitea() error {

	sets.Logger.Info("Initializing Gitea")

	giteaexe := Join(BinPath("gitea"), "gitea.exe")

	// create gitea admin user
	gitconfig := Join(ConfigPath("gitea"), "gitea.ini")
	err1 := ExecAndWait(giteaexe,
		[]string{"admin", "user", "create", "--admin",
			"--username", "gitea_admin", "--password", "gitea_admin",
			"--email", "gitea_admin@example.com", "must-change-password",
			"false", "-c", gitconfig})

	err2 := CreateGithook()
	if err1 != nil {
		return err1
	}

	if err2 != nil {
		return err2
	}

	return nil

}

func ExecAndWait(exe string, params []string) error {

	cmd := exec.Command(exe, params...)
	cmd.Dir = filepath.Dir(exe)

	outfile, err := os.Create(Join(sets.LogPath, "init.log"))
	if err != nil {
		return err
	}
	defer outfile.Close()

	cmderr := outfile //os.Stderr

	cmdout := outfile //os.Stdout

	cmd.Stdout = cmdout
	cmd.Stderr = cmderr

	err0 := cmd.Start()
	if err0 != nil {
		return err0
	}

	err1 := cmd.Wait()

	return err1
}
