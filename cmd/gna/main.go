package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/horrorsaur/gimmenews/internal/nntp"
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
	signal.Notify(sigs, os.Interrupt, syscall.SIGKILL, syscall.SIGINT)

	flag.Parse()

	if host == "" && !DEBUG {
		fmt.Println("The host config option is required.")
		flag.Usage()
		os.Exit(1)
	}

	f := newLogFile("client.log")
	defer f.Close()

	// WIP remove
	if DEBUG {
		host = "127.0.0.1"
		f = os.Stdout
	}

	clogger := log.New(f, "[CLIENT] ", 2)
	c := nntp.NewClient(host, insecure, clogger)
	defer c.Close()

	c.SendCapabilities("")

	// if v := c.SupportsAuth(); v {
	// 	c.SendAuthInfo(os.Getenv("GNA_USERNAME"), os.Getenv("GNA_PASSWORD"))
	// }

	<-sigs
	log.Printf("Received quit! Calling disconnect on client...")
	c.Quit()
}

func newLogFile(fileName string) *os.File {
	logFilePath := filepath.Join(defaultLogPath, fileName)
	file, err := os.Create(logFilePath)
	if err != nil {
		log.Printf("couldnt create log file: %s", err)
	}
	log.Printf("created log file '%s' at '%s'", fileName, defaultLogPath)
	return file
}
