package settings

import (
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
