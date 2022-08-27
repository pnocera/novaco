package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/hashicorp/consul/api"
	"github.com/pnocera/novaco/internal/cmdparams"
	"github.com/pnocera/novaco/internal/gitkit"
)

func main() {
	// Configure git hooks
	// hooks := &gitkit.HookScripts{
	// 	PreReceive: `echo "Hello World!"`,
	// }

	// configPath := []string{}
	// flag.Var((*config.StringFlag)(&configPath), "config", "config")
	// flag.Parse()

	cfg, err := cmdparams.GetGitConfig("primary")
	if err != nil {
		log.Fatal(err)
	}

	err = cfg.SetupSSL()
	if err != nil {
		log.Fatal(err)
	}

	//setup consul
	client, _ := api.NewClient(api.DefaultConfig())

	schema := "http"
	if cfg.TlsCertPath != "" {
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
	if err := http.ListenAndServeTLS(
		fmt.Sprintf(":%d", cfg.Port), cfg.TlsCertPath, cfg.TlsKeyPath, nil); err != nil {
		log.Fatal(err)
	}
}
