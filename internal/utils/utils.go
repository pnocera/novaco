package utils

import (
	"os"
	"strings"
	"text/template"

	"github.com/pnocera/novaco/internal/settings"
)

var sets = settings.GetSettings()

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

func StringsContains(s []string, e string) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}
