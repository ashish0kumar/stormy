package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashish0kumar/stormy/internal/weather"
	"github.com/k0kubun/go-ansi"
)

// version is set during build time using -ldflags
var version = "dev"

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		// reset cursor visibility for live mode
		_, _ = ansi.Print("\x1b[?25h")
		os.Exit(1)
	}()
}

func main() {
	// Parse command line flags
	flags := weather.ParseFlags()

	// Handle version flag
	if flags.Version {
		fmt.Printf("stormy version %s\n", version)
		os.Exit(0)
	}

	// Read/create config
	config := weather.ReadConfig()

	// Override config with command line flags if provided
	weather.ApplyFlags(&config, flags)

	// Check if city is set
	if config.City == "" {
		fmt.Fprintln(os.Stderr, "Error: City must be set in the config file or via command line flags")
		fmt.Fprintln(os.Stderr, "Config file location:", weather.GetConfigPath())
		fmt.Fprintf(os.Stderr, "Run '%s --help' for usage information.\n", os.Args[0])
		os.Exit(1)
	}

	// Check if API key and city are set
	if config.Provider == weather.ProviderOpenWeatherMap && config.ApiKey == "" {
		fmt.Fprintln(os.Stderr, "Error: API key must be set in the config file when using OpenWeatherMap")
		fmt.Fprintln(os.Stderr, "Get your API key from https://openweathermap.org/api")
		fmt.Fprintln(os.Stderr, "Config file location:", weather.GetConfigPath())
		fmt.Fprintf(os.Stderr, "Run '%s --help' for usage information.\n", os.Args[0])
	}

	fetchAndDisplay(config, false)
}

// fetchAndDisplay fetches weather data and displays it according to the given configuration.
// clear determines whether the screen should be cleared before displaying updated information.
func fetchAndDisplay(config weather.Config, clear bool) {
	// Fetch weather data
	weatherData, err := weather.FetchWeather(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch weather data: %v\n", err)
		fmt.Fprintln(os.Stderr, "Please check your internet connection and API key.")
		os.Exit(1)
	}

	// Clear screen in live mode
	if clear {
		_, _ = ansi.Printf("\x1b[%dA\x1b[J", 7) // maximum number of displayed lines
	}

	// Display the weather
	weather.DisplayWeather(weatherData, config)

	// Loop in live mode
	if !config.LiveMode {
		return
	}
	if !clear {
		// hide cursor on live mode startup
		_, _ = ansi.Print("\x1b[?25l")
	}
	time.Sleep(15 * time.Second)
	fetchAndDisplay(config, true)
}
