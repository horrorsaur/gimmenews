package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/horrorsaur/gimmenews/internal/nntp"
)

// gimmenews articles

var (
	insecure *bool = flag.Bool("insecure", false, "This prompts the nntp client to connect over the traditional port, 119. This is NOT encouraged as traffic is typically unencrypted here")
)

func main() {
	flag.Parse()

	c := nntp.NewClient("news.eternal-september.org", "119", insecure, log.New(os.Stdout, "[CLIENT] ", 2))
	defer c.Close()

	for {
		select {
		case <-time.After(3 * time.Second):
			c.GetCapabilities()

			time.Sleep(2 * time.Second)

			c.List()
		}

		break
	}

	fmt.Println("Exiting...")
}
