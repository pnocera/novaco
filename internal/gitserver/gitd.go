package gitserver

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"net/http"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strconv"
	"strings"

	"github.com/pnocera/novaco/internal/config"
)

// http://commandcenter.blogspot.com/2014/01/self-referential-functions-and-design.html
type option func(*handler)

// Internal handler
type handler struct {
}

func NewHandler(config *config.Config) option {
	cfg = config
	return func(l *handler) {
	}
}

// Handler configures the handler and returns an HTTP handler function.
func Handler(h http.Handler, opts ...option) http.Handler {

	// Default configuration.
	handler := &handler{
		// reposPath: reposPath,
	}

	// Sets users specified configurations, overriding default ones.
	for _, opt := range opts {
		opt(handler)
	}

	handlers := map[*regexp.Regexp]func(http.ResponseWriter, *http.Request, string){
		regexp.MustCompile("(.*?)/git-upload-pack$"):  uploadPack,
		regexp.MustCompile("(.*?)/git-receive-pack$"): receivePack,
		regexp.MustCompile("(.*?)/info/refs$"):        infoRefs,
	}

	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		for re, fn := range handlers {
			if m := re.FindStringSubmatch(req.URL.Path); m != nil {
				repoPath := m[1]
				fn(w, req, repoPath)
				return
			}
		}
		h.ServeHTTP(w, req)
	})
}

// uploadPack runs git-upload-pack in a safe manner.
func uploadPack(w http.ResponseWriter, req *http.Request, repoPath string) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	service := "git-upload-pack"
	process := filepath.Join(cfg.GitBinPath, service+".exe")
	cwd := filepath.Join(cfg.ReposPath, repoPath)

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

// receivePack runs git-receive-pack in a safe manner.
func receivePack(w http.ResponseWriter, req *http.Request, repoPath string) {
	if req.Method != "POST" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}
	service := "git-receive-pack"

	process := filepath.Join(cfg.GitBinPath, service+".exe")
	cwd := filepath.Join(cfg.ReposPath, repoPath)

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

// infoRefs returns Git object refs.
func infoRefs(w http.ResponseWriter, req *http.Request, repoPath string) {
	if req.Method != "GET" {
		w.WriteHeader(http.StatusMethodNotAllowed)
		w.Write([]byte("Method Not Allowed"))
		return
	}

	service := req.URL.Query().Get("service")

	process := filepath.Join(cfg.GitBinPath, service+".exe")
	cwd := filepath.Join(cfg.ReposPath, repoPath)

	if service != "git-receive-pack" && service != "git-upload-pack" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Bad Request"))
		return
	}

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

// runCommand executes a shell command and pipes its output to HTTP response writer.
// DO NOT expose this function directly to end users as it will create a security breach.

// packetWrite returns bytes of a git packet containing the given string
func packetWrite(str string) []byte {
	s := strconv.FormatInt(int64((len(str) + 4)), 16)

	m := len(s) % 4
	if m != 0 {
		s = strings.Repeat("0", 4-m) + s
	}

	return []byte(s + str)
}

func packetFlush() []byte {
	return []byte("0000")
}

// decompress unzips request body if it is compressed by Git clients.
func decompress(r *http.Request) (io.Reader, error) {
	encoding := r.Header.Get("Content-Encoding")
	if encoding != "gzip" && encoding != "x-gzip" {
		return r.Body, nil
	}

	return gzip.NewReader(r.Body)
}

// sanitize Sanitizes name to avoid overwriting sensitive system files
// or executing forbidden binaries
func sanitize(name string) string {
	// Gets rid of volume drive label in Windows
	if len(name) > 1 && name[1] == ':' && runtime.GOOS == "windows" {
		name = name[2:]
	}

	name = filepath.Clean(name)
	name = filepath.ToSlash(name)
	for strings.HasPrefix(name, "../") {
		name = name[3:]
	}
	return name
}

// checkGitVersion checks a given Git version and returns whether or not
// the required version is installed in the system.
func CheckGitVersion(major, minor, patch int) bool {
	git := filepath.Join(cfg.GitBinPath, "git.exe")

	cmd := exec.Command(git, "--version")
	var stdout string
	var err error
	if stdout, _, err = runAndLog(cmd); err != nil {
		log.Printf("[ERROR] %v", err)
		return false
	}

	output := strings.Split(stdout, "\n")
	if len(output) < 2 {
		log.Printf("[DEBUG] git version output: %v", output)
		return false
	}

	parts := strings.Split(output[0], " ")
	if len(parts) < 3 {
		log.Printf("[DEBUG] git version parts: %v", parts)
		return false
	}

	version := strings.Split(parts[2], ".")
	major2, _ := strconv.Atoi(version[0])
	minor2, _ := strconv.Atoi(version[1])
	patch2, _ := strconv.Atoi(version[2])

	if major2 < major && minor2 < minor && patch2 < patch {
		log.Printf("[INFO] git version not supported: %d.%d.%d < %d.%d.%d", major2, minor2, patch2, major, minor, patch)
		return false
	}

	return true
}

// Borrowed from https://github.com/mitchellh/packer/blob/master/builder/vmware/common/driver.go
// runAndLog executes Git commands and logs output.
func runAndLog(cmd *exec.Cmd) (string, string, error) {
	var stdout, stderr bytes.Buffer

	log.Printf("[GitD] Executing: %s %v", cmd.Path, cmd.Args[1:])
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	err := cmd.Run()

	stdoutString := strings.TrimSpace(stdout.String())
	stderrString := strings.TrimSpace(stderr.String())

	if _, ok := err.(*exec.ExitError); ok {
		message := stderrString
		if message == "" {
			message = stdoutString
		}

		err = fmt.Errorf("[GitD] error: %s", message)
	}

	log.Printf("stdout: %s", stdoutString)
	log.Printf("stderr: %s", stderrString)

	// Replace these for Windows, we only want to deal with Unix
	// style line endings.
	returnStdout := strings.Replace(stdout.String(), "\r\n", "\n", -1)
	returnStderr := strings.Replace(stderr.String(), "\r\n", "\n", -1)

	return returnStdout, returnStderr, err
}
