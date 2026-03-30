package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "weather",
	Short: "Get weather forecasts from your terminal",
	Long: `weather-cli - Current conditions and forecasts for any city.
	
	Examples:
		weather current Madrid
		weather current "New York"
		weather forecast Tokyo --days 3
		weather config --set-key YOUR_API_KEY
		weather config --set-units imperial
		weather config --show`,
	SilenceUsage:  true,
	SilenceErrors: true,
}

// Entry point called from main.go
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.AddCommand(currentCmd)
	rootCmd.AddCommand(forecastCmd)
	rootCmd.AddCommand(configCmd)
}
