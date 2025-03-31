package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

// formatTime formats a Unix timestamp according to the config
func formatTime(timestamp int64, config Config, adjustment time.Duration) string {
	t := time.Unix(timestamp, 0).Add(adjustment)

	if config.TimeFormat == "12" {
		return t.Format("3:04 PM")
	}
	return t.Format("15:04")
}

// getWindDirectionSymbol converts wind degrees to a direction symbol
func getWindDirectionSymbol(degrees int) string {
	directions := []string{"↑", "↗", "→", "↘", "↓", "↙", "←", "↖"}
	index := int((float64(degrees)+22.5)/45.0) % 8
	return directions[index]
}

// displayWeather renders the weather data with ASCII art
func displayWeather(weather *Weather, config Config, adjustment time.Duration) {
	// Get main weather condition
	var mainWeather string
	var description string
	if len(weather.Weather) > 0 {
		mainWeather = weather.Weather[0].Main
		description = weather.Weather[0].Description
	} else {
		mainWeather = "Unknown"
		description = "unknown conditions"
	}

	// Determine units based on config
	windSpeedUnits := "m/s"
	tempUnit := "°C"

	switch config.Units {
	case "metric":
		windSpeedUnits = "m/s"
		tempUnit = "°C"
	case "imperial":
		windSpeedUnits = "mph"
		tempUnit = "°F"
	case "standard":
		windSpeedUnits = "m/s"
		tempUnit = "K"
	}

	// Format temperature and wind speed
	tempStr := fmt.Sprintf("Temperature: %.1f%s", weather.Main.Temp, tempUnit)
	windDir := getWindDirectionSymbol(weather.Wind.Deg)
	windSpeedStr := fmt.Sprintf("Wind: %.1f %s %s", weather.Wind.Speed, windSpeedUnits, windDir)
	humidityStr := fmt.Sprintf("Humidity: %d%%", weather.Main.Humidity)

	// Apply colors if enabled
	tempDisplay := tempStr
	windSpeedDisplay := windSpeedStr
	humidityDisplay := humidityStr

	if config.UseColors {
		tempDisplay = color.RedString(tempStr)
		windSpeedDisplay = color.GreenString(windSpeedStr)
		humidityDisplay = color.CyanString(humidityStr)
	}

	// Calculate sunrise and sunset times
	sunriseString := formatTime(weather.Sys.Sunrise, config, adjustment)
	sunsetString := formatTime(weather.Sys.Sunset, config, adjustment)

	// Format sunrise and sunset with colors if enabled
	sunriseDisplay := fmt.Sprintf("Sunrise: %s", sunriseString)
	sunsetDisplay := fmt.Sprintf("Sunset: %s", sunsetString)

	if config.UseColors {
		sunriseDisplay = color.YellowString(sunriseDisplay)
		sunsetDisplay = color.BlueString(sunsetDisplay)
	}

	// Date display
	dateLabel := ""
	dateValue := ""

	if config.ShowDate {
		dateLabel = "Date: "
		dateValue = time.Now().Format("01/02/06")
	}

	// Apply colors to date if enabled
	dateLabelDisplay := dateLabel
	dateValueDisplay := dateValue

	if config.UseColors {
		dateLabelDisplay = color.WhiteString(dateLabel)
		dateValueDisplay = color.WhiteString(dateValue)
	}

	// City name display
	cityDisplay := ""
	if config.ShowCityName {
		cityName := weather.Name
		if config.City != "" && cityName == "" {
			cityName = config.City
		}

		if config.UseColors {
			cityDisplay = color.GreenString(color.New(color.Bold).Sprintf("City: %s", cityName))
		} else {
			cityDisplay = fmt.Sprintf("City: %s", cityName)
		}
	}

	// Weather condition display
	var weatherDisplay string
	if config.UseColors {
		weatherDisplay = getColoredWeatherText(mainWeather, description)
	} else {
		weatherDisplay = fmt.Sprintf("Weather: %s", description)
	}

	// Display ascii art based on weather condition
	displayWeatherArt(mainWeather, cityDisplay, weatherDisplay, tempDisplay,
		windSpeedDisplay, humidityDisplay, sunriseDisplay,
		sunsetDisplay, dateLabelDisplay, dateValueDisplay)
}

// getColoredWeatherText returns a colored weather text based on the condition
func getColoredWeatherText(mainWeather, description string) string {
	switch mainWeather {
	case "Clear":
		return color.YellowString(color.New(color.Bold).Sprintf("Weather: %s", description))
	case "Clouds":
		return color.MagentaString(color.New(color.Bold).Sprintf("Weather: %s", description))
	case "Rain":
		return color.BlueString(color.New(color.Bold).Sprintf("Weather: %s", description))
	case "Snow":
		return color.CyanString(color.New(color.Bold).Sprintf("Weather: %s", description))
	case "Thunderstorm":
		return color.New(color.Bold).Sprintf("Weather: %s", description)
	default:
		return color.RedString(color.New(color.Bold).Sprintf("Weather: %s", description))
	}
}

// displayWeatherArt shows ASCII art based on the weather condition
func displayWeatherArt(mainWeather, cityDisplay, weatherDisplay, tempDisplay,
	windSpeedDisplay, humidityDisplay, sunriseDisplay,
	sunsetDisplay, dateLabelDisplay, dateValueDisplay string) {

	artTemplates := map[string]string{
		"Clear": `             %s
   \   /     %s
    .-.      %s
 ‒ (   ) ‒   %s
    ʻ-ʻ      %s
   /   \     %s
             %s
             %s%s`,
		"Clouds": `               %s
     .--.      %s
  .-(    ).    %s
 (___.__)__)   %s
               %s
               %s
               %s
               %s%s`,
		"Rain": `               %s
     .--.      %s
  .-(    ).    %s
 (___.__)__)   %s
  ʻ‚ʻ‚ʻ‚ʻ‚ʻ    %s
               %s
               %s
               %s%s`,
		"Snow": `               %s
     .--.      %s
  .-(    ).    %s
 (___.__)__)   %s
   * * * *     %s
  * * * *      %s
               %s
               %s%s`,
		"Thunderstorm": `               %s
     .--.      %s
  .-(    ).    %s
 (___.__)__)   %s
    /_  /_     %s
     /  /      %s
               %s
               %s%s`,
	}

	// Default art if weather condition not found
	art, exists := artTemplates[mainWeather]
	if !exists {
		art = artTemplates["Clouds"]
	}

	fmt.Println(fmt.Sprintf(art, cityDisplay, weatherDisplay, tempDisplay,
		windSpeedDisplay, humidityDisplay, sunriseDisplay, sunsetDisplay,
		dateLabelDisplay, dateValueDisplay))
}
