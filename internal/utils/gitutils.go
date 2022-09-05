package utils

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

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

//initalize gitea
func InitGitea(runtype string) error {

	log.Println("Initializing Gitea")

	assets, err := Assets()
	if err != nil {
		return err
	}

	giteaexe := Join(assets, "bin/git/gitea.exe")

	err1 := ExecAndWait(giteaexe, []string{"admin", "user", "create", "--admin", "--username", "gitea_admin", "--password", "gitea_admin", "--email", "gitea_admin@example.com", "must-change-password", "false", "-c", Join(assets, "config/git/app.ini")})

	if err1 != nil {
		return err1
	}

	// if err != nil {
	// 	log.Printf("Could not get assets : %v", err)
	// 	return err
	// }
	// dbpath := Join(assets, "data/git/git."+runtype+".db")
	// _, err = os.Stat(dbpath)
	// if os.IsNotExist(err) {
	// 	log.Println("Creating Gitea database")

	// 	giteaexe := Join(assets, "bin/git/gitea.exe")
	// 	err0 := ExecAndWait(giteaexe, []string{"migrate", "-c", Join(assets, "config/git/app.ini")})

	// 	if err0 != nil {
	// 		log.Printf("Could not create Gitea database : %v", err0)
	// 		return err0
	// 	}
	// 	err1 := ExecAndWait(giteaexe, []string{"admin", "user", "create", "--admin", "--username", "gitea_admin", "--password", "gitea_admin", "--email", "gitea_admin@example.com", "must-change-password", "false", "-c", Join(assets, "config/git/app.ini")})

	// 	if err1 != nil {
	// 		log.Printf("Could not create Gitea admin user : %v", err1)
	// 		return err1
	// 	}

	// }
	return nil

}

func ExecAndWait(exe string, params []string) error {

	cmd := exec.Command(exe, params...)
	cmd.Dir = filepath.Dir(exe)

	cmderr := os.Stderr

	cmdout := os.Stdout

	cmd.Stdout = cmdout
	cmd.Stderr = cmderr

	err0 := cmd.Start()
	if err0 != nil {
		return err0
	}

	err1 := cmd.Wait()

	return err1
}
