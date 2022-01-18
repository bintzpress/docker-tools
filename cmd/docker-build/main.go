package main

import (
	"github.com/bintzpress/docker-build/internal/buildConfig"
	"github.com/bintzpress/docker-build/internal/commandExecutor"
	"github.com/bintzpress/docker-build/internal/localConfig"
	"github.com/bintzpress/docker-build/internal/setConfig"

	"errors"
	"flag"
	"fmt"
	"os"
	"strings"
)

func runDockerBuild(workspaceDir string, lc *localConfig.LocalConfig) error {
	bc, err := buildConfig.LoadConfig(workspaceDir, lc)
	if err == nil {
		err = commandExecutor.Execute(workspaceDir, bc)
	}
	return err
}

func main() {
	var dir string
	flag.StringVar(&dir, "d", ".", "Specify base directory. Default is current directory.")
	flag.Usage = func() {
		fmt.Printf("Usage:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if !strings.HasSuffix(dir, string(os.PathSeparator)) {
		dir = dir + string(os.PathSeparator)
	}
	var lc *localConfig.LocalConfig

	var err error
	lc, err = localConfig.LoadConfig(dir + ".env")
	if err != nil && errors.Is(err, os.ErrNotExist) {
		err = nil                         // ignore not found
		lc = localConfig.NewLocalConfig() // create an empty local config since no .env file
	}

	bs, err := setConfig.LoadConfig(dir)
	if err != nil && errors.Is(err, os.ErrNotExist) {
		// set file. just read single docker-builder.yml file
		err = runDockerBuild(dir, lc)
	} else if err == nil {
		var workspaceDir string
		var i int
		for i = 0; i < len(bs.Set); i++ {
			workspaceDir = dir + bs.Set[i] + string(os.PathSeparator)
			err = runDockerBuild(workspaceDir, lc)
			if err != nil {
				break
			}
		}
	}

	if err != nil {
		fmt.Println(err)
	}
}
