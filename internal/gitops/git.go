package gitops

import (
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	cmd "github.com/ShinyTrinkets/overseer"
	"github.com/pnocera/novaco/internal/settings"
	"github.com/pnocera/novaco/internal/utils"
	"github.com/sethvargo/go-password/password"
)

var sets = settings.GetSettings()

// Create self-signed certificate for Git TLS
func CreateGitSelfSignedKeyCert() (string, string, error) {

	key, cert, err := utils.MakeCert()

	if err != nil {
		return "", "", err
	}
	outkey := utils.Join(utils.DataPath("git"), "localhost.key")
	outcert := utils.Join(utils.DataPath("git"), "localhost.crt")
	err = os.WriteFile(outkey, []byte(key), 0644)
	if err != nil {
		return "", "", err
	}
	err = os.WriteFile(outcert, []byte(cert), 0644)
	if err != nil {
		return "", "", err
	}

	return outkey, outcert, nil

}

func WatchStatus(statusFeed chan *cmd.ProcessJSON) {

	go func() {
		for state := range statusFeed {
			if state.ID == "gitea" && state.State == "running" {
				// do relevant git initialization here
				gitadd := sets.GetGitAddress()
				err := WaitForUrl(gitadd)

				if err != nil {
					sets.Logger.Error("Error waiting for gitea", err)
				}

				consuladd := sets.GetConsulAddress()
				err = WaitForUrl(consuladd)

				if err != nil {
					sets.Logger.Error("Error waiting for consul url", err)
				} else {
					err = WaitForUrl(sets.GetConsulAddress())
					if err != nil {
						sets.Logger.Error("Error waiting for consul url", err)
					} else {
						err = InitGitea()
						if err != nil {
							sets.Logger.Error("Error initializing gitea", err)
						}
					}
				}
			}
		}
	}()
}

// wait for url to return a 200
func WaitForUrl(url string) error {
	for {
		resp, err := http.Get(url)
		if err != nil {
			time.Sleep(100 * time.Millisecond)
			continue
		}
		if resp.StatusCode == 200 {
			return nil
		}
	}
}

// initalize gitea
func InitGitea() error {

	sets.Logger.Info("Initializing Gitea")

	res, err := password.Generate(12, 4, 8, false, false)
	if err != nil {
		sets.Logger.Error("error generating password for gitea", err)
		return err
	}

	giteaexe := utils.Join(utils.BinPath("gitea"), "gitea.exe")

	// create gitea admin user
	gitconfig := utils.Join(utils.ConfigPath("gitea"), "gitea.ini")
	err = ExecAndWait(giteaexe,
		[]string{"admin", "user", "create", "--admin",
			"--username", "gitea_admin", "--password", res,
			"--email", "gitea_admin@example.com", "must-change-password",
			"false", "-c", gitconfig})

	if err == nil {

		err = CreateGithook(res)

	}

	return err

}

func ExecAndWait(exe string, params []string) error {

	cmd := exec.Command(exe, params...)
	cmd.Dir = filepath.Dir(exe)

	outfile, err := os.Create(utils.Join(sets.LogPath, "init.log"))
	if err != nil {
		return err
	}
	defer outfile.Close()

	cmderr := outfile //os.Stderr

	cmdout := outfile //os.Stdout

	cmd.Stdout = cmdout
	cmd.Stderr = cmderr

	err0 := cmd.Start()
	if err0 != nil {
		return err0
	}

	err1 := cmd.Wait()

	return err1
}
