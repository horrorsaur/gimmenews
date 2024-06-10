package main

import (
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/horrorsaur/gimmenews/internal/nncp/commands"
)

var (
	userHome string

	defaultCfgPath string = "/etc/nnncp.hjson"
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
	if _, err := os.Stat(defaultCfgPath); errors.Is(err, os.ErrNotExist) {
		log.Printf("calling nncp-cfgnew to generate initial configuration")

		config := commands.NewCfg()
		fmt.Printf("config: %v", config)
	}

	fmt.Println("found conf file, skiping generate")
}
