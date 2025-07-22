package weather

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/BurntSushi/toml"
)

// Config holds the application configuration
type Config struct {
	Provider     string `toml:"provider"`
	ApiKey       string `toml:"api_key"`
	City         string `toml:"city"`
	Units        string `toml:"units"`
	ShowCityName bool   `toml:"showcityname"`
	UseColors    bool   `toml:"use_colors"`
	LiveMode     bool   `toml:"live_mode"`
	Compact      bool   `toml:"compact"`
}

// Flags holds command line flags
type Flags struct {
	City, Units            string
	Compact, Help, Version bool
}

const (
	UnitMetric   = "metric"
	UnitImperial = "imperial"
)

var validUnits = [...]string{UnitMetric, UnitImperial}

// DefaultConfig returns a new Config with default values
func DefaultConfig() Config {
	return Config{
		Provider:     ProviderOpenMeteo,
		ApiKey:       "",
		City:         "",
		Units:        UnitMetric,
		ShowCityName: false,
		UseColors:    true,
		LiveMode:     false,
		Compact:      false,
	}
}

// GetConfigPath returns the path to the config file following XDG Base Directory Specification
func GetConfigPath() string {
	var configDir string

	if runtime.GOOS == "windows" {
		// Windows: Use AppData directory
		dir, err := os.UserConfigDir()
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to get config directory:", err)
			dir, err = os.UserHomeDir()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Failed to get home directory:", err)
				return ""
			}
		}
		configDir = filepath.Join(dir, "stormy")
	} else {
		// Linux/macOS: Follow XDG Base Directory Specification
		xdgConfigHome := os.Getenv("XDG_CONFIG_HOME")
		if xdgConfigHome != "" {
			// Use XDG_CONFIG_HOME if set
			configDir = filepath.Join(xdgConfigHome, "stormy")
		} else {
			// Fall back to ~/.config/stormy
			homeDir, err := os.UserHomeDir()
			if err != nil {
				_, _ = fmt.Fprintln(os.Stderr, "Failed to get home directory:", err)
				return ""
			}
			configDir = filepath.Join(homeDir, ".config", "stormy")
		}
	}

	return filepath.Join(configDir, "stormy.toml")
}

// ValidateConfig checks if the config is valid
func ValidateConfig(config *Config) {
	const defaultProvider = ProviderOpenMeteo
	const defaultUnit = UnitMetric

	// Validate provider
	if !slices.Contains(providers[:], config.Provider) {
		_, _ = fmt.Fprintf(
			os.Stderr, "Warning: Invalid provider in config. Using '%s' as default.\n", defaultProvider,
		)
		config.Provider = defaultProvider
	}

	// Validate units
	if !slices.Contains(validUnits[:], config.Units) {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: Invalid units in config. Using '%s' as default.\n", defaultUnit)
		config.Units = defaultUnit
	}

	// Validate API key requirement
	if config.Provider == ProviderOpenWeatherMap && config.ApiKey == "" {
		_, _ = fmt.Fprintf(os.Stderr, "Warning: 'api_key' is required for %s provider.\n", ProviderOpenWeatherMap)
	}
}

// ReadConfig reads/creates the config file and returns the configuration
func ReadConfig() Config {
	configPath := GetConfigPath()
	defaultConfig := DefaultConfig()
	if configPath == "" {
		return defaultConfig
	}

	// Create the directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err = os.MkdirAll(configDir, 0755); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to create config directory:", err)
			return defaultConfig
		}
	}

	// Check if the config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		var file *os.File
		file, err = os.Create(configPath)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to create config file:", err)
			return defaultConfig
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(file)

		if err = toml.NewEncoder(file).Encode(defaultConfig); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to write default config:", err)
			return defaultConfig
		}

		fmt.Printf("No config detected, config created at %s.\n", configPath)
		fmt.Println("Please edit the configuration file to add your API key and city.")
		return defaultConfig
	}

	// Read existing config
	var config Config
	data, err := os.ReadFile(configPath)
	if err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to read config file:", err)
		return defaultConfig
	}

	if err = toml.Unmarshal(data, &config); err != nil {
		_, _ = fmt.Fprintln(os.Stderr, "Failed to parse config file, using defaults with available values:", err)

		// Try to load partial config
		var partialConfig map[string]any

		if err = toml.Unmarshal(data, &partialConfig); err == nil {
			// Apply any valid values from partial config
			if provider, ok := partialConfig["provider"].(string); ok {
				defaultConfig.Provider = provider
			}
			if apiKey, ok := partialConfig["api_key"].(string); ok {
				defaultConfig.ApiKey = apiKey
			}
			if city, ok := partialConfig["city"].(string); ok {
				defaultConfig.City = city
			}
			if units, ok := partialConfig["units"].(string); ok {
				defaultConfig.Units = units
			}
			if showCityName, ok := partialConfig["showcityname"].(bool); ok {
				defaultConfig.ShowCityName = showCityName
			}
			if useColors, ok := partialConfig["use_colors"].(bool); ok {
				defaultConfig.UseColors = useColors
			}
			if liveMode, ok := partialConfig["live_mode"].(bool); ok {
				defaultConfig.LiveMode = liveMode
			}
			if compact, ok := partialConfig["compact"].(bool); ok {
				defaultConfig.Compact = compact
			}
		}

		// Write corrected config back
		var file *os.File
		file, err = os.Create(configPath)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to update config file:", err)
			return defaultConfig
		}
		defer func(f *os.File) {
			_ = f.Close()
		}(file)

		if err = toml.NewEncoder(file).Encode(defaultConfig); err != nil {
			_, _ = fmt.Fprintln(os.Stderr, "Failed to write merged config:", err)
		}

		config = defaultConfig
	}

	// Validate configuration
	ValidateConfig(&config)

	return config
}

// ParseFlags parses command line flags
func ParseFlags() (flags Flags) {
	flag.StringVar(&flags.City, "city", "", "City to get weather for")
	flag.StringVar(&flags.Units, "units", "", fmt.Sprintf("Units (%s)", strings.Join(validUnits[:], ", ")))
	flag.BoolVar(&flags.Compact, "compact", false, "Compact display mode")
	flag.BoolVar(&flags.Help, "help", false, "Show help")
	flag.BoolVar(&flags.Version, "version", false, "Show version information")

	// Add usage information
	flag.Usage = func() {
		_, _ = fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		_, _ = fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		_, _ = fmt.Fprintln(os.Stderr, "\nConfig file is located at:", GetConfigPath())
	}

	flag.Parse()

	if flags.Help {
		flag.Usage()
		os.Exit(0)
	}

	return
}

// ApplyFlags applies command line flags to the config
func ApplyFlags(config *Config, flags Flags) {
	if flags.City != "" {
		config.City = flags.City
	}
	if flags.Units != "" {
		config.Units = flags.Units
		ValidateConfig(config)
	}
	if flags.Compact {
		config.Compact = true
	}
}
