package weather

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

const (
	ProviderOpenWeatherMap = "OpenWeatherMap"
	ProviderOpenMeteo      = "OpenMeteo"
)

var providers = [...]string{ProviderOpenWeatherMap, ProviderOpenMeteo}

var ErrUnsupportedQuery = errors.New("unsupported query")

type GeoResult struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

type GeoResponse struct {
	Results []GeoResult `json:"results"`
}

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

type OpenWeatherMapGeolocationResult struct {
	City       string            `json:"name"`
	LocalNames map[string]string `json:"local_names"`
	Latitude   float64           `json:"lat"`
	Longitude  float64           `json:"lon"`
	Country    string            `json:"country"`
	State      string            `json:"state"`
}

type OpenWeatherMapWeather struct {
	Coordinates struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temperature          float64 `json:"temp"`
		FeelsLikeTemperature float64 `json:"feels_like"`
		TempMin              float64 `json:"temp_min"`
		TempMax              float64 `json:"temp_max"`
		Pressure             int     `json:"pressure"`
		Humidity             int     `json:"humidity"`
		SeaLevelPressure     int     `json:"sea_level"`
		GroundLevelPressure  int     `json:"grnd_level"`
	} `json:"main"`
	VisibilityDistance int `json:"visibility"`
	Wind               struct {
		Speed   float64 `json:"speed"`
		Degrees int     `json:"deg"`
		Gust    float64 `json:"gust"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Rain struct {
		Precipitations float64 `json:"1h,omitempty"`
	} `json:"rain,omitempty"`
	Snow struct {
		Precipitations float64 `json:"1h,omitempty"`
	} `json:"snow,omitempty"`
	CalculationDate int `json:"dt"`
	Sys             struct {
		Type        int    `json:"type"`
		Id          int    `json:"id"`
		Country     string `json:"country"`
		SunriseTime int    `json:"sunrise"`
		SunsetTime  int    `json:"sunset"`
	} `json:"sys"`
	TimezoneShift int    `json:"timezone"`
	ID            int    `json:"id"`
	City          string `json:"name"`
	Cod           int    `json:"cod"`
}

// Weather holds the weather data returned by the API
type Weather struct {
	Weather []struct {
		ID                int
		Main, Description string
	}
	Main struct {
		Temp     float64
		Humidity int
	}
	Wind struct {
		Speed float64
		Deg   int
	}
	Rain struct {
		OneHour float64
	}
	Clouds struct {
		All int
	}
	Pop  float64
	Name string
	Dt   int64
}

func fetchAndUnmarshal[T any](u string, args ...any) (out T, err error) {
	var resp *http.Response
	resp, err = http.Get(fmt.Sprintf(u, args...))
	if err != nil {
		return
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	if resp.StatusCode == http.StatusUnauthorized {
		err = fmt.Errorf("invalid API key - please check your configuration")
		return
	} else if resp.StatusCode == http.StatusNotFound {
		err = fmt.Errorf("no data found - please check your input")
		return
	} else if resp.StatusCode != http.StatusOK {
		err = fmt.Errorf("API returned status code %d", resp.StatusCode)
		return
	}

	err = json.NewDecoder(resp.Body).Decode(&out)
	return
}

func GetFirstGeoResult(encodedCity string) (*GeoResult, error) {
	geo, err := fetchAndUnmarshal[GeoResponse](
		"https://geocoding-api.open-meteo.com/v1/search?name=%s&count=1", encodedCity,
	)
	if err != nil {
		return nil, err
	}

	if len(geo.Results) == 0 {
		return nil, fmt.Errorf("no results found for city query \"%s\"", encodedCity)
	}

	return &geo.Results[0], nil
}

func CodeToSentence(code int) string {
	switch code {
	case 0, 1, 2, 3:
		return ConditionClear
	case 45, 48, 51, 53, 55, 56, 57:
		return ConditionClouds
	case 61, 63, 65, 66, 67:
		return ConditionRain
	case 71, 73, 75, 77:
		return ConditionSnow
	case 80, 81, 82:
		return ConditionRain
	case 85, 86:
		return ConditionSnow
	case 95, 96, 99:
		return ConditionThunderstorm
	default:
		return "Unknown weather code"
	}
}

func ConvertOpenMeteoToWeather(om OpenMeteoWeather, cityName string) Weather {
	return Weather{
		Weather: []struct {
			ID          int
			Main        string
			Description string
		}{
			{
				ID:          om.Current.WeatherCode,
				Main:        CodeToSentence(om.Current.WeatherCode),
				Description: CodeToSentence(om.Current.WeatherCode),
			},
		},
		Main: struct {
			Temp     float64
			Humidity int
		}{
			Temp:     om.Current.Temperature2m,
			Humidity: om.Current.RelativeHumidity2m,
		},
		Wind: struct {
			Speed float64
			Deg   int
		}{
			Speed: om.Current.WindSpeed10m,
			Deg:   om.Current.WindDirection10m,
		},
		Rain: struct {
			OneHour float64
		}{
			OneHour: om.Current.Precipitation,
		},
		Clouds: struct {
			All int
		}{
			All: 0, // Not provided by Open-Meteo
		},
		Pop:  0, // Not provided by Open-Meteo
		Name: cityName,
		Dt:   0, // You could parse om.Current.Time to a Unix timestamp if needed
	}
}

func ConvertOpenWeatherMapToWeather(om OpenWeatherMapWeather, cityName string) Weather {
	return Weather{
		Weather: []struct {
			ID          int
			Main        string
			Description string
		}(om.Weather),
		Main: struct {
			Temp     float64
			Humidity int
		}{
			Temp:     om.Main.Temperature,
			Humidity: om.Main.Humidity,
		},
		Wind: struct {
			Speed float64
			Deg   int
		}{
			Speed: om.Wind.Speed,
			Deg:   om.Wind.Degrees,
		},
		Rain: struct {
			OneHour float64
		}{
			OneHour: om.Rain.Precipitations,
		},
		Clouds: struct {
			All int
		}{
			All: om.Clouds.All,
		},
		Pop:  0, // Not provided (anymore) by OpenWeatherMap
		Name: cityName,
		Dt:   int64(om.CalculationDate),
	}
}

func FetchWeatherOpenMeteo(config Config) (*Weather, error) {
	var latitude, longitude any
	var city string

	if config.City != "" {
		// URL-encode the city name before passing to geocoding
		encodedCity := url.QueryEscape(config.City)

		cityGeo, err := GetFirstGeoResult(encodedCity)
		if err != nil {
			if strings.Contains(config.City, " ") || strings.Contains(config.City, ",") {
				return nil, fmt.Errorf("geocoding failed - %w: %w", ErrUnsupportedQuery, err)
			}
			return nil, fmt.Errorf("geocoding failed: %w", err)
		}

		latitude, longitude, city = cityGeo.Latitude, cityGeo.Longitude, cityGeo.Name
	} else {
		latitude, longitude = config.Latitude, config.Longitude
	}

	openMeteoWeather, err := fetchAndUnmarshal[OpenMeteoWeather](
		"https://api.open-meteo.com/v1/forecast?latitude=%v&longitude=%v&current=temperature_2m,weather_code,precipitation,relative_humidity_2m,wind_speed_10m,wind_direction_10m&wind_speed_unit=kmh&temperature_unit=celsius",
		latitude,
		longitude,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch or decode data: %w", err)
	}

	weather := ConvertOpenMeteoToWeather(openMeteoWeather, city)

	return &weather, nil
}

func FetchWeatherOpenWeatherMap(config Config) (*Weather, error) {
	var longitude, latitude any
	var city string

	if config.City != "" {
		// URL encode the city parameter
		encodedCity := url.QueryEscape(config.City)

		// Geocoding
		geoResult, err := fetchAndUnmarshal[[]OpenWeatherMapGeolocationResult](
			"https://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s",
			encodedCity,
			config.ApiKey,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to fetch or decode data for geocoding: %w", err)
		}

		if len(geoResult) == 0 {
			return nil, fmt.Errorf("no results found for city %s", config.City)
		}

		latitude, longitude, city = geoResult[0].Latitude, geoResult[0].Longitude, geoResult[0].City
	} else {
		latitude, longitude = config.Longitude, config.Latitude
	}

	// Actual weather
	openWeatherMapWeather, err := fetchAndUnmarshal[OpenWeatherMapWeather](
		"https://api.openweathermap.org/data/2.5/weather?lat=%v&lon=%v&units=metric&appid=%s",
		latitude,
		longitude,
		config.ApiKey,
	)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch or decode data: %w", err)
	}

	weather := ConvertOpenWeatherMapToWeather(openWeatherMapWeather, city)

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
