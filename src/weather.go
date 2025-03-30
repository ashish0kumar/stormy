package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Weather holds the weather data returned by the API
type Weather struct {
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Main struct {
		Temp     float64 `json:"temp"`
		Humidity int     `json:"humidity"`
	} `json:"main"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Sys struct {
		Sunrise int64 `json:"sunrise"`
		Sunset  int64 `json:"sunset"`
	} `json:"sys"`
	Name string `json:"name"`
	Dt   int64  `json:"dt"`
}

// WeatherCache holds cached weather data
type WeatherCache struct {
	Weather   Weather `json:"weather"`
	CacheTime int64   `json:"cache_time"`
}

// fetchWeather fetches weather data from the OpenWeatherMap API
func fetchWeather(config Config) (*Weather, error) {
	// URL encode the city parameter
	encodedCity := url.QueryEscape(config.City)

	apiURL := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&units=%s&APPID=%s",
		encodedCity,
		config.Units,
		config.ApiKey,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, fmt.Errorf("invalid API key - please check your configuration")
	} else if resp.StatusCode == http.StatusNotFound {
		return nil, fmt.Errorf("city '%s' not found - please check the spelling", config.City)
	} else if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status code %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var weather Weather
	if err := json.Unmarshal(body, &weather); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	return &weather, nil
}

// getCachedWeather tries to get weather data from cache
func getCachedWeather(config Config) (*Weather, bool) {
	cachePath := GetCachePath()
	if cachePath == "" {
		return nil, false
	}

	// Check if cache file exists
	if _, err := os.Stat(cachePath); os.IsNotExist(err) {
		return nil, false
	}

	// Read cache file
	data, err := os.ReadFile(cachePath)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to read cache file:", err)
		return nil, false
	}

	var cache WeatherCache
	if err := json.Unmarshal(data, &cache); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to parse cache file:", err)
		return nil, false
	}

	// Check if cache is still valid
	now := time.Now().Unix()
	cacheAge := now - cache.CacheTime
	if cacheAge > config.CacheDuration*60 {
		return nil, false
	}

	// If the city in the config doesn't match the cached city, invalidate the cache
	if cache.Weather.Name != "" && config.City != "" &&
		!strings.EqualFold(cache.Weather.Name, config.City) {
		return nil, false
	}

	return &cache.Weather, true
}

// cacheWeather saves weather data to cache
func cacheWeather(weather *Weather) {
	cachePath := GetCachePath()
	if cachePath == "" {
		return
	}

	// Create directory if it doesn't exist
	cacheDir := filepath.Dir(cachePath)
	if _, err := os.Stat(cacheDir); os.IsNotExist(err) {
		if err := os.MkdirAll(cacheDir, 0755); err != nil {
			fmt.Fprintln(os.Stderr, "Failed to create cache directory:", err)
			return
		}
	}

	// Create cache data
	cache := WeatherCache{
		Weather:   *weather,
		CacheTime: time.Now().Unix(),
	}

	// Marshal cache data
	data, err := json.Marshal(cache)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Failed to marshal cache data:", err)
		return
	}

	// Write cache file
	if err := os.WriteFile(cachePath, data, 0644); err != nil {
		fmt.Fprintln(os.Stderr, "Failed to write cache file:", err)
	}
}
