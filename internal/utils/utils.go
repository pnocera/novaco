package utils

import (
	"net"
	"os"
	"path/filepath"
	"text/template"
)

func Join(elem ...string) string {
	return filepath.ToSlash(filepath.Join(elem...))
}

func CurPath(dir string) string {
	ex, err := os.Executable()
	if err != nil {
		return dir
	}
	return Join(filepath.Dir(ex), dir)
}

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
