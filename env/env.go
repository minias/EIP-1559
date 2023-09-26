package env

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v2"
)

// Glocal veriable in YAML contents
var Conf = Env{}

// Env struct
type Env struct {
	Debug      bool       `yaml:"debug"`
	BlockChain BlockChain `yaml:"blockChain"`
}

// BlockChain struct
type BlockChain struct {
	RPCURL           string `yaml:"RPCURL"`   // rpc
	SCANCURL         string `yaml:"SCANCURL"` // scanurl
	CONTRACT_ADDRESS string `yaml:"CONTRACT_ADDRESS"`
	MINIAS_ADDRESS   string `yaml:"MINIAS_ADDRESS"` // It's Me address
	TRANS_TOPICS     string `yaml:"TRANS_TOPICS"`
	BALACE_TOPICS    string `yaml:"BALACE_TOPICS"`
}

// InitProfile godoc
func InitProfile() string {
	var profile string

	//not yes yml load
	fmt.Printf("GO_PROFILE: %s\n", os.Getenv("GO_PROFILE"))

	// Impotant!!!
	profile = os.Getenv("GO_PROFILE")
	if len(profile) <= 0 {
		profile = "prod"
	}
	// Check Directory
	return "./env/" + profile + ".yml"
}

// ReadConfig godoc
func ReadConfig(profile string) {
	// YAML read file
	data, err := os.ReadFile(profile)
	if err != nil {
		fmt.Printf("Failed to read YAML file: %v\n", err)
		return
	}
	// YAML parsing
	err = yaml.Unmarshal(data, &Conf)
	if err != nil {
		fmt.Printf("Failed to parse YAML: %v\n", err)
		return
	}
}
