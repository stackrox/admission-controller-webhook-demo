package main

import (
	"os"
	"strings"
	"sync"
)

// Config contains runtime configuration needed for the controller.
type Config struct {
	IgnoredNamespaces []string
}

var (
	config     *Config
	configSync sync.Once
)

// GetConfig returns Config struct, and initializes it first time
func GetConfig() *Config {
	configSync.Do(func() {
		config = &Config{
			IgnoredNamespaces: strings.Split(os.Getenv("IGNORED_NAMESPACES"), ","),
		}
	})

	return config
}
