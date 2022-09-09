package utils

import (
	"os"
	"path/filepath"

	"github.com/pnocera/novaco/internal/settings"
)

var logger *KLogger = NewKLogger("utils")

func Join(elem ...string) string {
	return filepath.ToSlash(filepath.Join(elem...))
}

func Assets() string {
	ex, err := os.Executable()
	if err != nil {
		logger.Error("Error getting executable path: %v", err)
		return ""
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
		return filepath.Dir(ex)
	}
	return assets

}

func DataPathDefault() string {

	return Join(filepath.Dir(Assets()), "data")

}

func DataPath(app string) string {
	dataPath := settings.GetSettings().DataPath
	appdata := Join(dataPath, app)
	_, err := os.Stat(appdata)
	if os.IsNotExist(err) {
		err = os.MkdirAll(appdata, os.ModePerm)
		if err != nil {
			logger.Error("Error creating data path: %v", err)
		}
	}
	return appdata
}

func LogPathDefault() string {

	return Join(filepath.Dir(Assets()), "logs")

}

func LogPath(app string) string {
	logPath := settings.GetSettings().LogPath
	applog := Join(logPath, app)
	_, err := os.Stat(applog)
	if os.IsNotExist(err) {
		err = os.MkdirAll(applog, os.ModePerm)
		if err != nil {
			logger.Error("Error creating data path: %v", err)
		}
	}
	return applog
}

func BinPath(app string) string {
	return Join(Assets(), "bin", app)
}

func ConfigPath(app string) string {
	cfg := Join(Assets(), "config", app)
	_, err := os.Stat(cfg)
	if os.IsNotExist(err) {
		err = os.MkdirAll(cfg, os.ModePerm)
		if err != nil {
			logger.Error("Error creating config path: %v", err)
		}
	}
	return cfg
}

func TemplatePath(file string) string {
	return Join(Assets(), "templates", file)
}
