package weather

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	ProviderOpenWeatherMap = "OpenWeatherMap"
	ProviderOpenMeteo      = "OpenMeteo"
)

type OpenMeteoWeather struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Current   struct {
		Time               string  `json:"time"`
		Interval           int     `json:"interval"`
		Temperature2m      float64 `json:"temperature_2m"`
		WeatherCode        int     `json:"weather_code"`
		Precipitation      float64 `json:"precipitation"`
		RelativeHumidity2m int     `json:"relative_humidity_2m"`
		WindSpeed10m       float64 `json:"wind_speed_10m"`
		WindDirection10m   int     `json:"wind_direction_10m"`
	} `json:"current"`
}

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

func WeatherCodeToSentence(code int) string {
	switch code {
	case 0:
		return "Clear"
	case 1, 2, 3:
		return "Clear"
	case 45, 48:
		return "Clouds"
	case 51, 53, 55:
		return "Clouds"
	case 56, 57:
		return "Clouds"
	case 61, 63, 65:
		return "Rain"
	case 66, 67:
		return "Rain"
	case 71, 73, 75:
		return "Snow"
	case 77:
		return "Snow"
	case 80, 81, 82:
		return "Rain"
	case 85, 86:
		return "Snow"
	case 95:
		return "Thunderstorm"
	case 96, 99:
		return "Thunderstorm"
	default:
		return "Unknown weather code"
	}
}

func ConvertOpenMeteoToWeather(om OpenMeteoWeather, cityName string) Weather {
	return Weather{
		Weather: []struct {
			ID          int    `json:"id"`
			Main        string `json:"main"`
			Description string `json:"description"`
		}{
			{
				ID:          om.Current.WeatherCode,
				Main:        WeatherCodeToSentence(om.Current.WeatherCode),
				Description: WeatherCodeToSentence(om.Current.WeatherCode),
			},
		},
		Main: struct {
			Temp     float64 `json:"temp"`
			Humidity int     `json:"humidity"`
		}{
			Temp:     om.Current.Temperature2m,
			Humidity: om.Current.RelativeHumidity2m,
		},
		Wind: struct {
			Speed float64 `json:"speed"`
			Deg   int     `json:"deg"`
		}{
			Speed: om.Current.WindSpeed10m,
			Deg:   om.Current.WindDirection10m,
		},
		Rain: struct {
			OneHour float64 `json:"1h"`
		}{
			OneHour: om.Current.Precipitation,
		},
		Clouds: struct {
			All int `json:"all"`
		}{
			All: 0, // Not provided by Open-Meteo
		},
		Pop:  0, // Not provided by Open-Meteo
		Name: cityName,
		Dt:   0, // You could parse om.Current.Time to a Unix timestamp if needed
	}
}

func FetchWeatherOpenMeteo(config Config) (*Weather, error) {
	// URL-encode the city name before passing to geocoding
	encodedCity := url.QueryEscape(config.City)

	cityGeo, err := GetFirstGeoResult(encodedCity)
	if err != nil {
		return nil, fmt.Errorf("geocoding failed: %w", err)
	}

	apiURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%f&longitude=%f&current=temperature_2m,weather_code,precipitation,relative_humidity_2m,wind_speed_10m,wind_direction_10m&wind_speed_unit=kmh&temperature_unit=celsius",
		cityGeo.Latitude,
		cityGeo.Longitude,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch weather data: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	var openMeteoWeather OpenMeteoWeather
	if err := json.Unmarshal(body, &openMeteoWeather); err != nil {
		return nil, fmt.Errorf("failed to parse JSON: %w", err)
	}

	weather := ConvertOpenMeteoToWeather(openMeteoWeather, cityGeo.Name)

	return &weather, nil
}

func FetchWeatherOpenWeatherMap(config Config) (*Weather, error) {
	// URL encode the city parameter
	encodedCity := url.QueryEscape(config.City)

	apiURL := fmt.Sprintf(
		"https://api.openweathermap.org/data/2.5/weather?q=%s&units=metric&APPID=%s",
		encodedCity,
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

// FetchWeather fetches weather data from the configured provider
func FetchWeather(config Config) (*Weather, error) {
	if config.Provider == ProviderOpenMeteo {
		return FetchWeatherOpenMeteo(config)
	} else {
		return FetchWeatherOpenWeatherMap(config)
	}
}
