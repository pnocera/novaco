package cmdparams

import (
	"path/filepath"
	"runtime"

	"github.com/pnocera/novaco/internal/settings"
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
		LogLevel:                 "DEBUG",
		StorageConsulAddress:     settings.GetSettings().LeaderServerIP + ":8500",
		StorageConsulPath:        "vault/",
		TcpAddress:               utils.IP() + ":8200",
		TcpTlsDisable:            1,
		TelemetryStatsdAddress:   utils.IP() + ":8125",
		TelemetryDisableHostname: true,
		UiEnabled:                true,
	}

	configtemplate := utils.TemplatePath("vault.hcl")
	configoutput := utils.Join(utils.ConfigPath("vault"), "vault.hcl")

	err := utils.Render(configtemplate, configoutput, vaultconfigparams)

	if err != nil {
		return nil, err
	}

	exefile := utils.Join(utils.BinPath("vault"), "vault_"+runtime.GOARCH+".exe")

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
