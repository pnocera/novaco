package config

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/hashicorp/hcl"
	"github.com/pnocera/novaco/internal/utils"
)

// ParseConfigFile returns an agent.Config from parsed from a file.
func ParseGitConfigFile(path string) (*GitConfig, error) {
	// slurp
	var buf bytes.Buffer
	path, err := filepath.Abs(path)
	if err != nil {
		return nil, err
	}

	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	if _, err := io.Copy(&buf, f); err != nil {
		return nil, err
	}

	// parse
	c := &GitConfig{}

	err = hcl.Decode(c, buf.String())
	if err != nil {
		return nil, fmt.Errorf("failed to decode HCL file %s: %w", path, err)
	}

	// Set client template config or its members to nil if not set.
	finalizeGitTemplateConfig(c)

	return c, nil
}

func finalizeGitTemplateConfig(config *GitConfig) {
	// if config.GitBinPath.IsEmpty() {
	// 	config.GitBinPath = nil
	// }
}

func LoadGitConfig(path string) (*GitConfig, error) {
	fi, err := os.Stat(path)
	if err != nil {
		return nil, err
	}

	if fi.IsDir() {
		return LoadGitConfigDir(path)
	}

	cleaned := filepath.Clean(path)
	config, err := ParseGitConfigFile(cleaned)
	if err != nil {
		return nil, fmt.Errorf("Error loading %s: %s", cleaned, err)
	}

	return config, nil
}

func LoadGitConfigDir(dir string) (*GitConfig, error) {
	f, err := os.Open(dir)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	fi, err := f.Stat()
	if err != nil {
		return nil, err
	}
	if !fi.IsDir() {
		return nil, fmt.Errorf(
			"configuration path must be a directory: %s", dir)
	}

	var files []string
	err = nil
	for err != io.EOF {
		var fis []os.FileInfo
		fis, err = f.Readdir(128)
		if err != nil && err != io.EOF {
			return nil, err
		}

		for _, fi := range fis {
			// Ignore directories
			if fi.IsDir() {
				continue
			}

			// Only care about files that are valid to load.
			name := fi.Name()
			skip := true
			if strings.HasSuffix(name, ".hcl") {
				skip = false
			} else if strings.HasSuffix(name, ".json") {
				skip = false
			}
			if skip || utils.IsTemporaryFile(name) {
				continue
			}

			path := filepath.Join(dir, name)
			files = append(files, path)
		}
	}

	// Fast-path if we have no files
	if len(files) == 0 {
		return &GitConfig{}, nil
	}

	sort.Strings(files)

	var result *GitConfig
	for _, f := range files {
		config, err := ParseGitConfigFile(f)
		if err != nil {
			return nil, fmt.Errorf("Error loading %s: %s", f, err)
		}

		if result == nil {
			result = config
		} else {
			result = result.Merge(config)
		}
	}

	return result, nil
}
