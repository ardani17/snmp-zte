package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
)

// Config represents application configuration
type Config struct {
	Server  ServerConfig  `json:"server"`
	Redis   RedisConfig   `json:"redis"`
	OLTs    []OLTConfig   `json:"olts"`
}

// ServerConfig represents HTTP server configuration
type ServerConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
}

func (c ServerConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// RedisConfig represents Redis connection configuration
type RedisConfig struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	Password string `json:"password"`
	DB       int    `json:"db"`
}

func (c RedisConfig) Addr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}

// OLTConfig represents an OLT device configuration
type OLTConfig struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Model       string `json:"model"`       // C320, C300, C600
	IPAddress   string `json:"ip_address"`
	Port        int    `json:"port"`
	Community   string `json:"community"`
	BoardCount  int    `json:"board_count"`
	PonPerBoard int    `json:"pon_per_board"`
}

var (
	cfg     *Config
	cfgOnce sync.Once
	cfgPath string
)

// SetConfigPath sets the configuration file path
func SetConfigPath(path string) {
	cfgPath = path
}

// Load loads configuration from file
func Load() (*Config, error) {
	var err error
	cfgOnce.Do(func() {
		cfg, err = loadConfig()
	})
	return cfg, err
}

func loadConfig() (*Config, error) {
	// Default config path
	if cfgPath == "" {
		cfgPath = "config/olts.json"
	}

	// Check if config file exists
	if _, err := os.Stat(cfgPath); os.IsNotExist(err) {
		// Create default config
		return createDefaultConfig()
	}

	// Read config file
	data, err := os.ReadFile(cfgPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	// Set defaults
	if cfg.Server.Port == 0 {
		cfg.Server.Port = 8080
	}

	return &cfg, nil
}

func createDefaultConfig() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Host: "0.0.0.0",
			Port: 8080,
		},
		Redis: RedisConfig{
			Host: "localhost",
			Port: 6379,
			DB:   0,
		},
		OLTs: []OLTConfig{},
	}

	// Ensure directory exists
	dir := filepath.Dir(cfgPath)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create config directory: %w", err)
	}

	// Write default config
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return nil, fmt.Errorf("failed to marshal default config: %w", err)
	}

	if err := os.WriteFile(cfgPath, data, 0644); err != nil {
		return nil, fmt.Errorf("failed to write default config: %w", err)
	}

	return cfg, nil
}

// Save saves configuration to file
func Save(cfg *Config) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(cfgPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}

	return nil
}

// GetOLT returns OLT configuration by ID
func (c *Config) GetOLT(id string) (*OLTConfig, error) {
	for _, olt := range c.OLTs {
		if olt.ID == id {
			return &olt, nil
		}
	}
	return nil, fmt.Errorf("OLT not found: %s", id)
}

// AddOLT adds a new OLT configuration
func (c *Config) AddOLT(olt OLTConfig) error {
	// Check if ID already exists
	for _, existing := range c.OLTs {
		if existing.ID == olt.ID {
			return fmt.Errorf("OLT ID already exists: %s", olt.ID)
		}
	}
	c.OLTs = append(c.OLTs, olt)
	return nil
}

// UpdateOLT updates an existing OLT configuration
func (c *Config) UpdateOLT(id string, olt OLTConfig) error {
	for i, existing := range c.OLTs {
		if existing.ID == id {
			c.OLTs[i] = olt
			return nil
		}
	}
	return fmt.Errorf("OLT not found: %s", id)
}

// DeleteOLT removes an OLT configuration
func (c *Config) DeleteOLT(id string) error {
	for i, olt := range c.OLTs {
		if olt.ID == id {
			c.OLTs = append(c.OLTs[:i], c.OLTs[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("OLT not found: %s", id)
}
