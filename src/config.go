package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/BurntSushi/toml"
)

// Config holds the application configuration
type Config struct {
	ApiKey        string `toml:"api_key"`
	City          string `toml:"city"`
	Units         string `toml:"units"`
	TimePlus      int64  `toml:"timeplus"`
	TimeMinus     int64  `toml:"timeminus"`
	ShowCityName  bool   `toml:"showcityname"`
	ShowDate      bool   `toml:"showdate"`
	TimeFormat    string `toml:"timeformat"`
	UseColors     bool   `toml:"use_colors"`
	CacheDuration int64  `toml:"cache_duration"`
}

// Flags holds command line flags
type Flags struct {
	ApiKey  string
	City    string
	Units   string
	NoCache bool
	Help    bool
}

// DefaultConfig returns a new Config with default values
func DefaultConfig() Config {
	return Config{
		ApiKey:        "",
		City:          "",
		Units:         "metric",
		TimePlus:      0,
		TimeMinus:     0,
		ShowCityName:  false,
		ShowDate:      false,
		TimeFormat:    "24",
		UseColors:     false,
		CacheDuration: 30,
	}
}

// GetConfigPath returns the path to the config file
func GetConfigPath() string {
	var configDir string

	if runtime.GOOS == "windows" {
		dir, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get config directory:", err)
			dir, err = os.UserHomeDir()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to get home directory:", err)
				return ""
			}
			return filepath.Join(dir, "stormy", "stormy.toml")
		}
		configDir = filepath.Join(dir, "stormy")
	} else {
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get home directory:", err)
			return ""
		}
		configDir = filepath.Join(dir, ".config", "stormy")
	}

	return filepath.Join(configDir, "stormy.toml")
}

// GetCachePath returns the path to the cache file
func GetCachePath() string {
	var configDir string

	if runtime.GOOS == "windows" {
		dir, err := os.UserConfigDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get config directory:", err)
			dir, err = os.UserHomeDir()
			if err != nil {
				fmt.Fprintln(os.Stderr, "Failed to get home directory:", err)
				return ""
			}
			return filepath.Join(dir, "stormy", "cache.json")
		}
		configDir = filepath.Join(dir, "stormy")
	} else {
		dir, err := os.UserHomeDir()
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to get home directory:", err)
			return ""
		}
		configDir = filepath.Join(dir, ".config", "stormy")
	}

	return filepath.Join(configDir, "cache.json")
}

// ValidateConfig checks if the config is valid
func ValidateConfig(config *Config) {
	// Validate units
	validUnits := map[string]bool{
		"metric":   true,
		"imperial": true,
		"standard": true,
	}

	if !validUnits[config.Units] {
		fmt.Fprintln(os.Stderr, "Warning: Invalid units in config. Using 'metric' as default.")
		config.Units = "metric"
	}

	// Validate time format
	validTimeFormats := map[string]bool{
		"12": true,
		"24": true,
	}

	if !validTimeFormats[config.TimeFormat] {
		fmt.Fprintln(os.Stderr, "Warning: Invalid time format in config. Using '24' as default.")
		config.TimeFormat = "24"
	}

	// Validate cache duration
	if config.CacheDuration < 0 {
		fmt.Fprintln(os.Stderr, "Warning: Negative cache duration in config. Using 30 minutes as default.")
		config.CacheDuration = 30
	}
}

// ReadConfig reads/creates the config file and returns the configuration
func ReadConfig() Config {
	configPath := GetConfigPath()
	if configPath == "" {
		return DefaultConfig()
	}

	// Create directory if it doesn't exist
	configDir := filepath.Dir(configPath)
	if _, err := os.Stat(configDir); os.IsNotExist(err) {
		if err := os.MkdirAll(configDir, 0755); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create config directory:", err)
			return DefaultConfig()
		}
	}

	// Check if config file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		// Create default config
		defaultConfig := DefaultConfig()
		file, err := os.Create(configPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create config file:", err)
			return defaultConfig
		}
		defer file.Close()

		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(defaultConfig); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write default config:", err)
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
		fmt.Fprintln(os.Stderr, "Failed to read config file:", err)
		return DefaultConfig()
	}

	if err := toml.Unmarshal(data, &config); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse config file, using defaults with available values:", err)

		// Try to load partial config
		defaultConfig := DefaultConfig()
		var partialConfig map[string]interface{}

		if err := toml.Unmarshal(data, &partialConfig); err == nil {
			// Apply any valid values from partial config
			if apiKey, ok := partialConfig["api_key"].(string); ok {
				defaultConfig.ApiKey = apiKey
			}
			if city, ok := partialConfig["city"].(string); ok {
				defaultConfig.City = city
			}
			if units, ok := partialConfig["units"].(string); ok {
				defaultConfig.Units = units
			}
			if timeplus, ok := partialConfig["timeplus"].(int64); ok {
				defaultConfig.TimePlus = timeplus
			}
			if timeminus, ok := partialConfig["timeminus"].(int64); ok {
				defaultConfig.TimeMinus = timeminus
			}
			if showcityname, ok := partialConfig["showcityname"].(bool); ok {
				defaultConfig.ShowCityName = showcityname
			}
			if showdate, ok := partialConfig["showdate"].(bool); ok {
				defaultConfig.ShowDate = showdate
			}
			if timeformat, ok := partialConfig["timeformat"].(string); ok {
				defaultConfig.TimeFormat = timeformat
			}
			if useColors, ok := partialConfig["use_colors"].(bool); ok {
				defaultConfig.UseColors = useColors
			}
			if cacheDuration, ok := partialConfig["cache_duration"].(int64); ok {
				defaultConfig.CacheDuration = cacheDuration
			}
		}

		// Write corrected config back
		file, err := os.Create(configPath)
		if err != nil {
			fmt.Fprintln(os.Stderr, "Failed to update config file:", err)
			return defaultConfig
		}
		defer file.Close()

		encoder := toml.NewEncoder(file)
		if err := encoder.Encode(defaultConfig); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to write merged config:", err)
		}

		config = defaultConfig
	}

	// Validate configuration
	ValidateConfig(&config)

	return config
}

// parseFlags parses command line flags
func parseFlags() Flags {
	flags := Flags{}

	flag.StringVar(&flags.ApiKey, "key", "", "OpenWeatherMap API key")
	flag.StringVar(&flags.City, "city", "", "City to get weather for")
	flag.StringVar(&flags.Units, "units", "", "Units (metric, imperial, standard)")
	flag.BoolVar(&flags.NoCache, "no-cache", false, "Disable cache")
	flag.BoolVar(&flags.Help, "help", false, "Show help")

	// Add usage information
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options]\n\n", os.Args[0])
		fmt.Fprintln(os.Stderr, "Options:")
		flag.PrintDefaults()
		fmt.Fprintln(os.Stderr, "\nConfig file is located at:", GetConfigPath())
	}

	flag.Parse()

	if flags.Help {
		flag.Usage()
		os.Exit(0)
	}

	return flags
}

// applyFlags applies command line flags to the config
func applyFlags(config *Config, flags Flags) {
	if flags.ApiKey != "" {
		config.ApiKey = flags.ApiKey
	}
	if flags.City != "" {
		config.City = flags.City
	}
	if flags.Units != "" {
		config.Units = flags.Units
		ValidateConfig(config)
	}
}
