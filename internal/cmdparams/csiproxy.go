package cmdparams

import (
	"path/filepath"

	"github.com/pnocera/novaco/internal/utils"
)

func GetProxyParams() (*ProgramParams, error) {

	exefile := utils.Join(utils.BinPath("csi-proxy"), "csi-proxy.exe")

	return &ProgramParams{
		ID:          "csiproxy",
		DirPath:     filepath.Dir(exefile),
		ExeFullname: exefile,
		AdditionalParams: []string{
			"-kubelet-path",
			utils.DataPath("csi-proxy"),
		},
	}, nil
}
