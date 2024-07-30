package main

import (
	"fmt"
	"log"
	"time"

	"github.com/BurntSushi/toml"
)

// Config holds the configuration values from config.toml
type Config struct {
	Server   ServerConfig
	Paths    PathsConfig
	RabbitMQ RabbitMQConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Address      string
	Port         int
	ReadTimeout  int
	WriteTimeout int
}

// PathsConfig holds the paths configuration
type PathsConfig struct {
	DownloadDirectory string `toml:"download_directory"`
	WebDirectory      string `toml:"web_directory"`
}

// RabbitMQConfig holds the RabbitMQ configuration
type RabbitMQConfig struct {
	URI       string `toml:"uri"`
	QueueName string `toml:"queue_name"`
}

var config Config

// LoadConfig loads the configuration from the given file
func LoadConfig(file string) error {
	if _, err := toml.DecodeFile(file, &config); err != nil {
		return err
	}
	// Log the loaded configuration to verify it
	log.Printf("Loaded configuration: %+v\n", config)
	return nil
}

// GetAddress returns the server address with port
func (c *Config) GetAddress() string {
	return fmt.Sprintf("%s:%d", c.Server.Address, c.Server.Port)
}

// GetReadTimeout returns the read timeout duration
func (c *Config) GetReadTimeout() time.Duration {
	return time.Duration(c.Server.ReadTimeout) * time.Second
}

// GetWriteTimeout returns the write timeout duration
func (c *Config) GetWriteTimeout() time.Duration {
	return time.Duration(c.Server.WriteTimeout) * time.Second
}
