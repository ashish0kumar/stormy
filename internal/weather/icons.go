package weather

// getIcon returns the ASCII art for a given weather condition with proper spacing
func getIcon(name string, useColors bool) []string {
	// Base monochrome icons with consistent 7-line height (top/bottom padding)
	icon := map[string][]string{
		"Unknown": {
			"             ",
			"    .-.      ",
			"     __)     ",
			"    (        ",
			"     `-’     ",
			"      •      ",
			"             ",
		},
		"Sunny": {
			"             ",
			"    \\   /    ",
			"     .-.     ",
			"  ― (   ) ―  ",
			"     `-’     ",
			"    /   \\    ",
			"             ",
		},
		"PartlyCloudy": {
			"             ",
			"   \\  /      ",
			" _ /\"\".-.    ",
			"   \\_(   ).  ",
			"   /(___(__) ",
			"             ",
			"             ",
		},
		"Cloudy": {
			"             ",
			"             ",
			"     .--.    ",
			"  .-(    ).  ",
			" (___.__)__) ",
			"             ",
			"             ",
		},
		"VeryCloudy": {
			"             ",
			"             ",
			"     .--.    ",
			"  .-(    ).  ",
			" (___.__)__) ",
			"             ",
			"             ",
		},
		"LightShowers": {
			"             ",
			" _`/\"\".-.    ",
			"  ,\\_(   ).  ",
			"   /(___(__) ",
			"     ' ' ' ' ",
			"    ' ' ' '  ",
			"             ",
		},
		"HeavyShowers": {
			"             ",
			" _`/\"\".-.    ",
			"  ,\\_(   ).  ",
			"   /(___(__) ",
			"   ‚'‚'‚'‚'  ",
			"   ‚'‚'‚'‚'  ",
			"             ",
		},
		"LightSnow": {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    *  *  *  ",
			"   *  *  *   ",
			"             ",
		},
		"HeavySnow": {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"   * * * *   ",
			"  * * * *    ",
			"             ",
		},
		"Thunderstorm": {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    ⚡\"\"⚡\"\" ",
			"  ‚'‚'‚'‚'   ",
			"             ",
		},
		"Fog": {
			"             ",
			"             ",
			" _ - _ - _ - ",
			"  _ - _ - _  ",
			" _ - _ - _ - ",
			"             ",
			"             ",
		},
	}

	if !useColors {
		return icon[name]
	}

	// Colored versions with same spacing
	coloredIcon := map[string][]string{
		"Sunny": {
			"             ",
			"\033[38;5;226m    \\   /    \033[0m",
			"\033[38;5;226m     .-.     \033[0m",
			"\033[38;5;226m  ― (   ) ―  \033[0m",
			"\033[38;5;226m     `-’     \033[0m",
			"\033[38;5;226m    /   \\    \033[0m",
			"             ",
		},
		"PartlyCloudy": {
			"             ",
			"\033[38;5;226m   \\  /\033[0m      ",
			"\033[38;5;226m _ /\"\"\033[38;5;250m.-.    \033[0m",
			"\033[38;5;226m   \\_\033[38;5;250m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
			"             ",
			"             ",
		},
		"Cloudy": {
			"             ",
			"             ",
			"\033[38;5;250m     .--.    \033[0m",
			"\033[38;5;250m  .-(    ).  \033[0m",
			"\033[38;5;250m (___.__)__) \033[0m",
			"             ",
			"             ",
		},
		"VeryCloudy": {
			"             ",
			"             ",
			"\033[38;5;240;1m     .--.    \033[0m",
			"\033[38;5;240;1m  .-(    ).  \033[0m",
			"\033[38;5;240;1m (___.__)__) \033[0m",
			"             ",
			"             ",
		},
		"LightShowers": {
			"             ",
			"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
			"\033[38;5;111m     ' ' ' ' \033[0m",
			"\033[38;5;111m    ' ' ' '  \033[0m",
			"             ",
		},
		"HeavyShowers": {
			"             ",
			"\033[38;5;226m _`/\"\"\033[38;5;240;1m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;240;1m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;240;1m(___(__) \033[0m",
			"\033[38;5;21;1m   ‚'‚'‚'‚'  \033[0m",
			"\033[38;5;21;1m   ‚'‚'‚'‚'  \033[0m",
			"             ",
		},
		"LightSnow": {
			"             ",
			"\033[38;5;250m     .-.     \033[0m",
			"\033[38;5;250m    (   ).   \033[0m",
			"\033[38;5;250m   (___(__)  \033[0m",
			"\033[38;5;255m    *  *  *  \033[0m",
			"\033[38;5;255m   *  *  *   \033[0m",
			"             ",
		},
		"HeavySnow": {
			"             ",
			"\033[38;5;240;1m     .-.     \033[0m",
			"\033[38;5;240;1m    (   ).   \033[0m",
			"\033[38;5;240;1m   (___(__)  \033[0m",
			"\033[38;5;255;1m   * * * *   \033[0m",
			"\033[38;5;255;1m  * * * *    \033[0m",
			"             ",
		},
		"Thunderstorm": {
			"             ",
			"\033[38;5;240;1m     .-.     \033[0m",
			"\033[38;5;240;1m    (   ).   \033[0m",
			"\033[38;5;240;1m   (___(__)  \033[0m",
			"\033[38;5;228;5m    ⚡\033[38;5;111;25m\"\"\033[38;5;228;5m⚡\033[38;5;111;25m\"\" \033[0m",
			"\033[38;5;21;1m  ‚'‚'‚'‚'   \033[0m",
			"             ",
		},
		"Fog": {
			"             ",
			"             ",
			"\033[38;5;251m _ - _ - _ - \033[0m",
			"\033[38;5;251m  _ - _ - _  \033[0m",
			"\033[38;5;251m _ - _ - _ - \033[0m",
			"             ",
			"             ",
		},
	}

	if ci, ok := coloredIcon[name]; ok {
		return ci
	}
	return icon[name]
}

