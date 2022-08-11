package main

import (
	"log"
	"net/http"

	"github.com/pnocera/novaco/internal/gitkit"
)

func main() {
	// Configure git hooks
	hooks := &gitkit.HookScripts{
		PreReceive: `echo "Hello World!"`,
	}

	// Configure git service
	service := gitkit.New(gitkit.Config{
		Dir:        "e:/Projects/nomad/server_install/assets/data/git",
		AutoCreate: true,
		AutoHooks:  false,
		Hooks:      hooks,
		GitPath:    "e:/Projects/nomad/server_install/assets/bin/git/git.exe",
	})

	// Configure git server. Will create git repos path if it does not exist.
	// If hooks are set, it will also update all repos with new version of hook scripts.
	if err := service.Setup(); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", service)

	// Start HTTP server
	if err := http.ListenAndServe(":8888", nil); err != nil {
		log.Fatal(err)
	}
}
