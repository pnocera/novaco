package config

type GitConfig struct {
	RunUser      string `hcl:"run_user"`      // User the service is running as
	RunMode      string `hcl:"run_mode"`      // Service run mode "dev", "prod" or "test"
	GitPath      string `hcl:"git_path"`      // Git binary full path
	DatabasePath string `hcl:"database_path"` // Path to database db
	RepoPath     string `hcl:"repo_path"`     // Path to git repositories
	LogMode      string `hcl:"log_mode"`      // Log mode "console", "file" or "syslog"
	LogLevel     string `hcl:"log_level"`     // Log level
	LogPath      string `hcl:"log_path"`      // Log path
	Domain       string `hcl:"domain"`        // Domain eg localhost
	HostIP       string `hcl:"hostip"`        // Host IP
	Port         int    `hcl:"port"`          // Port
	LfsPath      string `hcl:"lfs_path"`      // Path to git lfs
}

type TLSConfig struct {
	TlsCertPath string `hcl:"tls_cert_path"` // TLS cert path
	TlsKeyPath  string `hcl:"tls_key_path"`  // TLS key path
}

func (c *GitConfig) Merge(other *GitConfig) *GitConfig {
	if other.LogLevel != "" {
		c.LogLevel = other.LogLevel
	}

	if other.LogPath != "" {
		c.LogPath = other.LogPath
	}
	if other.RunUser != "" {
		c.RunUser = other.RunUser
	}
	if other.RunMode != "" {
		c.RunMode = other.RunMode
	}
	if other.GitPath != "" {
		c.GitPath = other.GitPath
	}
	if other.DatabasePath != "" {
		c.DatabasePath = other.DatabasePath
	}
	if other.Domain != "" {
		c.Domain = other.Domain
	}
	if other.HostIP != "" {
		c.HostIP = other.HostIP
	}
	if other.Port != 0 {
		c.Port = other.Port
	}

	return c
}
