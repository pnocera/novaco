// This Source Code Form is subject to the terms of the Mozilla Public
// License, version 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

package main

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/c4milo/handlers/logger"
	"github.com/hashicorp/logutils"
	"github.com/pnocera/novaco/internal/config"
	"github.com/pnocera/novaco/internal/gitserver"
	"gopkg.in/tylerb/graceful.v1"
)

// Version is injected in build time and defined in the Makefile
var Version string = "1.0.0"

// Name is injected in build time and defined in the Makefile
var Name string = "git-backend"

func main0() {
	cfg := config.NewConfig([]string{})
	var logWriter = os.Stderr

	filter := &logutils.LevelFilter{
		Levels:   []logutils.LogLevel{"DEBUG", "WARN", "ERROR"},
		MinLevel: logutils.LogLevel(cfg.LogLevel),
		Writer:   logWriter,
	}

	log.SetOutput(filter)

	mux := http.DefaultServeMux
	rack := gitserver.Handler(mux, gitserver.NewHandler(cfg))
	rack = logger.Handler(rack, logger.AppName(Name))

	timeout, err := time.ParseDuration(cfg.ShutdownTimeout)
	if err != nil {
		log.Fatalf("[ERROR] %v", err)
	}

	log.Printf("[INFO] Listening on %s...", cfg.ServerAddr)
	log.Printf("[INFO] Serving Git repositories over HTTP from %s", cfg.ReposPath)

	gitserver.CheckGitVersion(2, 2, 1)
	graceful.Run(cfg.ServerAddr(), timeout, rack)
}
