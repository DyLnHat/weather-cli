# 🌤️ weather-cli

A terminal weather tool built with Go. Get current conditions and 5-day forecasts for any city in the world, right from your command line.

![Go](https://img.shields.io/badge/Go-1.21-00ADD8?style=flat&logo=go)
![License](https://img.shields.io/badge/license-MIT-green?style=flat)
![Docker](https://img.shields.io/badge/Docker-ready-2496ED?style=flat&logo=docker)

---

## Demo

```
☀️  Madrid, ES
    Clear sky

🌡️  Temperature:          12.5°C
🤔  Feels like:           11.6°C
📉  Min / Max:            12.0°C  /  13.7°C

💧  Humidity:             57%
💨  Wind:                 3.2 m/s SW
📊  Pressure:             1027 hPa
👁️  Visibility:           10.0 km

🌅  Sunrise:              08:02
🌇  Sunset:               20:36

Last updated: 30 Mar 2026 · 22:00
```

---

## Features

- 🌡️ Current weather — temperature, feels like, min/max, humidity, wind, pressure, visibility, sunrise & sunset
- 📅 5-day forecast — daily summary with precipitation probability and wind speed
- 🎨 Color-coded temperatures — 🔴 hot / 🔵 cold / 🟡 mild
- ⚙️ Persistent config — API key and units saved to `~/.weather-cli/config.json`
- 🐳 Docker ready — multi-stage build, image under 20MB

---

## Installation

### Option A — Build from source

**Requirements:** Go 1.21+

```bash
git clone https://github.com/DyLnHat/weather-cli.git
cd weather-cli
go build -o weather .        # Linux / macOS
go build -o weather.exe .    # Windows
```

### Option B — Docker

```bash
docker build -t weather-cli .
```

---

## Setup

Get a free API key at [openweathermap.org/api](https://openweathermap.org/api) and run:

```bash
# Linux / macOS
./weather config --set-key YOUR_API_KEY

# Windows
.\weather.exe config --set-key YOUR_API_KEY
```

> ⚠️ New API keys can take up to 10 minutes to activate after registration.

---

## Usage

### Current weather

```bash
./weather current Madrid
./weather current "New York"
./weather current Tokyo
```

### Forecast

```bash
./weather forecast Madrid              # 5 days (default)
./weather forecast London --days 3     # 1 to 5 days
```

### Configuration

```bash
./weather config --set-key YOUR_API_KEY        # Save API key
./weather config --set-units imperial          # metric | imperial | standard
./weather config --show                        # Show current config
```

### Docker

```bash
# Mount your config so the API key persists between runs
docker run --rm \
  -v ~/.weather-cli:/root/.weather-cli \
  weather-cli current Madrid

docker run --rm \
  -v ~/.weather-cli:/root/.weather-cli \
  weather-cli forecast Tokyo --days 3
```

---

## Project Structure

```
weather-cli/
├── cmd/
│   ├── root.go          # Root Cobra command
│   ├── current.go       # `weather current` command
│   ├── forecast.go      # `weather forecast` command
│   └── config.go        # `weather config` command
├── internal/
│   ├── api/
│   │   └── openweather.go   # OpenWeatherMap HTTP client
│   ├── config/
│   │   └── config.go        # JSON config read/write
│   └── display/
│       └── display.go       # Colored terminal output
├── main.go
├── Dockerfile
└── go.mod
```

---

## Tech Stack

| Tool | Purpose |
|---|---|
| [Go 1.21](https://go.dev) | Language |
| [Cobra](https://github.com/spf13/cobra) | CLI framework |
| [fatih/color](https://github.com/fatih/color) | Terminal colors |
| [OpenWeatherMap API](https://openweathermap.org/api) | Weather data |
| [Docker](https://www.docker.com) | Containerisation |

---

## License

MIT — see [LICENSE](LICENSE) for details.
