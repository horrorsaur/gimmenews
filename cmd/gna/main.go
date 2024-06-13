package main

import (
	"log"
	"os"
	"time"

	"github.com/horrorsaur/gimmenews/internal/nntp"
)

// gimmenews articles

func main() {
	c := nntp.NewClient("news.eternal-september.org", "119", log.New(os.Stdout, "[CLIENT] ", 2))
	defer c.Close()

	for {
		select {
		case <-time.After(5 * time.Second):
			c.GetCapabilities()
		}

		break
	}
}
