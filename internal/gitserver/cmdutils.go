package gitserver

import (
	"io"
	"log"
	"os/exec"
)

func runCommand(w io.Writer, r io.Reader, cmd *exec.Cmd) {
	// if cmd.Dir != "" {
	// 	cmd.Dir = sanitize(cmd.Dir)
	// }

	log.Printf("[DEBUG] Running command from %s: %s %s ", cmd.Dir, cmd.Path, cmd.Args)

	stdin, err := cmd.StdinPipe()
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	stdout, err := cmd.StdoutPipe()
	if err != nil {
		log.Printf("[ERROR] %v", err)
	}

	if err := cmd.Start(); err != nil {
		log.Printf("[ERROR] %v", err)
	}

	io.Copy(stdin, r)
	io.Copy(w, stdout)
	cmd.Wait()
}
