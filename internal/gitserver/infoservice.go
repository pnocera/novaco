package gitserver

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
)

func infoserviceHandler(w http.ResponseWriter, r *http.Request) {
	userName, repoName, _ := GetParamValues(r)
	service := FindService(r)
	if ok := IsRestricted(service); ok {
		log.Println("Operation not permitted")
		http.Error(w, "Operation not permitted", http.StatusForbidden)
		return
	}
	execPath := RepoPath(userName, repoName)
	if ok := IsExistingRepository(execPath); !ok {
		log.Println("repository not found")
		http.Error(w, "repository not found", http.StatusNotFound)
		return
	}

	cmd := exec.Command(cfg.GitExePath(), service, "--stateless-rpc", "--advertise-refs", execPath)
	_, stdout, stderr, ok := GetChildPipes(cmd, w)
	if !ok {
		return
	}

	if err := cmd.Start(); err != nil {
		log.Println("Error while spawning:", err)
		http.Error(w, "Error while spawning", http.StatusInternalServerError)
		return
	}

	contentType := fmt.Sprintf("application/x-git-%s-advertisement", service)
	SetHeader(w, contentType)
	w.Write([]byte(CreateFirstPKTLine(service)))
	go io.Copy(w, stdout)
	go io.Copy(w, stderr)
	if err := cmd.Wait(); err != nil {
		log.Println("Error while waiting:", err)
		return
	}
}

func infoserviceHandler2(w http.ResponseWriter, req *http.Request) {

	userName, repoName, _ := GetParamValues(req)
	service := FindService(req)

	if ok := IsRestricted(service); ok {
		log.Println("Operation not permitted")
		http.Error(w, "Operation not permitted", http.StatusForbidden)
		return
	}
	repoPath := cfg.GetRepoPath(userName, repoName)
	if ok := IsExistingRepository(repoPath); !ok {
		log.Println("repository not found")
		http.Error(w, "repository not found", http.StatusNotFound)
		return
	}

	process := cfg.GetExePath(service)
	cwd := repoPath

	log.Printf("[DEBUG] infoHandler process: %s", process)
	log.Printf("[DEBUG] infoHandler cwd: %s", cwd)

	headers := w.Header()
	headers.Add("Content-Type", fmt.Sprintf("application/x-%s-advertisement", service))
	w.WriteHeader(http.StatusOK)

	w.Write(packetWrite(fmt.Sprintf("# service=%s\n", service)))
	w.Write(packetFlush())

	cmd := exec.Command(process, "--stateless-rpc", "--advertise-refs", ".")
	cmd.Dir = cwd

	body, err := decompress(req)
	if err != nil {
		log.Printf("[ERROR] Error attempting to decompress request body: %+v", err)
		body = req.Body
	}

	runCommand(w, body, cmd)
	req.Body.Close()
}
