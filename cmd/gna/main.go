package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"

	"github.com/horrorsaur/gimmenews/internal/nntp"
	"github.com/horrorsaur/gimmenews/internal/utils"
	"github.com/joho/godotenv"
)

// gimmenews articles (gna)

var (
	defaultLogPath string
	host           string
	insecure       bool

	DEBUG bool = false
)

func init() {
	osCacheDir, err := os.UserCacheDir()
	if err != nil {
		log.Fatal(err)
	}

	flag.BoolVar(
		&DEBUG,
		"debug",
		false,
		"Tells the client to run in DEBUG mode",
	)

	flag.BoolVar(
		&insecure,
		"insecure",
		false,
		"This prompts the nntp client to connect over the traditional port, 119. This is NOT encouraged as traffic is typically unencrypted here",
	)

	flag.StringVar(
		&defaultLogPath,
		"log",
		filepath.Join(osCacheDir, "gimmenews", "logs"),
		"Passing this flag prompts the client to save logs to the passed in directory",
	)

	flag.StringVar(
		&host,
		"host",
		"",
		"Tells the nntp client to connect to this host (REQUIRED)",
	)
}

func main() {
	godotenv.Load()

	sigs := make(chan os.Signal)
	signal.Notify(sigs, os.Interrupt, os.Kill)

	flag.Parse()

	if host == "" {
		fmt.Println("The host config option is required.")
		flag.Usage()
	}

	f := utils.NewLogFile("client.log", defaultLogPath)
	defer f.Close()

	c := nntp.NewClient(host, insecure)
	defer c.Close()

	c.SendCapabilities("")

	// if v := c.SupportsAuth(); v {
	// 	c.SendAuthInfo(os.Getenv("GNA_USERNAME"), os.Getenv("GNA_PASSWORD"))
	// }

	<-sigs
	log.Printf("Received quit! Calling disconnect on client...")

	c.Quit()
}
