# stormy

Neofetch-like, minimalistic, and customizable weather-fetching CLI based on
[rainy](https://github.com/liveslol/rainy)

<img src="./assets/screenshot.png" width="70%">

## Features

- Current weather conditions with ASCII art representation
- Temperature, wind, humidity, and precipitation information
- Customizable units (metric, imperial, standard)
- Local configuration file
- Color support for terminals
- Compact display mode

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

`stormy` will create a default configuration file on first run:

- Linux/macOS: `~/.config/stormy/stormy.toml`
- Windows: `%APPDATA%\stormy\stormy.toml`

### Configuration Options

- `api_key`: Your OpenWeatherMap API key.
- `city`: The city for which to fetch weather data.
- `units`: Units for temperature and wind speed (`metric`, `imperial` or
  `standard`).
- `showcityname`: Whether to display the city name (`true` or `false`).
- `use_colors`: Enables and disables text colors (`true` or `false`).
- `compact`: Use a more compact display format (`true` or `false`).

### Example Config

```toml
api_key = "your_api_key"
city = "New Delhi"
units = "metric"
showcityname = false
use_colors = false
compact = false
```

## Usage

```bash
# Basic usage
stormy

# Specify city via command line
stormy --city "New York"

# Use imperial units
stormy --units imperial

# Use compact display mode
stormy --compact

# Show help
stormy --help
```

## Display Examples

| ![basic](./assets/basic.png)     | ![color](./assets/color.png) |
| -------------------------------- | ---------------------------- |
| ![compact](./assets/compact.png) | ![units](./assets/units.png) |
| ![](./assets/few_clouds.png)     | ![](./assets/light_rain.png) |

## Weather Icons

The application displays different ASCII art based on weather conditions:

```sh
Clear/Sunny             Partly Cloudy
             
    \   /                 \  /      
     .-.                _ /"".-.    
  ― (   ) ―               \_(   ).  
     `-’                  /(___(__) 
    /   \                       
                             

Cloudy                  Very Cloudy
             
     .--.                    .--.    
  .-(    ).               .-(    ).  
 (___.__)__)             (___.__)__) 
                             

Light Showers           Heavy Showers
             
 _`/"".-.              _`/"".-.    
  ,\_(   ).             ,\_(   ).  
   /(___(__)             /(___(__) 
     ' ' ' '              ‚'‚'‚'‚'  
    ' ' ' '              ‚'‚'‚'‚'  
                             

Light Snow              Heavy Snow
             
     .-.                  .-.     
    (   ).               (   ).   
   (___(__)             (___(__)  
    *  *  *             * * * *   
   *  *  *             * * * *    
                             

Thunderstorm              Fog
             
     .-.                
    (   ).               
   (___(__)             _ - _ - _ - 
   ⚡""⚡""              _ - _ - _  
  ‚'‚'‚'‚'             _ - _ - _ -
```

Colors are applied when the `use_colors` option is enabled.

## Acknowledgements

- [OpenWeatherMap](https://openweathermap.org/) for providing weather data
- [rainy](https://github.com/liveslol/rainy) for the overall idea, structure and
  design of the project
- [wttr.in](https://github.com/chubin/wttr.in?tab=readme-ov-file) for the ASCII
  weather icons

## License

[MIT](LICENSE)
