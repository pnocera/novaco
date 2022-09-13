package cmdparams

import (
	"path/filepath"

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

func GetVaultProgramParams() (*ProgramParams, error) {

	vaultconfigparams := VaultConfigParams{
		LogLevel:                 sets.UppercaseLogLevel(),
		StorageConsulAddress:     sets.LeaderServerIP + ":" + sets.ConsulPort,
		StorageConsulPath:        "vault/",
		TcpAddress:               sets.IP() + ":" + sets.VaultPort,
		TcpTlsDisable:            1,
		TelemetryStatsdAddress:   sets.IP() + ":8125",
		TelemetryDisableHostname: true,
		UiEnabled:                true,
	}

	configtemplate := utils.TemplatePath("vault.hcl")
	configoutput := utils.Join(utils.ConfigPath("vault"), "vault.hcl")

	err := utils.Render(configtemplate, configoutput, vaultconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(utils.BinPath("vault"), "vault.exe")

	additionalparams := []string{
		"server",
		"-config=" + utils.ConfigPath("vault"),
	}

	return &ProgramParams{
		ID:               "vault",
		DirPath:          filepath.Dir(exefile),
		ExeFullname:      exefile,
		AdditionalParams: additionalparams,
	}, nil
}
