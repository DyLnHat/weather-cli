package cmd

import (
	"fmt"
	"strings"

	"github.com/DyLnHat/weather-cli/internal/api"
	"github.com/DyLnHat/weather-cli/internal/config"
	"github.com/DyLnHat/weather-cli/internal/display"
	"github.com/spf13/cobra"
)

var currentCmd = &cobra.Command{
	Use:   "current [city]",
	Short: "Show current weather for a city",
	Example: `  weather current Madrid
  		weather current "New York"
  		weather current London`,
	Args: cobra.MinimumNArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		city := strings.Join(args, " ")

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("could not load config: %w", err)
		}

		// No API key set
		if cfg.APIKey == "" {
			display.PrintError("API key not set — run: weather config --set-key <YOUR_KEY>\nGet a free key at: https://openweathermap.org/api")
			return nil
		}

		client := api.NewClient(cfg.APIKey, cfg.Units)
		weather, err := client.GetCurrent(city)
		if err != nil {
			display.PrintError(err.Error())
			return nil
		}

		display.PrintCurrent(weather, cfg.Units)
		return nil
	},
}
