package cmd

import (
	"fmt"

	"github.com/DyLnHat/weather-cli/internal/config"
	"github.com/DyLnHat/weather-cli/internal/display"
	"github.com/spf13/cobra"
)

var (
	setKey   string
	setUnits string
	showCfg  bool
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage weather-cli configuration",
	Example: `weather config --set-key abc123xyz
  		weather config --set-units imperial
  		weather config --set-units metric
  		weather config --show`,

	RunE: func(cmd *cobra.Command, args []string) error {
		if !showCfg && setKey == "" && setUnits == "" {
			return fmt.Errorf("provide at least one flag: --set-key, --set-units or --show")
		}

		cfg, err := config.Load()
		if err != nil {
			return fmt.Errorf("could not load config: %w", err)
		}

		// Show current config
		if showCfg {
			display.PrintConfig(cfg.APIKey, cfg.Units)
			return nil
		}

		//	Save API key
		if setKey != "" {
			cfg.APIKey = setKey
			if err := config.Save(cfg); err != nil {
				return fmt.Errorf("could not save config: %w", err)
			}
			preview := setKey
			if len(preview) > 8 {
				preview = preview[:8] + "..."
			}
			display.PrintSuccess(fmt.Sprintf("API key saved (%s)", preview))
		}

		// Save units
		if setUnits != "" {
			switch setUnits {
			case "metric", "imperial", "standard":
				cfg.Units = setUnits
				if err := config.Save(cfg); err != nil {
					return fmt.Errorf("could not save config: %w", err)
				}
				display.PrintSuccess(fmt.Sprintf("Units set to: %s", setUnits))
			default:
				return fmt.Errorf("invalid units %q — use: metric, imperial or standard", setUnits)
			}
		}

		return nil
	},
}

func init() {
	configCmd.Flags().StringVar(&setKey, "set-key", "", "Set your OpenWeatherMap API key")
	configCmd.Flags().StringVar(&setUnits, "set-units", "", "Set temperature units (metric|imperial|standard)")
	configCmd.Flags().BoolVar(&showCfg, "show", false, "Show current configuration")
}
