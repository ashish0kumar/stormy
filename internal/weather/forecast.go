package weather

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type OpenMeteo7DayResponse struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Daily     struct {
		Time                     []string  `json:"time"`
		WeatherCode              []int     `json:"weather_code"`
		Temperature2mMin         []float64 `json:"temperature_2m_min"`
		Temperature2mMax         []float64 `json:"temperature_2m_max"`
		PrecipitationProbability []int     `json:"precipitation_probability_max"`
	} `json:"daily"`
}

type DayForecast struct {
	Day       string
	Condition string
	Icon      []string
	MinTemp   int
	MaxTemp   int
	Precip    int
}

var (
	smallSunIcon = []string{
		"\033[38;5;226m   \\  /   \033[0m",
		"\033[38;5;226m  --(  )-- \033[0m",
		"\033[38;5;226m   /  \\   \033[0m",
	}
	smallPartlyCloudyIcon = []string{
		"\033[38;5;226m  \\  /   \033[0m",
		"\033[38;5;226m _/\"\"\033[38;5;250m.-. \033[0m",
		"\033[38;5;226m  \\_\033[38;5;250m(   )\033[0m",
	}
	smallCloudyIcon = []string{
		"\033[38;5;250m  .--.   \033[0m",
		"\033[38;5;250m .-(  ). \033[0m",
		"\033[38;5;250m(___.__) \033[0m",
	}
	smallVeryCloudyIcon = []string{
		"\033[38;5;240;1m  .--.   \033[0m",
		"\033[38;5;240;1m .-(  ). \033[0m",
		"\033[38;5;240;1m(___.__) \033[0m",
	}
	smallLightShowersIcon = []string{
		"\033[38;5;226m _`/\"\"\033[38;5;250m.-. \033[0m",
		"\033[38;5;226m  ,\\_\033[38;5;250m(   )\033[0m",
		"\033[38;5;111m  ' ' ' ' \033[0m",
	}
	smallHeavyShowersIcon = []string{
		"\033[38;5;226m _`/\"\"\033[38;5;240;1m.-. \033[0m",
		"\033[38;5;226m  ,\\_\033[38;5;240;1m(   )\033[0m",
		"\033[38;5;21;1m  ‚'‚'‚'‚' \033[0m",
	}
	smallLightSnowIcon = []string{
		"\033[38;5;250m  .-.   \033[0m",
		"\033[38;5;250m (   ). \033[0m",
		"\033[38;5;255m *  *  * \033[0m",
	}
	smallHeavySnowIcon = []string{
		"\033[38;5;240;1m  .-.   \033[0m",
		"\033[38;5;240;1m (   ). \033[0m",
		"\033[38;5;255;1m * * * * \033[0m",
	}
	smallThunderIcon = []string{
		"\033[38;5;240;1m  .-.   \033[0m",
		"\033[38;5;240;1m (   ). \033[0m",
		"\033[38;5;228;5m ⚡\"\"⚡\"\" \033[0m",
	}
	smallFogIcon = []string{
		"\033[38;5;251m _ - _ - \033[0m",
		"\033[38;5;251m  _ - _  \033[0m",
		"\033[38;5;251m _ - _ - \033[0m",
	}
	smallUnknownIcon = []string{
		"  .-.   ",
		"  (   )  ",
		"  `-’   ",
	}
)

func mapWeatherCodeToSmallIcon(code int) []string {
	switch code {
	case 0:
		return smallSunIcon
	case 1:
		return smallPartlyCloudyIcon
	case 2, 3:
		return smallCloudyIcon
	case 45, 48:
		return smallFogIcon
	case 51, 53, 55, 56, 57:
		return smallLightShowersIcon
	case 61, 63, 65, 66, 67:
		return smallLightShowersIcon
	case 71, 73, 75, 77:
		return smallLightSnowIcon
	case 80, 81, 82:
		return smallLightShowersIcon
	case 85, 86:
		return smallHeavySnowIcon
	case 95, 96, 99:
		return smallThunderIcon
	default:
		return smallUnknownIcon
	}
}

