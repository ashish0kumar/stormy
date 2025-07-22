package weather

var (
	// icon base monochrome icons with consistent 7-line height (top/bottom padding)
	icon = map[string][]string{
		ConditionUnknown: {
			"             ",
			"    .-.      ",
			"     __)     ",
			"    (        ",
			"     `-’     ",
			"      •      ",
			"             ",
		},
		ConditionSunny: {
			"             ",
			"    \\   /    ",
			"     .-.     ",
			"  ― (   ) ―  ",
			"     `-’     ",
			"    /   \\    ",
			"             ",
		},
		ConditionPartlyCloudy: {
			"             ",
			"   \\  /      ",
			" _ /\"\".-.    ",
			"   \\_(   ).  ",
			"   /(___(__) ",
			"             ",
			"             ",
		},
		ConditionCloudy: {
			"             ",
			"             ",
			"     .--.    ",
			"  .-(    ).  ",
			" (___.__)__) ",
			"             ",
			"             ",
		},
		ConditionVeryCloudy: {
			"             ",
			"             ",
			"     .--.    ",
			"  .-(    ).  ",
			" (___.__)__) ",
			"             ",
			"             ",
		},
		ConditionLightShowers: {
			"             ",
			" _`/\"\".-.    ",
			"  ,\\_(   ).  ",
			"   /(___(__) ",
			"     ' ' ' ' ",
			"    ' ' ' '  ",
			"             ",
		},
		ConditionHeavyShowers: {
			"             ",
			" _`/\"\".-.    ",
			"  ,\\_(   ).  ",
			"   /(___(__) ",
			"   ‚'‚'‚'‚'  ",
			"   ‚'‚'‚'‚'  ",
			"             ",
		},
		ConditionLightSnow: {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    *  *  *  ",
			"   *  *  *   ",
			"             ",
		},
		ConditionHeavySnow: {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"   * * * *   ",
			"  * * * *    ",
			"             ",
		},
		ConditionThunderstorm: {
			"             ",
			"     .-.     ",
			"    (   ).   ",
			"   (___(__)  ",
			"    ⚡\"\"⚡\"\" ",
			"  ‚'‚'‚'‚'   ",
			"             ",
		},
		ConditionFog: {
			"             ",
			"             ",
			" _ - _ - _ - ",
			"  _ - _ - _  ",
			" _ - _ - _ - ",
			"             ",
			"             ",
		},
	}

	// coloredIcon colored icons with the same spacing
	coloredIcon = map[string][]string{
		ConditionSunny: {
			"             ",
			"\033[38;5;226m    \\   /    \033[0m",
			"\033[38;5;226m     .-.     \033[0m",
			"\033[38;5;226m  ― (   ) ―  \033[0m",
			"\033[38;5;226m     `-’     \033[0m",
			"\033[38;5;226m    /   \\    \033[0m",
			"             ",
		},
		ConditionPartlyCloudy: {
			"             ",
			"\033[38;5;226m   \\  /\033[0m      ",
			"\033[38;5;226m _ /\"\"\033[38;5;250m.-.    \033[0m",
			"\033[38;5;226m   \\_\033[38;5;250m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
			"             ",
			"             ",
		},
		ConditionCloudy: {
			"             ",
			"             ",
			"\033[38;5;250m     .--.    \033[0m",
			"\033[38;5;250m  .-(    ).  \033[0m",
			"\033[38;5;250m (___.__)__) \033[0m",
			"             ",
			"             ",
		},
		ConditionVeryCloudy: {
			"             ",
			"             ",
			"\033[38;5;240;1m     .--.    \033[0m",
			"\033[38;5;240;1m  .-(    ).  \033[0m",
			"\033[38;5;240;1m (___.__)__) \033[0m",
			"             ",
			"             ",
		},
		ConditionLightShowers: {
			"             ",
			"\033[38;5;226m _`/\"\"\033[38;5;250m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;250m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;250m(___(__) \033[0m",
			"\033[38;5;111m     ' ' ' ' \033[0m",
			"\033[38;5;111m    ' ' ' '  \033[0m",
			"             ",
		},
		ConditionHeavyShowers: {
			"             ",
			"\033[38;5;226m _`/\"\"\033[38;5;240;1m.-.    \033[0m",
			"\033[38;5;226m  ,\\_\033[38;5;240;1m(   ).  \033[0m",
			"\033[38;5;226m   /\033[38;5;240;1m(___(__) \033[0m",
			"\033[38;5;21;1m   ‚'‚'‚'‚'  \033[0m",
			"\033[38;5;21;1m   ‚'‚'‚'‚'  \033[0m",
			"             ",
		},
		ConditionLightSnow: {
			"             ",
			"\033[38;5;250m     .-.     \033[0m",
			"\033[38;5;250m    (   ).   \033[0m",
			"\033[38;5;250m   (___(__)  \033[0m",
			"\033[38;5;255m    *  *  *  \033[0m",
			"\033[38;5;255m   *  *  *   \033[0m",
			"             ",
		},
		ConditionHeavySnow: {
			"             ",
			"\033[38;5;240;1m     .-.     \033[0m",
			"\033[38;5;240;1m    (   ).   \033[0m",
			"\033[38;5;240;1m   (___(__)  \033[0m",
			"\033[38;5;255;1m   * * * *   \033[0m",
			"\033[38;5;255;1m  * * * *    \033[0m",
			"             ",
		},
		ConditionThunderstorm: {
			"             ",
			"\033[38;5;240;1m     .-.     \033[0m",
			"\033[38;5;240;1m    (   ).   \033[0m",
			"\033[38;5;240;1m   (___(__)  \033[0m",
			"\033[38;5;228;5m    ⚡\033[38;5;111;25m\"\"\033[38;5;228;5m⚡\033[38;5;111;25m\"\" \033[0m",
			"\033[38;5;21;1m  ‚'‚'‚'‚'   \033[0m",
			"             ",
		},
		ConditionFog: {
			"             ",
			"             ",
			"\033[38;5;251m _ - _ - _ - \033[0m",
			"\033[38;5;251m  _ - _ - _  \033[0m",
			"\033[38;5;251m _ - _ - _ - \033[0m",
			"             ",
			"             ",
		},
	}
)

