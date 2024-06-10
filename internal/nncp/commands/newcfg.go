package commands

import (
	"encoding/json"
	"github.com/hjson/hjson-go/v4"
	"log"
)

// http://www.nncpgo.org/Configuration.html
// http://www.nncpgo.org/CfgNeigh.html

type (
	// The 'neigh' section describes remote peers & contains the definitions for our registered news servers
	Config struct {
		// The filepath of the 'nncp.hjson' file
		Path string

		SpoolPath string   `json:"spool"`
		LogPath   string   `json:"log"`
		Neighbors Neighbor `json:"neigh"`
	}

	// Describes an entry in the 'nncp' neighbor configuration section
	//
	// The first entry is 'self', as these are our public keys that can be passed off to someone
	//
	// Additional entries are human-readable names of our neighbors. For example, 'clive: { id: ... }'
	Neighbor struct {
		Id       string         `json:"id"`
		Exchpub  string         `json:"exchpub"`
		Signpub  string         `json:"signpub"`
		Noisepub string         `json:"noisepub"`
		Handles  ProgramHandles `json:"exec"`
	}

	// Describes an entry in the 'exec' configuration section. Neighbor's can have multiple handles
	// each with their own command line arguments.
	//
	// For example, a neighbor that is allowed to call 'whoami' can be described as so:
	// exec: {
	// 	whoami: ["/usr/bin/whoami", "--help"]
	// }
	ProgramHandles struct {
		Cmd map[string][]string
	}
)

// Expands the 'neigh' configuration object
func (c *Config) AddNeighborEntry() (*Neighbor, error) { return nil, nil }

// The 'nncp-cfgnew' command generates a new hjson configuration
//
// Filepath is stored in $NNCPCFG (default "/etc/nncp.hjson")
func NewCfg() Config {
	var c Config
	newConfig := Command{Name: "nncp-cfgnew", Path: "/usr/bin/"}
	newConfig.load()

	dat, err := newConfig.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	if err := json.Unmarshal(dat, &c); err != nil {
		log.Fatal(err)
	}

	return c
}
