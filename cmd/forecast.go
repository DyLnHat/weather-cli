package cmd

import (
	"fmt"
	"strings"

	"github.com/DyLnHat/weather-cli/internal/api"
	"github.com/DyLnHat/weather-cli/internal/config"
	"github.com/DyLnHat/weather-cli/internal/display"
	"github.com/spf13/cobra"
)

var days int

var forecastCmd = &cobra.Command{
	Use:   "forecast [city]",
	Short: "Show a weather forecast for a city",
	Example: `  weather forecast Madrid
  		weather forecast "New York" --days 3
  		weather forecast Tokyo --days 5`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		city := strings.Join(args, " ")

		// Validate --days range
		if days < 1 || days > 5 {
			return fmt.Errorf("--days must be between 1 and 5")
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("could not load config: %w", err)
		}

		if cfg.APIKey == "" {
			display.PrintError("API key not set — run: weather config --set-key <YOUR_KEY>\nGet a free key at: https://openweathermap.org/api")
			return nil
		}

		client := api.NewClient(cfg.APIKey, cfg.Units)
		forecast, err := client.GetForecast(city, days)
		if err != nil {
			display.PrintError(err.Error())
			return nil
		}

		display.PrintForecast(forecast, cfg.Units)
		return nil
	},
}

func init() {
	forecastCmd.Flags().IntVarP(&days, "days", "d", 5, "Number of forecast days (1-5)")
}
