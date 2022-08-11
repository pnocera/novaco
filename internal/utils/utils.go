package utils

import (
	"net"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

func Join(elem ...string) string {
	return filepath.ToSlash(filepath.Join(elem...))
}

func Assets() (string, error) {
	ex, err := os.Executable()
	if err != nil {
		return "", err
	}

	assets := Join(ex, "assets")
	_, err = os.Stat(assets)

	if os.IsNotExist(err) {
		//try to get it from the parent folder
		assets = Join(filepath.Dir(ex), "assets")
	}
	_, err = os.Stat(assets)
	if os.IsNotExist(err) {
		//try to get it from the parent folder
		assets = Join(filepath.Dir(filepath.Dir(ex)), "assets")
	}

	_, err = os.Stat(assets)
	if os.IsNotExist(err) {
		//try to get it from the parent folder
		assets = Join(filepath.Dir(filepath.Dir(filepath.Dir(ex))), "assets")
	}
	_, err = os.Stat(assets)
	if os.IsNotExist(err) {
		return "", err
	}
	return assets, nil

}

// func curPath(dir string) (string, error) {
// 	ex, err := os.Executable()
// 	if err != nil {
// 		return "", err
// 	}
// 	return Join(filepath.Dir(ex), dir), nil
// }

func GetOutboundIP() (net.IP, error) {
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
