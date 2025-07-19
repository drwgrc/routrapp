package config

import (
	"fmt"
	"os"
	"path/filepath"
	"time"

	"routrapp-api/internal/utils/constants"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Server      ServerConfig   `yaml:"server"`
	CORS        CORSConfig     `yaml:"cors"`
	Database    DatabaseConfig `yaml:"database"`
	Environment string
}

type ServerConfig struct {
	Port         string        `yaml:"port"`
	ReadTimeout  time.Duration `yaml:"read_timeout"`
	WriteTimeout time.Duration `yaml:"write_timeout"`
}

type CORSConfig struct {
	FrontendURL string `yaml:"frontend_url"`
}

// Load loads the configuration from YAML files with environment variable expansion for production
func Load() *Config {
	config := &Config{
		Server: ServerConfig{
			Port:         constants.DefaultPort,
			ReadTimeout:  constants.DefaultReadTimeout,
			WriteTimeout: constants.DefaultWriteTimeout,
		},
		CORS: CORSConfig{
			FrontendURL: constants.DefaultFrontendURL,
		},
		Database: DatabaseConfig{
			Host:         constants.DefaultDBHost,
			Port:         constants.DefaultDBPort,
			User:         constants.DefaultDBUser,
			Password:     constants.DefaultDBPassword,
			DatabaseName: constants.DefaultDBName,
			SSLMode:      constants.DefaultDBSSLMode,
			MaxIdleConns: constants.DefaultDBMaxIdleConns,
			MaxOpenConns: constants.DefaultDBMaxOpenConns,
			ConnMaxLife:  time.Duration(constants.DefaultDBConnMaxLife) * time.Second,
		},
	}

	// Determine environment and load appropriate config files
	env := os.Getenv("ENVIRONMENT")
	if env == "" {
		env = "development"
	}
	config.Environment = env

	// Map environment names to config file suffixes
	configEnv := ""
	switch env {
	case "production":
		configEnv = "prod"
	case "staging":
		configEnv = "staging"
	case "development":
		configEnv = "dev"
	default:
		configEnv = "dev"
	}

	// Load configuration from YAML files
	loadConfigFromYAML(config, configEnv)

	// Expand environment variables for production configs only
	if env == "production" || env == "staging" {
		expandEnvVars(config)
	}

	return config
}

// expandEnvVars expands environment variables in config string fields (for production use)
func expandEnvVars(c *Config) {
	c.Server.Port = os.ExpandEnv(c.Server.Port)
	c.CORS.FrontendURL = os.ExpandEnv(c.CORS.FrontendURL)
	c.Database.Host = os.ExpandEnv(c.Database.Host)
	c.Database.Port = os.ExpandEnv(c.Database.Port)
	c.Database.User = os.ExpandEnv(c.Database.User)
	c.Database.Password = os.ExpandEnv(c.Database.Password)
	c.Database.DatabaseName = os.ExpandEnv(c.Database.DatabaseName)
	c.Database.SSLMode = os.ExpandEnv(c.Database.SSLMode)
}

// loadConfigFromYAML attempts to load configuration from YAML files
func loadConfigFromYAML(config *Config, env string) {
	fmt.Printf("🔧 Environment detected: %s\n", env)
	
	// First try to load the default config
	configPath := "configs/config.yaml"
	fmt.Printf("📄 Loading base config: %s\n", configPath)
	loadYAMLIfExists(config, configPath)

	// Then try to load environment-specific config to override defaults
	envConfigPath := fmt.Sprintf("configs/config.%s.yaml", env)
	fmt.Printf("📄 Loading environment config: %s\n", envConfigPath)
	loadYAMLIfExists(config, envConfigPath)
}

// loadYAMLIfExists loads the YAML configuration if the file exists
func loadYAMLIfExists(config *Config, path string) {
	// Possible file locations (for local dev and Docker container)
	possiblePaths := []string{
		path,                              // configs/config.yaml (Docker container)
		filepath.Join("backend", path),    // backend/configs/config.yaml (local dev)
		filepath.Join("./", path),         // ./configs/config.yaml (alternative)
	}

	var finalPath string
	var found bool

	for _, testPath := range possiblePaths {
		if _, err := os.Stat(testPath); err == nil {
			finalPath = testPath
			found = true
			break
		}
	}

	if !found {
		fmt.Printf("⚠️  Config file not found. Tried: %v\n", possiblePaths)
		return
	}

	fmt.Printf("✅ Loading config file: %s\n", finalPath)
	file, err := os.Open(finalPath)
	if err != nil {
		fmt.Printf("❌ Error opening config file %s: %v\n", finalPath, err)
		return
	}
	defer file.Close()

	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(config); err != nil {
		fmt.Printf("❌ Error decoding YAML file %s: %v\n", finalPath, err)
	} else {
		fmt.Printf("✅ Successfully loaded: %s\n", finalPath)
	}
} 