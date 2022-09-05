package cmdparams

import (
	"path/filepath"

	"github.com/pnocera/novaco/internal/utils"
)

func GetProxyParams(assets string, runtype string) (*ProgramParams, error) {

	exefile := utils.Join(assets, "bin/csi-proxy/csi-proxy.exe")

	return &ProgramParams{
		ID:          "gitea",
		DirPath:     filepath.Dir(exefile),
		ExeFullname: exefile,
		AdditionalParams: []string{
			"kubelet-path",
			utils.Join(assets, "data/csi-proxy"),
		},
		LogFile: utils.Join(assets, "logs/csi-proxy/csi-proxy."+runtype+".log"),
	}, nil
}
