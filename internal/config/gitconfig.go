package config

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"

	"github.com/pnocera/novaco/internal/hclencoder"
	"github.com/pnocera/novaco/internal/utils"
)

type GitConfig struct {
	LogLevel    string       `hcl:"log_level"`     // Log level
	Hostname    string       `hcl:"hostname"`      // Hostname
	Port        int          `hcl:"port"`          // Port
	TlsCertPath string       `hcl:"tls_cert_path"` // TLS cert path
	TlsKeyPath  string       `hcl:"tls_key_path"`  // TLS key path
	KeyDir      string       `hcl:"key_dir"`       // Directory for server ssh keys. Only used in SSH strategy.
	Dir         string       `hcl:"dir"`           // Directory that contains repositories
	GitPath     string       `hcl:"git_path"`      // Path to git binary
	GitUser     string       `hcl:"git_user"`      // User for ssh connections
	AutoCreate  bool         `hcl:"auto_create"`   // Automatically create repostories
	AutoHooks   bool         `hcl:"auto_hooks"`    // Automatically setup git hooks
	Hooks       *HookScripts `hcl:"hook_scripts"`  // Scripts for hooks/* directory
	Auth        bool         `hcl:"auth"`          // Require authentication
}

type TLSConfig struct {
	TlsCertPath string `hcl:"tls_cert_path"` // TLS cert path
	TlsKeyPath  string `hcl:"tls_key_path"`  // TLS key path
}

// HookScripts represents all repository server-size git hooks
type HookScripts struct {
	PreReceive  string `hcl:"pre_receive"`
	Update      string `hcl:"update"`
	PostReceive string `hcl:"post_receive"`
}

func (c *GitConfig) ServerAddr() string {
	return fmt.Sprintf("%s:%d", c.Hostname, c.Port)
}

func (c *GitConfig) GitExePath() string {
	return utils.Join(c.GitPath, "git.exe")
}

func (c *GitConfig) UploadPackExePath() string {
	return utils.Join(c.GitPath, "git-upload-pack.exe")
}

func (c *GitConfig) ReceivePackExePath() string {
	return utils.Join(c.GitPath, "git-receive-pack.exe")
}

func (c *GitConfig) GetExePath(service string) string {
	if service == "receive-pack" {
		return c.ReceivePackExePath()
	} else if service == "upload-pack" {
		return c.UploadPackExePath()
	} else {
		return ""
	}
}

func (c *GitConfig) Merge(other *GitConfig) *GitConfig {
	if other.LogLevel != "" {
		c.LogLevel = other.LogLevel
	}
	if other.Hostname != "" {
		c.Hostname = other.Hostname
	}
	if other.Port != 0 {
		c.Port = other.Port
	}
	if other.KeyDir != "" {
		c.KeyDir = other.KeyDir
	}
	if other.Dir != "" {
		c.Dir = other.Dir
	}
	if other.GitPath != "" {
		c.GitPath = other.GitPath
	}
	if other.GitUser != "" {
		c.GitUser = other.GitUser
	}
	if other.AutoCreate {
		c.AutoCreate = other.AutoCreate
	}
	if other.AutoHooks {
		c.AutoHooks = other.AutoHooks
	}
	if other.Hooks != nil {
		c.Hooks = &HookScripts{
			PreReceive:  other.Hooks.PreReceive,
			Update:      other.Hooks.Update,
			PostReceive: other.Hooks.PostReceive,
		}
	}
	if other.Auth {
		c.Auth = other.Auth
	}

	if other.TlsCertPath != "" {
		c.TlsCertPath = other.TlsCertPath
	}
	if other.TlsKeyPath != "" {
		c.TlsKeyPath = other.TlsKeyPath
	}

	return c
}

func GetDefaultGitConfig() *GitConfig {

	return &GitConfig{
		LogLevel:   "INFO",
		Hostname:   "localhost",
		Port:       0,
		KeyDir:     "",
		Dir:        "",
		GitPath:    "",
		GitUser:    "",
		AutoCreate: false,
		AutoHooks:  false,
		Hooks:      &HookScripts{},
		Auth:       false,
	}
}

func NewGitConfig(configPath []string) *GitConfig {

	cfg := GetDefaultGitConfig()

	log.Printf("[DEBUG] loading configPath from flags length %d", len(configPath))

	for _, path := range configPath {
		current, err := LoadGitConfig(path)
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

// Configure hook scripts in the repo base directory
func (c *HookScripts) SetupInDir(path string) error {
	basePath := utils.Join(path, "hooks")
	scripts := map[string]string{
		"pre-receive":  c.PreReceive,
		"update":       c.Update,
		"post-receive": c.PostReceive,
	}

	// Cleanup any existing hooks first
	hookFiles, err := ioutil.ReadDir(basePath)
	if err == nil {
		for _, file := range hookFiles {
			if err := os.Remove(utils.Join(basePath, file.Name())); err != nil {
				return err
			}
		}
	}

	// Write new hook files
	for name, script := range scripts {
		fullPath := utils.Join(basePath, name)

		// Dont create hook if there's no script content
		if script == "" {
			continue
		}

		if err := ioutil.WriteFile(fullPath, []byte(script), 0755); err != nil {
			log.Println("hook-update", err)
			return err
		}
	}

	return nil
}

func (c *GitConfig) KeyPath() string {
	return utils.Join(c.KeyDir, "gitkit.rsa")
}

func (c *GitConfig) SetupSSL() error {
	var err error
	if c.TlsCertPath != "" && c.TlsKeyPath != "" {
		//should be ok here
	} else {
		//we need to create a self-signed cert
		key, cert, err := utils.CreateGitSelfSignedKeyCert()
		if err != nil {
			return err
		}
		c.TlsCertPath = cert
		c.TlsKeyPath = key
		//save it to custom dir
		tlsconfig := &TLSConfig{
			TlsCertPath: c.TlsCertPath,
			TlsKeyPath:  c.TlsKeyPath,
		}
		bytes, err := hclencoder.Encode(tlsconfig)
		if err != nil {
			return err
		}
		assets, err := utils.Assets()
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(utils.Join(assets, "config/git/custom/tls.hcl"), bytes, 0644)
		if err != nil {
			return err
		}
	}
	return err
}

func (c *GitConfig) Setup() error {
	if _, err := os.Stat(c.Dir); err != nil {
		if err = os.Mkdir(c.Dir, 0755); err != nil {
			return err
		}
	}

	if c.AutoHooks {
		return c.setupHooks()
	}

	return nil
}

func (c *GitConfig) setupHooks() error {
	files, err := ioutil.ReadDir(c.Dir)
	if err != nil {
		return err
	}

	for _, file := range files {
		if !file.IsDir() {
			continue
		}

		path := utils.Join(c.Dir, file.Name())

		if err := c.Hooks.SetupInDir(path); err != nil {
			return err
		}
	}

	return nil
}
