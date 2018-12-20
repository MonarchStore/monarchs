package config

import (
	"fmt"
	"os"

	"github.com/alecthomas/kingpin"
	log "github.com/sirupsen/logrus"
)

const (
	passwordMask = "******"
)

var (
	// Version is the current version, generated at build time
	Version = "unknown"
)

type Config struct {
	ListenAddress string
	ListenPort    int
	LogFormat     string
	LogLevel      string
	LogOutput     string
}

var defaultConfig = &Config{
	ListenAddress: "0.0.0.0",
	ListenPort:    6789,
	LogFormat:     "ascii",
	LogLevel:      "debug",
	LogOutput:     "stderr",
}

func NewConfig() *Config {
	cfg := Config(*defaultConfig)
	return &cfg
}

func NewConfigFromArgs(args []string) *Config {
	cfg := NewConfig()
	err := cfg.ParseFlags(args)
	if err != nil {
		log.Fatalf("Failed to generate config")
	}
	return cfg
}

func (cfg *Config) ParseFlags(args []string) error {
	app := kingpin.New("monarchs", "A hierarchical, NoSQL, in-memory data store with a RESTful API.")
	app.Version(Version)
	app.DefaultEnvars()

	// Basic Server Configs
	// --------------------
	// MONARCHS_ADDR
	app.Flag("addr", "The address/interface to listen on").
		Default(defaultConfig.ListenAddress).
		StringVar(&cfg.ListenAddress)
	// MONARCHS_PORT
	app.Flag("port", "The port to listen on").
		Default(string(defaultConfig.ListenPort)).
		IntVar(&cfg.ListenPort)

	// Logging Configs
	// ---------------
	// MONARCHS_LOG_LEVEL
	app.Flag("log-level", "The log level (trace|debug|info|warning|error|fatal|panic)").
		Default(defaultConfig.LogLevel).
		EnumVar(&cfg.LogLevel, allLogLevelsAsStrings()...)
	app.Flag("log-output", "The log output. Default: 'stderr' (also: 'stdout')").
		Default(defaultConfig.LogOutput).
		StringVar(&cfg.LogOutput)
	app.Flag("log-format", "The log format (ascii|json)").
		Default(defaultConfig.LogFormat).
		StringVar(&cfg.LogFormat)

	app.Parse(args)
	return nil
}

// Returns ':port'
func (cfg *Config) GetListenPort() string {
	return fmt.Sprintf(":%d", cfg.ListenPort)
}

// Returns 'hostname:port'
func (cfg *Config) GetListenAddress() string {
	return cfg.ListenAddress + cfg.GetListenPort()
}

func (cfg *Config) InitLogging() {
	// Set Log Level
	lvl, err := log.ParseLevel(cfg.LogLevel)
	if err != nil {
		lvl, _ = log.ParseLevel("debug")
	}
	log.SetLevel(lvl)

	// Set Log Output
	// Can be any io.Writer
	if cfg.LogOutput == "stdout" {
		log.SetOutput(os.Stdout)
	} else {
		log.SetOutput(os.Stderr)
	}

	// Set JSONFormatter, if desired
	if cfg.LogFormat == "json" {
		log.SetFormatter(&log.JSONFormatter{})
	}
	log.Printf("Logging Configured! [%s]\n", cfg.LogLevel)
}

// Get strings from the log package itself
func allLogLevelsAsStrings() (lvls []string) {
	for _, lvl := range log.AllLevels {
		lvls = append(lvls, lvl.String())
	}
	return
}
