package main

import (
	"fmt"
	"os"

	"github.com/ashish0kumar/stormy/internal/weather"
)

func main() {
	// Parse command line flags
	flags := weather.ParseFlags()

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

	// Fetch weather data
	weatherData, err := weather.FetchWeather(config)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to fetch weather data: %v\n", err)
		fmt.Fprintln(os.Stderr, "Please check your internet connection and API key.")
		os.Exit(1)
	}

	// Display the weather
	weather.DisplayWeather(weatherData, config)
}
