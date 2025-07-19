package weather

import (
	"fmt"
	"math"
	"strings"

	"github.com/fatih/color"
)

const (
	MpsToMph  = 2.23694
	KmphToMph = 0.621371
	MpsToKmph = 3.6
)

func celsiusToFahrenheit(c float64) float64 {
	return c*9/5 + 32
}

// getWindDirectionSymbol converts wind degrees to a direction symbol
func getWindDirectionSymbol(degrees int) string {
	directions := []string{"↑", "↗", "→", "↘", "↓", "↙", "←", "↖"}
	index := int((float64(degrees)+22.5)/45.0) % 8
	return directions[index]
}

// DisplayWeather renders the weather data with ASCII art
func DisplayWeather(weather *Weather, config Config) {
	// Get main weather condition
	var mainWeather, description string
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
	var windSpeedUnits, tempUnit string

	// Convert wind speed based on provider and units
	windSpeed := weather.Wind.Speed
	temperature := weather.Main.Temp

	switch config.Units {
	case "imperial":
		windSpeedUnits = "mph"
		tempUnit = "°F"

		// Convert temperature for both providers
		temperature = celsiusToFahrenheit(temperature)

		// Convert wind speed based on provider
		if config.Provider == ProviderOpenWeatherMap {
			windSpeed *= MpsToMph // m/s to mph
		} else {
			windSpeed *= KmphToMph // km/h to mph
		}

	default:
		windSpeedUnits = "km/h"
		tempUnit = "°C"

		// Convert wind speed for OpenWeatherMap provider
		if config.Provider == ProviderOpenWeatherMap {
			windSpeed *= MpsToKmph // m/s to km/h
		}
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

	labels := make([]string, 0, 6)
	values := make([]string, 0, cap(labels))

	// City name display
	cityName := weather.Name
	if config.City != "" && cityName == "" {
		cityName = config.City
	}

	if config.ShowCityName {
		if !config.Compact {
			labels = append(labels, "City")
			values = append(values, cityName)
		} else {
			// In compact mode, we'll just show the city name in the display function
		}
	}

	// Weather info
	if !config.Compact {
		labels = append(labels, "Weather ")
		values = append(values, description)

		labels = append(labels, "Temp ")
		values = append(values, fmt.Sprintf("%.1f%s", temperature, tempUnit))

		labels = append(labels, "Wind ")
		values = append(
			values, fmt.Sprintf("%.1f %s %s", windSpeed, windSpeedUnits, getWindDirectionSymbol(weather.Wind.Deg)),
		)

		labels = append(labels, "Humidity ")
		values = append(values, fmt.Sprintf("%d%%", weather.Main.Humidity))

		labels = append(labels, "Precip ")
		values = append(values, fmt.Sprintf("%.1f mm | %d%%", precipMM, popPercent))
	} else {
		// Compact mode doesn't use labels in the same way
		weatherDisplay := description
		tempDisplay := fmt.Sprintf("%.1f%s", temperature, tempUnit)
		windDisplay := fmt.Sprintf("%.1f%s %s", windSpeed, windSpeedUnits, getWindDirectionSymbol(weather.Wind.Deg))
		humidityDisplay := fmt.Sprintf("%d%%", weather.Main.Humidity)
		precipDisplay := fmt.Sprintf("%.1fmm | %d%%", precipMM, popPercent)

		// For compact mode, we'll just pass these values directly to the display function
		displayWeatherArtCompact(
			mainWeather,
			weatherID,
			cityName,
			weatherDisplay,
			tempDisplay,
			windDisplay,
			humidityDisplay,
			precipDisplay,
			config,
		)
		return
	}

	// For standard mode, display with aligned labels and values
	displayWeatherArtAligned(mainWeather, weatherID, labels, values, config)
}

// getColoredWeatherText returns a colored weather text based on the condition
func getColoredWeatherText(mainWeather, description string) string {
	text := description

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

// displayWeatherArtAligned shows ASCII art with vertically aligned labels and values
func displayWeatherArtAligned(mainWeather string, weatherID int, labels, values []string, config Config) {
	// Get the weather icon
	iconLines := getWeatherIcon(mainWeather, weatherID, config.UseColors)

	// Find the maximum label length for alignment
	maxLabelLen := 0
	for _, label := range labels {
		if len(label) > maxLabelLen {
			maxLabelLen = len(label)
		}
	}

	// Apply colors to values if enabled
	coloredValues := make([]string, len(values))
	coloredLabels := make([]string, len(labels))
	for i, value := range values {
		if config.UseColors {
			coloredLabels[i] = color.BlueString(labels[i])
		} else {
			coloredLabels[i] = labels[i]
		}

		if config.UseColors {
			switch i {
			case 0: // City
				coloredValues[i] = color.GreenString(color.New(color.Bold).Sprintf(value))
			case 1: // Weather description
				coloredValues[i] = getColoredWeatherText(mainWeather, value)
			case 2: // Temperature
				coloredValues[i] = color.RedString(value)
			case 3: // Wind
				coloredValues[i] = color.GreenString(value)
			case 4: // Humidity
				coloredValues[i] = color.CyanString(value)
			case 5: // Precipitation
				parts := strings.Split(value, "|")
				if len(parts) == 2 {
					coloredValues[i] = color.BlueString(strings.TrimSpace(parts[0])) + " | " + color.CyanString(strings.TrimSpace(parts[1]))
				} else {
					coloredValues[i] = color.BlueString(value)
				}
			default:
				coloredValues[i] = value
			}
		} else {
			coloredValues[i] = value
		}
	}

	// Prepare the text lines
	textLines := make([]string, 0, len(labels)+2)
	textLines = append(textLines, "") // Empty line to match icon top spacing

	// Add formatted lines with aligned labels and values
	for i := 0; i < len(labels); i++ {
		if config.UseColors {
			formattedLine := fmt.Sprintf("%-*s %s", maxLabelLen, labels[i], values[i])
			coloredLabel := color.BlueString(labels[i])
			formattedLine = strings.Replace(formattedLine, labels[i], coloredLabel, 1)
			formattedLine = strings.Replace(formattedLine, values[i], coloredValues[i], 1)
			textLines = append(textLines, formattedLine)
		} else {
			textLines = append(textLines, fmt.Sprintf("%-*s %s", maxLabelLen, labels[i], values[i]))
		}
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

// displayWeatherArtCompact shows ASCII art with compact formatting
func displayWeatherArtCompact(
	mainWeather string, weatherID int, cityName, weatherDisplay,
	tempDisplay, windDisplay, humidityDisplay, precipDisplay string, config Config,
) {

	// Get the weather icon
	iconLines := getWeatherIcon(mainWeather, weatherID, config.UseColors)

	// Apply colors if enabled
	if config.UseColors {
		if cityName != "" && config.ShowCityName {
			cityName = color.GreenString(color.New(color.Bold).Sprintf(cityName))
		}
		weatherDisplay = getColoredWeatherText(mainWeather, weatherDisplay)
		tempDisplay = color.RedString(tempDisplay)
		windDisplay = color.GreenString(windDisplay)
		humidityDisplay = color.CyanString(humidityDisplay)

		parts := strings.Split(precipDisplay, "|")
		if len(parts) == 2 {
			precipDisplay = color.BlueString(strings.TrimSpace(parts[0])) + " | " + color.CyanString(strings.TrimSpace(parts[1]))
		} else {
			precipDisplay = color.BlueString(precipDisplay)
		}
	}

	// Prepare the text lines
	textLines := make([]string, 0, 5)
	textLines = append(textLines, "") // Empty line to match icon top spacing

	if cityName != "" && config.ShowCityName {
		textLines = append(textLines, cityName)
	}

	textLines = append(textLines, weatherDisplay, tempDisplay, windDisplay, humidityDisplay)

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
