package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/consul/api"
	"github.com/pnocera/novaco/internal/config"
	"github.com/pnocera/novaco/internal/gitkit"
)

func main() {
	// Configure git hooks
	// hooks := &gitkit.HookScripts{
	// 	PreReceive: `echo "Hello World!"`,
	// }

	configPath := []string{}
	flag.Var((*config.StringFlag)(&configPath), "config", "config")
	flag.Parse()

	cfg := config.NewGitConfig(configPath)

	//setup consul
	client, _ := api.NewClient(api.DefaultConfig())

	schema := "http"
	if cfg.KeyDir != "" {
		schema = "https"
	}
	hcheck := fmt.Sprintf("%s://%s:%d/health", schema, cfg.Hostname, cfg.Port)

	client.Agent().ServiceRegister(&api.AgentServiceRegistration{
		Name: "gitserver",
		Tags: []string{"gitserver"},
		Port: cfg.Port,
		Check: &api.AgentServiceCheck{
			HTTP:     hcheck,
			Interval: "10s"},
	})

	// Configure git service
	service := gitkit.New(*cfg)

	// Configure git server. Will create git repos path if it does not exist.
	// If hooks are set, it will also update all repos with new version of hook scripts.
	if err := service.Setup(); err != nil {
		log.Fatal(err)
	}

	http.Handle("/", service)

	// Start HTTP server
	if err := http.ListenAndServe(
		fmt.Sprintf(":%d", cfg.Port), nil); err != nil {
		log.Fatal(err)
	}
}
