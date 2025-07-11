package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	ProviderOpenWeather = "OpenWeatherApi"
	ProviderOpenMeteo   = "OpenMeteo"
)

// Weather holds the weather data returned by the API
type Weather struct {
	Weather []struct {
		ID          int    `json:"id"`
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
	Rain struct {
		OneHour float64 `json:"1h"`
	} `json:"rain"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Pop  float64 `json:"pop"`
	Name string  `json:"name"`
	Dt   int64   `json:"dt"`
}

type GeoResult struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	// ...add more fields if needed
}

type GeoResponse struct {
	Results []GeoResult `json:"results"`
}

func GetFirstGeoResult(encodedCity string) (*GeoResult, error) {
	cityUrl := fmt.Sprintf("https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", encodedCity)

	resp, err := http.Get(cityUrl)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var geo GeoResponse
	if err := json.Unmarshal(body, &geo); err != nil {
		return nil, err
	}

	if len(geo.Results) == 0 {
		return nil, fmt.Errorf("no results found for city query %s", encodedCity)
	}

	return &geo.Results[0], nil
}

// FetchWeather fetches weather data from the OpenWeatherMap API
func FetchWeather(config Config) (*Weather, error) {
	// URL encode the city parameter
	encodedCity := url.QueryEscape(config.City)

	if config.Provider == ProviderOpenMeteo {

		cityGeo, err := GetFirstGeoResult(encodedCity)
		if err != nil {
			fmt.Println("Error:", err)
		}

		apiURL := fmt.Sprintf(
			"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,weather_code,precipitation,relative_humidity_2m,wind_speed_10m,wind_direction_10m&temperature_unit=fahrenheit",
			cityGeo.Latitude,
			cityGeo.Longitude,
		)

		resp, err := http.Get(apiURL)

		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("failed to read response body: %w", err)
		}

		fmt.Print(string(body))

	}

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
