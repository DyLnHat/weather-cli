package display

import (
	"fmt"
	"strings"
	"time"
	"unicode"

	"github.com/DyLnHat/weather-cli/internal/api"
	"github.com/fatih/color"
)

//	--|-- Color Palette --|--

var (
	colorTitle   = color.New(color.FgCyan, color.Bold)
	colorLabel   = color.New(color.FgHiBlack)
	colorValue   = color.New(color.FgWhite)
	colorHot     = color.New(color.FgRed, color.Bold)
	colorCold    = color.New(color.FgBlue, color.Bold)
	colorMild    = color.New(color.FgYellow, color.Bold)
	colorSuccess = color.New(color.FgGreen, color.Bold)
	colorError   = color.New(color.FgRed, color.Bold)
)

//	--|-- Helpers --|--

// reutrns an emoji for a given condition
func weatherIcon(condition string) string {
	icons := map[string]string{
		"Clear":        "☀️ ",
		"Clouds":       "☁️ ",
		"Rain":         "🌧️ ",
		"Drizzle":      "🌦️ ",
		"Thunderstorm": "⛈️ ",
		"Snow":         "❄️ ",
		"Mist":         "🌫️ ",
		"Fog":          "🌫️ ",
		"Haze":         "🌫️ ",
		"Smoke":        "🌫️ ",
		"Dust":         "🌪️ ",
		"Sand":         "🌪️ ",
		"Tornado":      "🌪️ ",
	}
	if icon, ok := icons[condition]; ok {
		return icon
	}
	return "🌡️ "
}

// returns the temperature colored by heat level
func tempColored(temp float64, units string) string {
	hotThreshold := 30.0
	coldThreshold := 10.0
	if units == "imperial" {
		hotThreshold = 86.0
		coldThreshold = 50.0
	}

	unit := "ºC"
	switch units {
	case "imperial":
		unit = "ºF"
	case "standard":
		unit = "K"
	}

	formatted := fmt.Sprintf("%.1f%s", temp, unit)
	switch {
	case temp >= hotThreshold:
		return colorHot.Sprint(formatted)
	case temp <= coldThreshold:
		return colorCold.Sprint(formatted)
	default:
		return colorMild.Sprint(formatted)
	}
}

// Converts degrees to a compass direction
func windDirection(deg int) string {
	dirs := []string{"N", "NE", "E", "SE", "S", "SW", "W", "NW"}
	return dirs[(deg+22)/45%8]
}

// Capitalises the first character of each word
func titleCase(s string) string {
	words := strings.Fields(s)
	for i, w := range words {
		runes := []rune(w)
		if len(runes) > 0 {
			runes[0] = unicode.ToUpper(runes[0])
		}
		words[i] = string(runes)
	}
	return strings.Join(words, " ")
}

// Prints a Two-column line with padding
func row(icon, labelStr, val string) {
	fmt.Printf(" %s %-22s %s\n", icon, colorLabel.Sprint(labelStr), val)
}

//	--|-- Public print functions --|--

// Renders the full current weather block
func PrintCurrent(w *api.CurrentWeather, units string) {
	icon := weatherIcon(w.Weather[0].Main)
	desc := titleCase(w.Weather[0].Description)
	sunrise := time.Unix(w.Sys.Sunrise, 0).Format("15:04")
	sunset := time.Unix(w.Sys.Sunset, 0).Format("15:04")
	updated := time.Unix(w.Dt, 0).Format("02 Jan 2006 ~ 15:04")

	fmt.Println()
	colorTitle.Printf("  %s %s, %s\n", icon, w.Name, w.Sys.Country)
	colorLabel.Printf("  %s\n", desc)
	fmt.Println()

	row("🌡️ ", "Temperature:", tempColored(w.Main.Temp, units))
	row("🤔", "Feels like:", tempColored(w.Main.FeelsLike, units))
	row("📉", "Min / Max:", fmt.Sprintf("%s  /  %s",
		tempColored(w.Main.TempMin, units),
		tempColored(w.Main.TempMax, units)))

	fmt.Println()

	row("💧", "Humidity:", colorValue.Sprintf("%d%%", w.Main.Humidity))
	row("💨", "Wind:", colorValue.Sprintf("%.1f m/s %s",
		w.Wind.Speed, windDirection(w.Wind.Deg)))
	row("📊", "Pressure:", colorValue.Sprintf("%d hPa", w.Main.Pressure))

	if w.Visibility > 0 {
		row("👁️ ", "Visibility:", colorValue.Sprintf("%.1f km", float64(w.Visibility)/1000))
	}

	fmt.Println()
	row("🌅", "Sunrise:", colorValue.Sprint(sunrise))
	row("🌇", "Sunset:", colorValue.Sprint(sunset))
	fmt.Println()

	colorLabel.Printf("  Last updated: %s\n\n", updated)
}

// Renders the N-day forecast table
func PrintForecast(f *api.Forecast, units string) {
	fmt.Println()
	colorTitle.Printf("  📅  Forecast — %s, %s\n\n", f.City.Name, f.City.Country)

	seen := map[string]bool{}
	printed := 0

	for _, item := range f.List {
		t := time.Unix(item.Dt, 0)
		day := t.Format("19-01-2026")
		hour := t.Hour()

		// One slot per day, prefer midday (11:00–15:00)
		if seen[day] {
			continue
		}
		if hour < 11 || hour > 15 {
			continue
		}
		seen[day] = true
		printed++

		icon := weatherIcon(item.Weather[0].Main)
		pop := int(item.Pop * 100)
		desc := titleCase(item.Weather[0].Description)
		dateStr := t.Format("Mon 02 Jan")

		fmt.Printf("  %s %-12s  %-10s  %-22s  💧 %3d%%  💨 %.1f m/s\n",
			icon,
			colorValue.Sprint(dateStr),
			tempColored(item.Main.Temp, units),
			colorLabel.Sprint(desc),
			pop,
			item.Wind.Speed,
		)
	}

	// if no midday slots found
	if printed == 0 {
		for _, item := range f.List {
			t := time.Unix(item.Dt, 0)
			day := t.Format("2006-01-02")
			if seen[day] {
				continue
			}
			seen[day] = true

			icon := weatherIcon(item.Weather[0].Main)
			pop := int(item.Pop * 100)
			desc := titleCase(item.Weather[0].Description)
			dateStr := t.Format("Mon 02 Jan")

			fmt.Printf("  %s %-12s  %-10s  %-22s  💧 %3d%%  💨 %.1f m/s\n",
				icon,
				colorValue.Sprint(dateStr),
				tempColored(item.Main.Temp, units),
				colorLabel.Sprint(desc),
				pop,
				item.Wind.Speed,
			)
		}
	}

	fmt.Println()
}

// Renders the current configuration
func PrintConfig(apiKey, units string) {
	fmt.Println()
	colorTitle.Printf("  ⚙️   Current Configuration\n\n")

	maskedKey := "not set"
	if len(apiKey) >= 8 {
		maskedKey = apiKey[:8] + strings.Repeat("*", len(apiKey)-8)
	} else if apiKey != " " {
		maskedKey = strings.Repeat("*", len(apiKey))
	}

	row("🔑", "API Key:", colorValue.Sprint(maskedKey))
	row("📏", "Units:", colorValue.Sprint(units))
	fmt.Println()
}

// Renders a green success message
func PrintSuccess(msg string) {
	fmt.Println()
	colorSuccess.Printf("  ✅  %s\n\n", msg)
}

// Renders a red error message
func PrintError(msg string) {
	fmt.Println()
	colorError.Printf("  ❌  %s\n\n", msg)
}
