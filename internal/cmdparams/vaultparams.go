package cmdparams

import (
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

func GetVaultProgramParams(assets string) (*ProgramParams, error) {

	vaultdir := utils.Join(assets, "vault")

	configtemplate := utils.Join(vaultdir, "templates/vault.server.hcl")
	configoutput := utils.Join(vaultdir, "config/vault.server.hcl")

	vaultconfigparams := VaultConfigParams{
		LogLevel:                 "DEBUG",
		StorageConsulAddress:     "127.0.0.1:8500",
		StorageConsulPath:        "vault/",
		TcpAddress:               "127.0.0.1:8200",
		TcpTlsDisable:            1,
		TelemetryStatsdAddress:   "127.0.0.1:8125",
		TelemetryDisableHostname: true,
		UiEnabled:                true,
	}

	err := utils.Render(configtemplate, configoutput, vaultconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(vaultdir, "vault_"+runtime.GOARCH+".exe")

	return &ProgramParams{
		DirPath:     vaultdir,
		ExeFullname: exefile,
		AdditionalParams: []string{
			"server",
			"-config=" + configoutput,
		},
	}, nil
}
