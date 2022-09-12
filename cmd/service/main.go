package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/pnocera/novaco/internal/service"
	"github.com/pnocera/novaco/internal/settings"
	"github.com/pnocera/novaco/internal/web"
)

var sets = settings.GetSettings()

func init() {

	regStringVar(&sets.DataPath, "data_dir", sets.DataPathDefault(), "data dir")
	regStringVar(&sets.LogPath, "logs_dir", sets.LogPathDefault(), "log dir")
	regStringVar(&sets.Runtypes, "run_types", "server,client", "things to run")
	regStringVar(&sets.LeaderServerIP, "leader_ip", sets.IP(), "leader IP")
	regStringVar(&sets.BindIPs, "bind_ips", sets.IP(), "Bind IP")
	regStringVar(&sets.LogLevel, "log_level", "DEBUG", "Log level")
	regStringVar(&sets.VaultPort, "vault_port", "8200", "Vault port")
	regStringVar(&sets.ConsulPort, "consul_port", "8500", "Consul port")
	regStringVar(&sets.NomadPort, "nomad_port", "4646", "Nomad port")
	regStringVar(&sets.GitPort, "git_port", "8888", "Git port")
	regStringVar(&sets.APIPort, "api_port", "7788", "API port")
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

	sets.Logger = settings.NewKLogger("main")

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
