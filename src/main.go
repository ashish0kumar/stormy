package main

import (
	"fmt"
	"os"
	"time"
)

func main() {
	// Parse command line flags
	flags := parseFlags()

	// Read/create config
	config := ReadConfig()

	// Override config with command line flags if provided
	applyFlags(&config, flags)

	// Check if API key and city are set
	if config.ApiKey == "" || config.City == "" {
		fmt.Fprintln(os.Stderr, "Error: API key and city name must be set in the config file or via command line flags.")
		fmt.Fprintln(os.Stderr, "Get your API key from https://openweathermap.org/api and set it in the config file.")
		fmt.Fprintf(os.Stderr, "Run '%s --help' for usage information.\n", os.Args[0])
		os.Exit(1)
	}

	// Fetch weather data
	weather, err := fetchWeather(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch weather data: %v\n", err)
		fmt.Fprintln(os.Stderr, "Please check your internet connection and API key.")
		os.Exit(1)
	}

	// Calculate timezone adjustment
	adjustment := time.Duration(config.TimePlus-config.TimeMinus) * time.Hour

	// Display the weather
	displayWeather(weather, config, adjustment)
}
