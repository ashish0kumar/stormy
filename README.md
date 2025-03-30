# Stormy

Neofetch-like, minimalistic, and customizable weather-fetching tool. 

> [!NOTE]
> Stormy's idea, structure and design is based off [Rainy](https://github.com/liveslol/rainy)

## Features

- Current weather conditions with ASCII art representation
- Temperature, wind, humidity, and sunrise/sunset times
- Customizable units (metric, imperial, standard)
- Local configuration file
- Weather data caching to reduce API calls
- Color support for terminals

## Installation

### Prerequisites

- Go 1.19 or higher
- An API key from [OpenWeatherMap](https://openweathermap.org/api)

### Build from Source

```bash
# Clone the repository
git clone https://github.com/ashish0kumar/stormy.git
cd stormy

# Build the application
go build -o stormy ./src

# Move to a directory in your PATH (optional)
sudo mv stormy /usr/local/bin/
```

## Configuration

Stormy will create a default configuration file on first run:

- Linux/macOS: `~/.config/stormy/stormy.toml`
- Windows: `%APPDATA%\stormy\stormy.toml`

You'll need to edit this file to add your API key and preferred city:

```toml
api_key = "api_key"      # Your OpenWeatherMap API Key
city = "New Delhi"       # City name
units = "metric"         # metric, imperial, or standard
timeplus = 0             # Hours to add to local time
timeminus = 0            # Hours to subtract from local time
showcityname = false     # Show city name in display
showdate = false         # Show date in display
timeformat = "24"        # 12 or 24 hour format
use_colors = false       # Use colors in terminal output
cache_duration = 30      # Cache duration in minutes
```

## Usage

```bash
# Basic usage
stormy

# Specify city via command line
stormy --city "New York"

# Use imperial units
stormy --units imperial

# Disable cache
stormy --no-cache

# Show help
stormy --help
```

## Examples


## Acknowledgements

- [OpenWeatherMap](https://openweathermap.org/) for providing weather data
- [rainy](https://github.com/liveslol/rainy) for the overall idea, structure and design of the project

## License

[MIT](LICENSE)