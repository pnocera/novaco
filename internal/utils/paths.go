package utils

import (
	"os"
	"path/filepath"
)

func Join(elem ...string) string {
	return filepath.ToSlash(filepath.Join(elem...))
}

func DataPath(app string) string {
	dataPath := sets.DataPath
	appdata := Join(dataPath, app)
	_, err := os.Stat(appdata)
	if os.IsNotExist(err) {
		err = os.MkdirAll(appdata, os.ModePerm)
		if err != nil {
			sets.Logger.Error("Error creating data path: %v", err)
		}
	}
	return appdata
}

func LogPath(app string) string {
	logPath := sets.LogPath
	applog := Join(logPath, app)
	_, err := os.Stat(applog)
	if os.IsNotExist(err) {
		err = os.MkdirAll(applog, os.ModePerm)
		if err != nil {
			sets.Logger.Error("Error creating data path: %v", err)
		}
	}
	return applog
}

func BinPath(app string) string {
	return Join(sets.Assets(), "bin", app)
}

func ConfigPath(app string) string {
	cfg := Join(sets.Assets(), "config", app)
	_, err := os.Stat(cfg)
	if os.IsNotExist(err) {
		err = os.MkdirAll(cfg, os.ModePerm)
		if err != nil {
			sets.Logger.Error("Error creating config path: %v", err)
		}
	}
	return cfg
}

func TemplatePath(file string) string {
	return Join(sets.Assets(), "templates", file)
}
