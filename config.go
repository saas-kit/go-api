package main

import (
	"fmt"
	"log"

	"github.com/jinzhu/configor"
)

// GetConfig returns a new instance of the Config structure
func GetConfig() *Config {
	cnf := &Config{}
	err := configor.New(&configor.Config{ErrorOnUnmatchedKeys: false}).
		Load(cnf, ".env")
	if err != nil {
		panic(err)
	}
	log.Println(fmt.Print(cnf))
	return cnf
}

// Config structure
type Config struct {
	Name string `env:"APP_NAME"`
	App  struct {
		Name    string `env:"APP_NAME"`
		Env     string `env:"APP_ENV"`
		BaseURL string `env:"APP_URL"`
		Version string `env:"APP_VERSION"`
		Port    int    `env:"APP_PORT"`
	}
}

// ListenAddr return address to listen and serve
func (c *Config) ListenAddr() string {
	return c.App.BaseURL
}

// APIVersion returns current API version
func (c *Config) APIVersion() string {
	return c.App.Version
}
