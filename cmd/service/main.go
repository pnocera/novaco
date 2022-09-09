package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pnocera/novaco/internal/service"
	"github.com/pnocera/novaco/internal/settings"
	"github.com/pnocera/novaco/internal/utils"
	"github.com/pnocera/novaco/internal/web"
)

func init() {
	settings := settings.GetSettings()
	regStringVar(&settings.DataPath, "data_dir", utils.DataPathDefault(), "data dir")
	regStringVar(&settings.LogPath, "logs_dir", utils.LogPathDefault(), "log dir")
	regStringVar(&settings.Runtypes, "run_types", "server,client", "things to run")
	regStringVar(&settings.LeaderServerIP, "leader_ip", utils.IP(), "leader IP")
	regStringVar(&settings.BindIPs, "bind_ips", utils.IP(), "Bind IP")
	regStringVar(&settings.LogLevel, "log_level", "DEBUG", "Log level")
	regStringVar(&settings.VaultPort, "vault_port", "8200", "Vault port")
	regStringVar(&settings.ConsulPort, "consul_port", "8500", "Consul port")
	regStringVar(&settings.NomadPort, "nomad_port", "4646", "Nomad port")
	regStringVar(&settings.GitPort, "git_port", "8888", "Git port")
	regStringVar(&settings.APIPort, "api_port", "7788", "API port")
}

func main() {

	flag.Parse()

	fmt.Println("Starting service")

	pathdir := getStringFlag("data_dir")
	_, err := os.Stat(pathdir)
	if os.IsNotExist(err) {
		err = os.MkdirAll(pathdir, os.ModePerm)
		if err != nil {
			fmt.Printf("Error creating data path: %v", err)
		}
	}

	go web.Serve()
	_ = service.StartNew()
}

func regStringVar(p *string, name string, value string, usage string) {
	if flag.Lookup(name) == nil {
		flag.StringVar(p, name, value, usage)
	}
}

func getStringFlag(name string) string {
	return flag.Lookup(name).Value.(flag.Getter).Get().(string)
}