// getWeatherIcon determines which icon to use based on weather condition
func getWeatherIcon(weatherMain string, weatherID int, useColors bool) []string {
	iconMap := map[int]string{
		// Thunderstorm
		200: "Thunderstorm",
		201: "Thunderstorm",
		202: "Thunderstorm",
		210: "Thunderstorm",
		211: "Thunderstorm",
		212: "Thunderstorm",
		221: "Thunderstorm",
		230: "Thunderstorm",
		231: "Thunderstorm",
		232: "Thunderstorm",

		// Drizzle
		300: "LightShowers",
		301: "LightShowers",
		302: "LightShowers",
		310: "LightShowers",
		311: "LightShowers",
		312: "LightShowers",
		313: "LightShowers",
		314: "LightShowers",
		321: "LightShowers",

		// Rain
		500: "LightShowers",
		501: "LightShowers",
		502: "HeavyShowers",
		503: "HeavyShowers",
		504: "HeavyShowers",
		511: "LightSnow",
		520: "LightShowers",
		521: "LightShowers",
		522: "HeavyShowers",
		531: "HeavyShowers",

		// Snow
		600: "LightSnow",
		601: "HeavySnow",
		602: "HeavySnow",
		611: "LightSnow",
		612: "LightSnow",
		613: "LightSnow",
		615: "LightSnow",
		616: "LightSnow",
		620: "LightSnow",
		621: "HeavySnow",
		622: "HeavySnow",

		// Atmosphere
		701: "Fog",
		711: "Fog",
		721: "Fog",
		731: "Fog",
		741: "Fog",
		751: "Fog",
		761: "Fog",
		762: "Fog",
		771: "Fog",
		781: "Fog",

		// Clear
		800: "Sunny",

		// Clouds
		801: "PartlyCloudy",
		802: "Cloudy",
		803: "VeryCloudy",
		804: "VeryCloudy",
	}

	iconName := "Unknown"
	if name, ok := iconMap[weatherID]; ok {
		iconName = name
	} else {
		switch weatherMain {
		case "Clear":
			iconName = "Sunny"
		case "Clouds":
			iconName = "Cloudy"
		case "Rain":
			iconName = "LightShowers"
		case "Drizzle":
			iconName = "LightShowers"
		case "Thunderstorm":
			iconName = "Thunderstorm"
		case "Snow":
			iconName = "LightSnow"
		case "Mist", "Smoke", "Haze", "Dust", "Fog", "Sand", "Ash", "Squall", "Tornado":
			iconName = "Fog"
		}
	}

	return getIcon(iconName, useColors)
}
