package cmdparams

import (
	"path/filepath"
	"runtime"

	"github.com/pnocera/novaco/internal/utils"
)

type VaultConfigParams struct {
	LogLevel                 string
	StorageConsulAddress     string
	StorageConsulPath        string
	TcpAddress               string
	TcpTlsDisable            int
	TelemetryStatsdAddress   string
	TelemetryDisableHostname bool
	UiEnabled                bool
}

func GetVaultProgramParams(assets string, runtype string) (*ProgramParams, error) {

	// vaultdir := utils.Join(assets, "vault")

	// configtemplate := utils.Join(vaultdir, "templates/vault.server.hcl")
	// configoutput := utils.Join(vaultdir, "config/vault.server.hcl")

	ip, err := utils.GetOutboundIP()
	if err != nil {
		return nil, err
	}

	vaultconfigparams := VaultConfigParams{
		LogLevel:                 "DEBUG",
		StorageConsulAddress:     ip.String() + ":8500",
		StorageConsulPath:        "vault/",
		TcpAddress:               ip.String() + ":8200",
		TcpTlsDisable:            1,
		TelemetryStatsdAddress:   ip.String() + ":8125",
		TelemetryDisableHostname: true,
		UiEnabled:                true,
	}

	configtemplate := utils.Join(assets, "templates/vault."+runtype+".hcl")
	configoutput := utils.Join(assets, "config/vault/auto/vault."+runtype+".hcl")

	err = utils.Render(configtemplate, configoutput, vaultconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(assets, "bin/vault/vault_"+runtime.GOARCH+".exe")

	return &ProgramParams{
		DirPath:     filepath.Dir(exefile),
		ExeFullname: exefile,
		AdditionalParams: []string{
			"server",
			"-config=" + configoutput,
			"-config=" + utils.Join(assets, "config/vault/custom"),
		},
		LogFile: utils.Join(assets, "logs/vault/vault."+runtype+".log"),
	}, nil
}