// getIcon returns the ASCII art for a given weather condition with proper spacing
func getIcon(name string, useColors bool) []string {
	if !useColors {
		return icon[name]
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
		200: ConditionThunderstorm,
		201: ConditionThunderstorm,
		202: ConditionThunderstorm,
		210: ConditionThunderstorm,
		211: ConditionThunderstorm,
		212: ConditionThunderstorm,
		221: ConditionThunderstorm,
		230: ConditionThunderstorm,
		231: ConditionThunderstorm,
		232: ConditionThunderstorm,

		// Drizzle
		300: ConditionLightShowers,
		301: ConditionLightShowers,
		302: ConditionLightShowers,
		310: ConditionLightShowers,
		311: ConditionLightShowers,
		312: ConditionLightShowers,
		313: ConditionLightShowers,
		314: ConditionLightShowers,
		321: ConditionLightShowers,

		// Rain
		500: ConditionLightShowers,
		501: ConditionLightShowers,
		502: ConditionHeavyShowers,
		503: ConditionHeavyShowers,
		504: ConditionHeavyShowers,
		511: ConditionLightSnow,
		520: ConditionLightShowers,
		521: ConditionLightShowers,
		522: ConditionHeavyShowers,
		531: ConditionHeavyShowers,

		// Snow
		600: ConditionLightSnow,
		601: ConditionHeavySnow,
		602: ConditionHeavySnow,
		611: ConditionLightSnow,
		612: ConditionLightSnow,
		613: ConditionLightSnow,
		615: ConditionLightSnow,
		616: ConditionLightSnow,
		620: ConditionLightSnow,
		621: ConditionHeavySnow,
		622: ConditionHeavySnow,

		// Atmosphere
		701: ConditionFog,
		711: ConditionFog,
		721: ConditionFog,
		731: ConditionFog,
		741: ConditionFog,
		751: ConditionFog,
		761: ConditionFog,
		762: ConditionFog,
		771: ConditionFog,
		781: ConditionFog,

		// Clear
		800: ConditionSunny,

		// Clouds
		801: ConditionPartlyCloudy,
		802: ConditionCloudy,
		803: ConditionVeryCloudy,
		804: ConditionVeryCloudy,
	}

	iconName := ConditionUnknown
	if name, ok := iconMap[weatherID]; ok {
		iconName = name
	} else {
		switch weatherMain {
		case ConditionClear:
			iconName = ConditionSunny
		case ConditionClouds:
			iconName = ConditionCloudy
		case ConditionRain:
			iconName = ConditionLightShowers
		case ConditionDrizzle:
			iconName = ConditionLightShowers
		case ConditionThunderstorm:
			iconName = ConditionThunderstorm
		case ConditionSnow:
			iconName = ConditionLightSnow
		case ConditionMist, ConditionSmoke, ConditionHaze, ConditionDust, ConditionFog, ConditionSand, ConditionAsh, ConditionSquall, ConditionTornado:
			iconName = ConditionFog
		}
	}

	return getIcon(iconName, useColors)
}
