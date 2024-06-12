package main

import (
	"errors"
	"log"
	"os"

	"github.com/horrorsaur/gimmenews/internal/nncp/commands"
)

var (
	userHome string

	defaultCfgPath string = "/etc/nncp.hjson"
)

type (
	App struct {
		config *commands.Config
	}
)

func init() {
	h, err := os.UserHomeDir()
	if err != nil {
		log.Fatal(err)
	}
	userHome = h
	nncpCfgVar := os.Getenv("NNCPCFG")
	if nncpCfgVar != "" {
		defaultCfgPath = nncpCfgVar
	}
}

func main() {
	app := &App{}
	_, err := os.Stat(defaultCfgPath)

	var nncpCfg commands.Config
	if errors.Is(err, os.ErrNotExist) {
		log.Printf("calling nncp-cfgnew to generate initial configuration")
		nncpCfg = commands.NewCfg(defaultCfgPath)
	}

	// grab relevant nncp configuration (default "/etc/nncp.hjson")
	nncpCfg, err = commands.Load(defaultCfgPath)
	if err != nil {
		log.Fatal(err)
	}

	app.config = &nncpCfg
}
