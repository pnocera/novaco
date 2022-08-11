package config

import (
	"fmt"
	"log"

	"github.com/pnocera/novaco/internal/utils"
)

type Config struct {
	LogLevel            string `hcl:"log_level"`
	ShutdownTimeout     string `hcl:"shutdown_timeout"`
	GitBinPath          string `hcl:"git_bin_path"`
	ReposPath           string `hcl:"repos_path"`
	Hostname            string `hcl:"hostname"`
	Port                int    `hcl:"port"`
	SSLEnabled          bool   `hcl:"ssl_enabled"`
	CertFilePath        string `hcl:"cert_file_path"`
	KeyFilePath         string `hcl:"key_file_path"`
	AuthEnabled         bool   `hcl:"auth_enabled"`
	PasswdFilePath      string `hcl:"passwd_file_path"`
	RestrictReceivePack bool   `hcl:"restrict_receive_pack"`
	RestrictUploadPack  bool   `hcl:"restrict_upload_pack"`
}

func (c *Config) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}

func (c *Config) GitExePath() string {
	return utils.Join(c.GitBinPath, "git.exe")
}

func (c *Config) UploadPackExePath() string {
	return utils.Join(c.GitBinPath, "git-upload-pack.exe")
}

func (c *Config) ReceivePackExePath() string {
	return utils.Join(c.GitBinPath, "git-receive-pack.exe")
}

func (c *Config) GetRepoPath(username string, repo string) string {
	return utils.Join(c.ReposPath, username, repo)
}

func (c *Config) GetExePath(service string) string {
	if service == "receive-pack" {
		return c.ReceivePackExePath()
	} else if service == "upload-pack" {
		return c.UploadPackExePath()
	} else {
		return ""
	}
}

func (c *Config) Merge(other *Config) *Config {
	if other.LogLevel != "" {
		c.LogLevel = other.LogLevel
	}
	if other.ShutdownTimeout != "" {
		c.ShutdownTimeout = other.ShutdownTimeout
	}
	if other.GitBinPath != "" {
		c.GitBinPath = other.GitBinPath
	}
	if other.ReposPath != "" {
		c.ReposPath = other.ReposPath
	}
	if other.Hostname != "" {
		c.Hostname = other.Hostname
	}
	if other.Port != 0 {
		c.Port = other.Port
	}
	if other.SSLEnabled != false {
		c.SSLEnabled = other.SSLEnabled
	}
	if other.AuthEnabled != false {
		c.AuthEnabled = other.AuthEnabled
	}
	if other.PasswdFilePath != "" {
		c.PasswdFilePath = other.PasswdFilePath
	}
	if other.RestrictReceivePack != false {
		c.RestrictReceivePack = other.RestrictReceivePack
	}
	if other.RestrictUploadPack != false {
		c.RestrictUploadPack = other.RestrictUploadPack
	}
	if other.CertFilePath != "" {
		c.CertFilePath = other.CertFilePath
	}
	if other.KeyFilePath != "" {
		c.KeyFilePath = other.KeyFilePath
	}

	return c
}

func GetDefaultConfig() *Config {

	return &Config{
		LogLevel:            "info",
		ShutdownTimeout:     "15s",
		Hostname:            "localhost",
		Port:                8080,
		SSLEnabled:          false,
		AuthEnabled:         false,
		PasswdFilePath:      "passwd",
		RestrictReceivePack: false,
		RestrictUploadPack:  false,
		GitBinPath:          "git",
		ReposPath:           "",
		CertFilePath:        "",
		KeyFilePath:         "",
	}
}

func NewConfig(configPath []string) *Config {

	cfg := GetDefaultConfig()

	// if len(configPath) == 0 {

	//flags := flag.NewFlagSet("gitserver", flag.ContinueOnError)

	log.Printf("[INFO] loading configPath from flags length %d", len(configPath))

	// } else {
	// 	for i := range configPath {
	// 		configPath[i] = strings.Split(configPath[i], "=")[1]
	// 	}
	// }

	for _, path := range configPath {
		current, err := LoadConfig(path)
		if err != nil {
			log.Printf("[ERROR] loading config file %s: %s", path, err)
			continue
		}
		if current == nil {
			continue
		}
		cfg = cfg.Merge(current)
	}

	return cfg
}
