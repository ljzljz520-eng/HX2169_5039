package main

import (
	"errors"
	"fmt"
	"net/url"
	"strings"
)

type RuntimeConfig struct {
	Address   string
	Database  string
	StaticDir string
}

func DefaultConfig() RuntimeConfig {
	return RuntimeConfig{Address: "127.0.0.1:8080", Database: "graduation-showcase.db", StaticDir: "web"}
}

func (c RuntimeConfig) Validate() error {
	if strings.TrimSpace(c.Address) == "" {
		return errors.New("address is required")
	}
	if strings.TrimSpace(c.Database) == "" {
		return errors.New("database is required")
	}
	return nil
}

func ParseAddress(value string) (string, error) {
	if !strings.Contains(value, ":") {
		return "", errors.New("address requires port")
	}
	parsed, err := url.Parse("http://" + value)
	if err != nil || parsed.Host == "" {
		return "", errors.New("invalid address")
	}
	return parsed.Host, nil
}

func DescribeConfig(config RuntimeConfig) string {
	return fmt.Sprintf("address=%s database=%s static=%s", config.Address, config.Database, config.StaticDir)
}

func WithDatabase(config RuntimeConfig, database string) RuntimeConfig {
	config.Database = database
	return config
}