func mapWeatherCodeToText(code int) string {
	switch code {
	case 0:
		return "Sunny"
	case 1:
		return "Partly Cloudy"
	case 2, 3:
		return "Cloudy"
	case 45, 48:
		return "Foggy"
	case 51, 53, 55, 56, 57:
		return "Drizzle"
	case 61, 63, 65, 66, 67:
		return "Rain"
	case 71, 73, 75, 77:
		return "Light Snow"
	case 80, 81, 82:
		return "Rain Showers"
	case 85, 86:
		return "Heavy Snow"
	case 95, 96, 99:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}

func Fetch7DayForecast(cityName string) ([]DayForecast, error) {
	geo, err := GetFirstGeoResult(url.QueryEscape(cityName))
	if err != nil {
		return nil, err
	}

	apiURL := fmt.Sprintf(
		"https://api.open-meteo.com/v1/forecast?latitude=%.6f&longitude=%.6f&daily=weather_code,temperature_2m_max,temperature_2m_min,precipitation_probability_max&forecast_days=7&timezone=auto",
		geo.Latitude, geo.Longitude,
	)

	resp, err := http.Get(apiURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var data OpenMeteo7DayResponse
	if err := json.NewDecoder(resp.Body).Decode(&data); err != nil {
		return nil, err
	}

	var forecasts []DayForecast
	for i := 0; i < len(data.Daily.Time) && i < 7; i++ {
		t, _ := time.Parse("2006-01-02", data.Daily.Time[i])
		dayName := t.Format("Mon")

		code := data.Daily.WeatherCode[i]
		condition := mapWeatherCodeToText(code)
		icon := mapWeatherCodeToSmallIcon(code)

		precip := 0
		if i < len(data.Daily.PrecipitationProbability) {
			precip = data.Daily.PrecipitationProbability[i]
		}

		forecasts = append(forecasts, DayForecast{
			Day:       dayName,
			Condition: condition,
			Icon:      icon,
			MinTemp:   int(data.Daily.Temperature2mMin[i]),
			MaxTemp:   int(data.Daily.Temperature2mMax[i]),
			Precip:    precip,
		})
	}

	return forecasts, nil
}

func Render7DayForecast(cityName string, forecasts []DayForecast) {
	if len(forecasts) == 0 {
		fmt.Println("No forecast data available.")
		return
	}

	reset := "\033[0m"
	bold := "\033[1m"
	gray := "\033[38;5;245m"
	yellow := "\033[38;5;226m"
	cyan := "\033[38;5;51m"
	orange := "\033[38;5;208m"

	coldColor := "\033[38;5;39;1m"
	hotColor := "\033[38;5;196;1m"

	globalMin := forecasts[0].MinTemp
	globalMax := forecasts[0].MaxTemp
	for _, f := range forecasts {
		if f.MinTemp < globalMin {
			globalMin = f.MinTemp
		}
		if f.MaxTemp > globalMax {
			globalMax = f.MaxTemp
		}
	}

	line := gray + "─────────────────────────────────────────────────────────────" + reset

	fmt.Printf(bold+"\n 📅 7-DAY EXTENDED FORECAST (%s)\n"+reset, strings.ToUpper(cityName))
	fmt.Println(line)
	fmt.Printf(" %s%-5s  %-12s  %-18s  %-14s  %-6s%s\n",
		bold, "DAY", "ICON", "CONDITION", "TEMP (MIN/MAX)", "PRECIP", reset)
	fmt.Println(line)

	for _, f := range forecasts {
		condColor := yellow
		if strings.Contains(f.Condition, "Rain") || strings.Contains(f.Condition, "Shower") || strings.Contains(f.Condition, "Drizzle") {
			condColor = cyan
		} else if strings.Contains(f.Condition, "Thunder") {
			condColor = orange
		} else if strings.Contains(f.Condition, "Cloud") || strings.Contains(f.Condition, "Fog") {
			condColor = gray
		} else if strings.Contains(f.Condition, "Snow") {
			condColor = "\033[38;5;255m"
		}

		minC := ""
		maxC := ""
		if f.MinTemp == globalMin {
			minC = coldColor
		}
		if f.MaxTemp == globalMax {
			maxC = hotColor
		}

		fmt.Printf(" %-5s  %s  %s%-18s%s  %s%2d°C%s / %s%2d°C%s        %2d%%\n",
			f.Day,
			f.Icon[0],
			condColor, f.Condition, reset,
			minC, f.MinTemp, reset,
			maxC, f.MaxTemp, reset,
			f.Precip,
		)

		fmt.Printf("        %s\n", f.Icon[1])
		fmt.Printf("        %s\n", f.Icon[2])

		fmt.Println(line)
	}
}