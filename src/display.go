package main

import (
	"fmt"
	"math"
	"strings"

	"github.com/fatih/color"
)

// getWindDirectionSymbol converts wind degrees to a direction symbol
func getWindDirectionSymbol(degrees int) string {
	directions := []string{"↑", "↗", "→", "↘", "↓", "↙", "←", "↖"}
	index := int((float64(degrees)+22.5)/45.0) % 8
	return directions[index]
}

// displayWeather renders the weather data with ASCII art
func displayWeather(weather *Weather, config Config) {
	// Get main weather condition
	var mainWeather string
	var description string
	var weatherID int
	if len(weather.Weather) > 0 {
		mainWeather = weather.Weather[0].Main
		description = weather.Weather[0].Description
		weatherID = weather.Weather[0].ID
	} else {
		mainWeather = "Unknown"
		description = "unknown conditions"
		weatherID = 0
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

	// Format precipitation info
	precipMM := 0.0
	if weather.Rain.OneHour > 0 {
		precipMM = weather.Rain.OneHour
	}

	popPercent := 0
	if weather.Pop > 0 {
		popPercent = int(math.Round(weather.Pop * 100))
	}

	// Format values differently for compact mode
	var tempStr, windStr, humidityStr, precipStr string
	if config.Compact {
		tempStr = fmt.Sprintf("%.1f%s", weather.Main.Temp, tempUnit)
		windStr = fmt.Sprintf("%.1f%s %s", weather.Wind.Speed, windSpeedUnits, getWindDirectionSymbol(weather.Wind.Deg))
		humidityStr = fmt.Sprintf("%d%%", weather.Main.Humidity)
		precipStr = fmt.Sprintf("%.1fmm | %d%%", precipMM, popPercent)
	} else {
		tempStr = fmt.Sprintf("Temperature: %.1f%s", weather.Main.Temp, tempUnit)
		windStr = fmt.Sprintf("Wind: %.1f %s %s", weather.Wind.Speed, windSpeedUnits, getWindDirectionSymbol(weather.Wind.Deg))
		humidityStr = fmt.Sprintf("Humidity: %d%%", weather.Main.Humidity)
		precipStr = fmt.Sprintf("Precip: %.1f mm | %d%%", precipMM, popPercent)
	}

	// Apply colors if enabled
	tempDisplay := tempStr
	windDisplay := windStr
	humidityDisplay := humidityStr
	precipDisplay := precipStr

	if config.UseColors {
		tempDisplay = color.RedString(tempStr)
		windDisplay = color.GreenString(windStr)
		humidityDisplay = color.CyanString(humidityStr)

		// Special handling for precipitation to color parts differently
		if config.Compact {
			parts := strings.Split(precipStr, "|")
			if len(parts) == 2 {
				precipDisplay = color.BlueString(parts[0]) + "|" + color.CyanString(parts[1])
			} else {
				precipDisplay = color.BlueString(precipStr)
			}
		} else {
			parts := strings.SplitN(precipStr, ":", 2)
			if len(parts) == 2 {
				precipDisplay = parts[0] + ":" + color.BlueString(parts[1])
			}
		}
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

		if config.Compact {
			cityDisplay = cityName // Just show the name without label
			if config.UseColors {
				cityDisplay = color.GreenString(color.New(color.Bold).Sprintf(cityName))
			}
		}
	}

	// Weather condition display
	var weatherDisplay string
	if config.Compact {
		weatherDisplay = description
	} else {
		weatherDisplay = fmt.Sprintf("Weather: %s", description)
	}

	if config.UseColors {
		weatherDisplay = getColoredWeatherText(mainWeather, description, config.Compact)
	}

	cityDisplay = strings.TrimSpace(cityDisplay)
	weatherDisplay = strings.TrimSpace(weatherDisplay)
	tempDisplay = strings.TrimSpace(tempDisplay)
	windDisplay = strings.TrimSpace(windDisplay)
	humidityDisplay = strings.TrimSpace(humidityDisplay)
	precipDisplay = strings.TrimSpace(precipDisplay)

	// Display the weather with icon
	displayWeatherArt(mainWeather, weatherID, cityDisplay, weatherDisplay,
		tempDisplay, windDisplay, humidityDisplay, precipDisplay, config)
}

// getColoredWeatherText returns a colored weather text based on the condition
func getColoredWeatherText(mainWeather, description string, compact bool) string {
	text := description
	if !compact {
		text = fmt.Sprintf("Weather: %s", description)
	}

	switch mainWeather {
	case "Clear":
		return color.YellowString(color.New(color.Bold).Sprintf(text))
	case "Clouds":
		return color.MagentaString(color.New(color.Bold).Sprintf(text))
	case "Rain":
		return color.BlueString(color.New(color.Bold).Sprintf(text))
	case "Snow":
		return color.CyanString(color.New(color.Bold).Sprintf(text))
	case "Thunderstorm":
		return color.New(color.Bold, color.BgRed).Sprintf(text)
	default:
		return color.RedString(color.New(color.Bold).Sprintf(text))
	}
}

// displayWeatherArt shows ASCII art based on the weather condition
func displayWeatherArt(mainWeather string, weatherID int, cityDisplay, weatherDisplay,
	tempDisplay, windDisplay, humidityDisplay, precipDisplay string, config Config) {

	// Get the weather icon
	iconLines := getWeatherIcon(mainWeather, weatherID, config.UseColors)

	// Prepare the text lines (remove date line)
	var textLines []string
	textLines = append(textLines, "") // Empty line to match icon top spacing

	if cityDisplay != "" {
		textLines = append(textLines, cityDisplay)
	}
	if weatherDisplay != "" {
		textLines = append(textLines, weatherDisplay)
	}
	textLines = append(textLines, tempDisplay, windDisplay, humidityDisplay)
	if precipDisplay != "" {
		textLines = append(textLines, precipDisplay)
	}

	textLines = append(textLines, "") // Empty line to match icon bottom spacing

	// Print icon and text lines together
	for i := 0; i < len(iconLines); i++ {
		iconLine := iconLines[i]
		textLine := ""
		if i < len(textLines) {
			textLine = textLines[i]
		}
		fmt.Printf("%s  %s\n", iconLine, textLine)
	}
}
