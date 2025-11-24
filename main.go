package main

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ashish0kumar/stormy/internal/weather"
	"github.com/fatih/color"
	"github.com/k0kubun/go-ansi"
	"golang.org/x/term"
)

// version is set during build time using -ldflags
var version = "dev"

func init() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM, syscall.SIGQUIT)
	go func() {
		stop := <-c
		// reset cursor visibility for live mode
		_, _ = ansi.Print("\x1b[?25h")
		switch stop {
		case os.Interrupt, syscall.SIGTERM:
			_, _ = fmt.Fprintf(
				os.Stderr, "\n\n[%s] Program interrupted. Bye!\n", color.New(color.FgRed).SprintFunc()("x"),
			)
			os.Exit(1)
		case syscall.SIGQUIT:
			fmt.Printf("\n\n[%s] Stopping now, bye!\n", color.New(color.FgGreen).SprintFunc()("✓"))
			os.Exit(0)
		}
	}()
}

func listenForQuit(stop chan struct{}) {
	var shouldExit, shouldInterrupt bool
	// Switch stdin into 'raw' mode
	oldState, errRaw := term.MakeRaw(int(os.Stdin.Fd()))
	if errRaw != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Error setting raw mode: %v\n", errRaw)
		return
	}
	defer func(fd int, state *term.State) {
		_ = term.Restore(fd, state)
		if shouldExit {
			_, _ = ansi.Print("\x1b[?25h") // restore cursor
			fmt.Printf("\n\n[%s] Stopping now, bye!\n", color.New(color.FgGreen).SprintFunc()("✓"))
			os.Exit(0)
		}
		if shouldInterrupt {
			_, _ = ansi.Print("\x1b[?25h") // restore cursor
			_, _ = fmt.Fprintf(
				os.Stderr, "\n\n[%s] Program interrupted. Bye!\n", color.New(color.FgRed).SprintFunc()("x"),
			)
			os.Exit(1)
		}
	}(int(os.Stdin.Fd()), oldState)

	buffer := make([]byte, 1)
	for {
		select {
		case <-stop:
			return
		default:
			n, err := os.Stdin.Read(buffer)
			if err != nil || n == 0 {
				continue
			}

			char := buffer[0]
			if char == 'q' || char == 'Q' {
				shouldExit = true
				return
			}
			if char == 3 { // Ctrl+C
				shouldInterrupt = true
				return
			}
		}
	}
}

func main() {
	// Parse command line flags
	flags := weather.ParseFlags()

	// Handle version flag
	if flags.Version {
		fmt.Printf("stormy version %s\n", version)
		os.Exit(0)
	}

	// Read/create config
	config := weather.ReadConfig()

	// Override config with command line flags if provided
	preFlagsConfig := config
	weather.ApplyFlags(&config, flags)

	scanner := bufio.NewScanner(os.Stdin)

	// Check if the city is set
	if config.City == "" {
		fmt.Printf("No city found in your configuration, please enter the city to check the weather for: ")
		scanner.Scan()
		newCity := scanner.Text()
		config.City = newCity
		preFlagsConfig.City = newCity
		err := weather.WriteConfig(preFlagsConfig, weather.GetConfigPath())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to update your stored config: %v\n", err)
			_, _ = fmt.Fprintln(os.Stderr, "You'll need to enter your city back next time if still not provided through the configuration or command line flags")
			_, _ = fmt.Fprintln(os.Stderr, "Proceeding with the provided city")
		}
	}

	// Check if the API key and city are set
	if config.Provider == weather.ProviderOpenWeatherMap && config.ApiKey == "" {
		fmt.Printf("No API key provided for %s, please enter it: ", config.Provider)
		scanner.Scan()
		apiKey := scanner.Text()
		config.ApiKey = apiKey
		preFlagsConfig.ApiKey = apiKey
		err := weather.WriteConfig(preFlagsConfig, weather.GetConfigPath())
		if err != nil {
			_, _ = fmt.Fprintf(os.Stderr, "Failed to update your stored config: %v\n", err)
			_, _ = fmt.Fprintln(os.Stderr, "You'll need to enter your API key back next time if still not provided through the configuration or command line flags")
			_, _ = fmt.Fprintln(os.Stderr, "Proceeding with the provided API key")
		}
	}

	fetchAndDisplay(config, false)
}

// fetchAndDisplay fetches weather data and displays it according to the given configuration.
// clearDisplay determines whether the screen should be cleared before displaying updated information.
func fetchAndDisplay(config weather.Config, clearDisplay bool) {
	// Fetch weather data
	weatherData, err := weather.FetchWeather(config)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "Failed to fetch weather data: %v\n", err)
		if errors.Is(err, weather.ErrUnsupportedQuery) {
			_, _ = fmt.Fprintf(
				os.Stderr, "Detailed queries are not supported by %s, try using other providers.\n",
				weather.ProviderOpenMeteo,
			)
		} else {
			_, _ = fmt.Fprintln(os.Stderr, "Please check your internet connection and API key.")
		}
		os.Exit(1)
	}

	// Clear screen in live mode
	if clearDisplay {
		_, _ = ansi.Printf("\x1b[%dA\x1b[J", 7) // maximum number of displayed lines
	}

	// Display the weather
	weather.DisplayWeather(weatherData, config)

	// Loop in live mode
	if !config.LiveMode {
		return
	}
	if !clearDisplay {
		// hide cursor on live mode startup
		_, _ = ansi.Print("\x1b[?25l")
	}
	// handle q press
	stop := make(chan struct{})
	go listenForQuit(stop)
	time.Sleep(15 * time.Second)
	stop <- struct{}{}
	fetchAndDisplay(config, true)
}
