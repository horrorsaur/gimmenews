package commands

import (
	"log"
	"os"

	"github.com/hjson/hjson-go/v4"
)

type (
	Config struct {
		// The filepath of the 'nncp.hjson' file
		Path string

		// Directory that contains the encrypted data to be ingested by nncp
		SpoolPath string `json:"spool"`

		// Directory containing the nncp logs (default is /)
		LogPath string `json:"log"`

		// The 'neigh' section contains information for remote peers
		//
		// Usenet providers will be registered here
		Neighbors map[string]Neighbor `json:"neigh"`
	}

	// Describes an entry in the 'nncp' neighbor configuration section
	//
	// The first entry is 'self', as these are our public keys that can be passed off to someone
	//
	// Additional entries are human-readable names of our neighbors. For example, 'clive: { id: ... }'
	Neighbor struct {
		Id       string      `json:"id"`
		Exchpub  string      `json:"exchpub"`
		Signpub  string      `json:"signpub"`
		Noisepub string      `json:"noisepub"`
		Exec     ExecHandles `json:"exec"`
	}

	// Describes an entry in the 'exec' configuration section. Neighbor's can have multiple handles
	// each with their own command line arguments.
	//
	// For example, a neighbor that is allowed to call 'whoami --help' can be described as so:
	// exec: {
	// 	whoami: ["/usr/bin/whoami", "--help"]
	// }
	ExecHandles map[string][]string
)

// Adds an entry to the neighbors key in the nncp configuration.
//
// This updates the global configuration (default "/etc/nncp.hjson")
func (c *Config) AddNeighbor(id, exchpub, signpub, noisepub string) {}

// Saves the current config (default "/etc/nncp.hjson")
//
// This will truncate the original file and write the given data
func (c *Config) Save() (bool, error) {
	dat, err := hjson.Marshal(c)
	if err != nil {
		return false, err
	}

	if err := os.WriteFile(c.Path, dat, os.ModePerm); err != nil {
		return false, err
	}
	return true, nil
}

// The 'nncp-cfgnew' command generates a new hjson configuration
//
// Filepath is stored in $NNCPCFG (default "/etc/nncp.hjson")
func NewCfg(path string) Config {
	var c Config
	newConfig := Command{Name: "nncp-cfgnew", Path: "/usr/bin/"}
	newConfig.load()

	dat, err := newConfig.CombinedOutput()
	if err != nil {
		log.Fatal(err)
	}

	if err := hjson.Unmarshal(dat, &c); err != nil {
		log.Fatal(err)
	}

	logger.Printf(`serialized config:
logpath: %s
spoolpath: %s
neighbors: %s
	`, c.LogPath, c.SpoolPath, c.Neighbors)

	c.Path = getDefaultConfigPath()

	return c
}

// Load reads "/etc/nncp.hjson" and serializes a Config to be edited in-memory
//
// Returns a new Config if filepath does not exist
func Load(path string) (Config, error) {
	var (
		c   Config
		err error
	)

	if path == "" {
		path = "/etc/nncp.hjson"
	}

	c.Path = path

	dat, err := os.ReadFile(path)
	if err != nil {
		return c, err
	}

	if err := hjson.Unmarshal(dat, &c); err != nil {
		return c, err
	}

	logger.Printf(`serialized config:
logpath: %s
spoolpath: %s
neighbors: %s
	`, c.LogPath, c.SpoolPath, c.Neighbors)

	return c, err
}

func getDefaultConfigPath() string {
	dir := os.Getenv("NNCPCFG")
	if dir == "" {
		dir = "/etc/nncp.hjson"
	}
	return dir
}
