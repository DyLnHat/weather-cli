package api

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"time"
)

const baseURL = "https://api.openweathermap.org/data/2.5"

// Sentinel errors for clean error handling in cmd layer
var (
	ErrInvalidAPIKey = errors.New("invalid API key — check your key at openweathermap.org")
	ErrCityNotFound  = errors.New("city not found — check the spelling")
	ErrTimeout       = errors.New("request timed out after 10s — try again") //change to 30
	ErrNoConnection  = errors.New("connection failed — check your internet connection")
)

// OpenWeather API Client
type Client struct {
	APIKey     string
	Units      string
	httpClient *http.Client
}

// NewClient with a 30s timeout
func NewClient(apiKey, units string) *Client {
	return &Client{
		APIKey: apiKey,
		Units:  units,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// GET request and HTTP error handling
func (c *Client) doGet(endpoint string) ([]byte, error) {
	resp, err := c.httpClient.Get(endpoint)
	if err != nil {
		if errors.Is(err, context_DeadlineExceeded(err)) {
			return nil, ErrTimeout
		}
		return nil, ErrNoConnection
	}
	defer resp.Body.Close()

	switch resp.StatusCode {
	case http.StatusOK:
		//	Continue
	case http.StatusUnauthorized:
		return nil, ErrInvalidAPIKey
	case http.StatusNotFound:
		return nil, ErrCityNotFound
	default:
		return nil, fmt.Errorf("unexpected response from API: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}
	return body, nil
}

// (context_DeadlineExceeded) checks if an error is a timeout
func context_DeadlineExceeded(err error) error {
	if err, ok := err.(interface{ Timeout() bool }); ok && err.Timeout() {
		return ErrTimeout
	}
	return nil
}

// --|-- Current Weather --|--
// Maps /data/2.5/weather JSON response
type CurrentWeather struct {
	Name string `json:"name"`
	Sys  struct {
		Country string `json:"country"`
		Sunrise int64  `json:"sunrise"`
		Sunset  int64  `json:"sunset"`
	} `json:"sys"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Humidity  int     `json:"humidity"`
		Pressure  int     `json:"pressure"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
	} `json:"wind"`
	Visibility int   `json:"visibility"`
	Dt         int64 `json:"dt"`
}

// Fetches the current weather for a city
func (c *Client) GetCurrent(city string) (*CurrentWeather, error) {
	endpoint := fmt.Sprintf("%s/weather?q=%s&appid=%s&units=%s",
		baseURL, url.QueryEscape(city), c.APIKey, c.Units)

	body, err := c.doGet(endpoint)
	if err != nil {
		if errors.Is(err, ErrCityNotFound) {
			return nil, fmt.Errorf("city %q not found - type the correct name", city)
		}
		return nil, err
	}

	var result CurrentWeather
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}

// --|-- Forecast --|--
// Maps a single 4h slot in the forecast response
type ForecastItem struct {
	Dt   int64 `json:"dt"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
	Pop float64 `json:"pop"` // Probability of precipitation
}

// Maps /data/2.5/forecast JSON response
type Forecast struct {
	City struct {
		Name    string `json:"name"`
		Country string `json:"country"`
	} `json:"city"`
	List []ForecastItem `json:"list"`
}

// Fetches the currect Forecast for the given number of days ( 1-5 )
func (c *Client) GetForecast(city string, days int) (*Forecast, error) {
	//	API returns 3h slots, days*8 = N days worth of data.
	cnt := days * 8
	endpoint := fmt.Sprintf("%s/forecast?q=%s&appid=%s&units=%s&cnt=%d",
		baseURL, url.QueryEscape(city), c.APIKey, c.Units, cnt)

	body, err := c.doGet(endpoint)
	if err != nil {
		if errors.Is(err, ErrCityNotFound) {
			return nil, fmt.Errorf("city %q not found — check the spelling", city)
		}
		return nil, err
	}

	var result Forecast
	if err := json.Unmarshal(body, &result); err != nil {
		return nil, fmt.Errorf("failed to parse response: %w", err)
	}
	return &result, nil
}
