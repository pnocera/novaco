package settings

import (
	"fmt"
	"net"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var once sync.Once

type settings struct {
	// Log Level
	LogLevel string
	// DataPath is the path where the data is stored
	DataPath string
	// LogPath is the path where the logs are stored
	LogPath string
	// Runtypes is a comma separated list of things to run (server and or client )
	Runtypes string
	// The leader Consul Nomad Server IP address(es)
	LeaderServerIP string
	// Bind IP address(es)
	BindIPs string
	// Vault port
	VaultPort string
	// Consul Port
	ConsulPort string
	// Nomad Port
	NomadPort string
	// Git Port
	GitPort string
	// API port
	APIPort string
	// main logger
	Logger *KLogger
}

var instance *settings

func GetSettings() *settings {
	once.Do(func() {
		instance = &settings{}
	})
	return instance
}

func (s *settings) FirstUppercaseLogLevel() string {
	return cases.Title(language.Und).String(s.LogLevel)
}

func (s *settings) UppercaseLogLevel() string {
	return strings.ToUpper(s.LogLevel)
}

func (s *settings) GetZapLevel() zapcore.Level {
	if strings.ToLower(s.LogLevel) == "debug" {
		return zap.DebugLevel
	}
	if strings.ToLower(s.LogLevel) == "info" {
		return zap.InfoLevel
	}
	if strings.ToLower(s.LogLevel) == "warn" {
		return zap.WarnLevel
	}
	if strings.ToLower(s.LogLevel) == "error" {
		return zap.ErrorLevel
	}
	if strings.ToLower(s.LogLevel) == "panic" {
		return zap.PanicLevel
	}
	if strings.ToLower(s.LogLevel) == "fatal" {
		return zap.FatalLevel
	}

	return zap.ErrorLevel
}

func (s *settings) Assets() string {

	ex, _ := os.Executable()

	assets := Join(ex, "assets")
	_, err := os.Stat(assets)

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
		return filepath.Dir(ex)
	}
	return assets

}

func Join(elem ...string) string {
	return filepath.ToSlash(filepath.Join(elem...))
}

func (s *settings) DataPathDefault() string {

	return Join(filepath.Dir(s.Assets()), "data")

}

func (s *settings) LogPathDefault() string {

	return Join(filepath.Dir(s.Assets()), "logs")

}

func (s *settings) GetGitAddress() string {
	return fmt.Sprintf("http://%s:%s", s.IP(), s.GitPort)
}

func (s *settings) GetAPIAddress() string {
	return fmt.Sprintf("http://%s:%s", s.IP(), s.APIPort)
}

func (s *settings) GetConsulAddress() string {
	return fmt.Sprintf("http://%s:%s", s.IP(), s.ConsulPort)
}

func (s *settings) IP() string {
	ip, err := getIP()
	if err != nil {
		return "127.0.0.1"
	}

	return ip.String()
}

func getIP() (net.IP, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP, nil
}
