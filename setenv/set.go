package setenv

import (
	"flag"
	"fmt"
	"os"
	"log"
	"github.com/go-yaml/yaml"
)

func SetupConfigToENV () {
	// Generate our config based on the config supplied
	// by the user in the flags

	configPath, err := ParseFlags()
	if err != nil {
		log.Fatal(err)
	}
	config, err := NewConfig(configPath)
	if err != nil {
		log.Fatal(err)
	}

	os.Setenv("WebsocketIP", config.WsServer.Host)
	os.Setenv("WebsocketPORT", config.WsServer.Port)
	os.Setenv("AuthenticationIP", config.AuthServer.Host)
	os.Setenv("AuthenticationPORT", config.AuthServer.Port)
	os.Setenv("DataBaseIP", config.WsServer.Host)
	os.Setenv("DataBasePORT", config.WsServer.Port)

}

// Config struct for webapp config
type Config struct {
	AuthServer struct {
		// Host is the local machine IP Address to bind the HTTP Server to
		Host string `yaml:"host"`

		// Port is the local machine TCP Port to bind the HTTP Server to
		Port string `yaml:"port"`
	} `yaml:"AuthenticationServer"`

	WsServer struct {
		// Host is the local machine IP Address to bind the Websocket Server to
		Host string `yaml:"host"`

		// Port is the local machine TCP Port to bind the Websocket Server to
		Port string `yaml:"port"`
	} `yaml:"WebSocketServer"`

	DbServer struct {
		// Host is the local machine IP Address for connection to the mongodb
		Host string `yaml:"host"`

		// Port is the local machine TCP Port for connection to the mongodb
		Port string `yaml:"port"`
	} `yaml:"DataBaseServer"`
}

// NewConfig returns a new decoded Config struct
func NewConfig(configPath string) (*Config, error) {
	// Create config structure
	config := &Config{}

	// Open config file
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// Init new YAML decode
	d := yaml.NewDecoder(file)

	// Start YAML decoding from file
	if err := d.Decode(&config); err != nil {
		return nil, err
	}

	return config, nil
}

// ValidateConfigPath just makes sure, that the path provided is a file,
// that can be read
func ValidateConfigPath(path string) error {
	s, err := os.Stat(path)
	if err != nil {
		return err
	}
	if s.IsDir() {
		return fmt.Errorf("'%s' is a directory, not a normal file", path)
	}
	return nil
}

// ParseFlags will create and parse the CLI flags
// and return the path to be used elsewhere
func ParseFlags() (string, error) {
	// String that contains the configured configuration path
	var configPath string

	// Set up a CLI flag called "-config" to allow users
	// to supply the configuration file
	flag.StringVar(&configPath, "config", "./config.yml", "path to config file")

	// Actually parse the flags
	flag.Parse()

	// Validate the path first
	if err := ValidateConfigPath(configPath); err != nil {
		return "", err
	}

	// Return the configuration path
	return configPath, nil
}