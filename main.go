package main

import (
	"os"

	"github.com/open-infra/osc/cmd"
	"github.com/open-infra/osc/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	_ "k8s.io/client-go/plugin/pkg/client/auth"
)

func init() {
	config.EnsurePath(config.OscLogs, config.DefaultDirMod)
}

func setupLog() {
	mod := os.O_CREATE | os.O_APPEND | os.O_WRONLY
	file, err := os.OpenFile(config.OscLogs, mod, config.DefaultFileMod)
	defer func() {
		_ = file.Close()
	}()
	if err != nil {
		panic(err)
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: file})
}

func main() {
	setupLog()
	cmd.Execute()
}
