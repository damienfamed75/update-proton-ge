package main

import (
	"os"
	"strings"

	"github.com/rs/zerolog/log"
)

const (
	// _protonDir = "/protonGE"
	_protonDir             = "/.local/proton-ge"
	_compatibilityToolsDir = "/.steam/root/compatibilitytools.d/"
	// _lynxCmd = "lynx"

	_envProtonDir             = "PROTON_GE_ARCHIVE"
	_envCompatibilityToolsDir = "COMPATIBILITY_TOOLS_DIR"
)

func init() {
	c := newConfig()

	found, ok := os.LookupEnv(_envProtonDir)
	if ok {
		c.protonDir = found
	}

	found, ok = os.LookupEnv(_envCompatibilityToolsDir)
	if ok {
		c.compatibilityToolsDir = found
	}

	cfg = c
}

var cfg *Config

type Config struct {
	protonDir             string
	homeDir               string
	compatibilityToolsDir string
}

func newConfig() *Config {
	homeDir, ok := os.LookupEnv("HOME")
	if !ok {
		log.Fatal().Msg("could not find home directory")
	}

	return &Config{
		homeDir:               homeDir,
		protonDir:             homeDir + _protonDir,
		compatibilityToolsDir: homeDir + _compatibilityToolsDir,
	}
}

func (c *Config) ProtonDir() string {
	return strings.TrimSuffix(c.protonDir, "/")
}

func (c *Config) CompatibilityToolsDir() string {
	return strings.TrimSuffix(c.compatibilityToolsDir, "/")
}

func (c *Config) Home() string {
	return c.homeDir
}
