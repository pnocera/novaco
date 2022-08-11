package gitserver

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os/exec"

	"github.com/pnocera/novaco/internal/utils"
)

func receivePackHandler(w http.ResponseWriter, r *http.Request) {
	userName, repoName, _ := GetParamValues(r)
	execPath := RepoPath(userName, repoName)

	cmd := exec.Command(cfg.GitExePath(), "receive-pack", "--stateless-rpc", execPath)
	stdin, stdout, stderr, ok := GetChildPipes(cmd, w)
	if !ok {
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println(err)
		http.Error(w, "Error while spawning", http.StatusInternalServerError)
		return
	}

	reqBody, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println("Error while reading request body:", err)
		http.Error(w, "Error while reading request body", http.StatusInternalServerError)
		return
	}
	stdin.Write(reqBody)

	contentType := "application/x-git-receive-pack-result"
	SetHeader(w, contentType)

	go io.Copy(w, stdout)
	go io.Copy(w, stderr)

	if err := cmd.Wait(); err != nil {
		log.Println("Error while waiting:", err)
		return
	}

}

func receivePackHandler2(w http.ResponseWriter, req *http.Request) {
	userName, repoName, _ := GetParamValues(req)

	service := "receive-pack"

	process := cfg.ReceivePackExePath()
	cwd := utils.Join(cfg.ReposPath, userName, repoName)

	headers := w.Header()
	headers.Add("Content-Type", fmt.Sprintf("application/x-%s-result", service))
	w.WriteHeader(http.StatusOK)

	cmd := exec.Command(process, "--stateless-rpc", ".")
	cmd.Dir = cwd

	body, err := decompress(req)
	if err != nil {
		log.Printf("[ERROR] Error attempting to decompress request body: %+v", err)
		body = req.Body
	}

	runCommand(w, body, cmd)
	req.Body.Close()
}
