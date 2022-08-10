package main

import (
	"flag"

	"github.com/pnocera/novaco/internal/config"
	"github.com/pnocera/novaco/internal/gitserver"
)

func main() {
	//var additionalparams []string = []string{}

	// if os.Getenv("DEBUG") == "1" {
	// 	assets, err := utils.Assets()
	// 	if err != nil {
	// 		panic(err)
	// 	}

	// 	gitparams, err := cmdparams.GetGitParams(assets, "primary")
	// 	if err != nil {
	// 		panic(err)
	// 	}
	// 	additionalparams = gitparams.AdditionalParams
	// }

	configPath := []string{}
	flag.Var((*config.StringFlag)(&configPath), "config", "config")
	flag.Parse()

	cfg := config.NewConfig(configPath)

	gitserver.GitServer(cfg)

}
