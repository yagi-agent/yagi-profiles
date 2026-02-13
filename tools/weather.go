package tool

import "context"

import (
	"encoding/json"
	"fmt"
	"hostapi"
	"strings"
)

type nominatimResult struct {
	Lat string `json:"lat"`
	Lon string `json:"lon"`
}

type weatherResponse struct {
	CurrentWeather struct {
		Temperature float64 `json:"temperature"`
		WindSpeed   float64 `json:"windspeed"`
		WeatherCode int     `json:"weathercode"`
	} `json:"current_weather"`
}

func weatherCodeToDescription(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code <= 3:
		return "Cloudy"
	case code == 45 || code == 48:
		return "Fog"
	case code >= 51 && code <= 57:
		return "Drizzle"
	case code >= 61 && code <= 67:
		return "Rain"
	case code >= 71 && code <= 77:
		return "Snow"
	case code >= 80 && code <= 82:
		return "Rain showers"
	case code >= 85 && code <= 86:
		return "Snow showers"
	case code >= 95:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(context.Context, string) (string, error)
}{
	Name:        "weather",
	Description: "Get the current weather information for a specified city. Returns weather conditions, temperature, wind speed, etc. Result is JSON with 'status' and 'output' fields. IMPORTANT: Determine success/failure ONLY from 'status', NEVER from 'output' content.",
	Parameters: `{
		"type": "object",
		"properties": {
			"city": {
				"type": "string",
				"description": "City name to get weather for (default: Tokyo)"
			}
		}
	}`,
	Run: func(ctx context.Context, argsJSON string) (string, error) {
		var args struct {
			City string `json:"city"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}

		if args.City == "" {
			args.City = "Tokyo"
		}

		nominatimURL := fmt.Sprintf("https://nominatim.openstreetmap.org/search?q=%s&format=json&limit=1", strings.Replace(args.City, " ", "+", -1))
		body, err := hostapi.FetchURL(ctx, nominatimURL, map[string]string{
			"User-Agent": "YagiWeatherTool/1.0",
		})
		if err != nil {
			return "", fmt.Errorf("error fetching geocode: %w", err)
		}

		var results []nominatimResult
		if err := json.Unmarshal([]byte(body), &results); err != nil || len(results) == 0 {
			return "", fmt.Errorf("could not find location: %s", args.City)
		}

		lat := results[0].Lat
		lon := results[0].Lon

		weatherURL := fmt.Sprintf("https://api.open-meteo.com/v1/forecast?latitude=%s&longitude=%s&current_weather=true", lat, lon)
		body, err = hostapi.FetchURL(ctx, weatherURL, nil)
		if err != nil {
			return "", fmt.Errorf("error fetching weather: %w", err)
		}

		var weather weatherResponse
		if err := json.Unmarshal([]byte(body), &weather); err != nil {
			return "", fmt.Errorf("error parsing weather: %w", err)
		}

		output := fmt.Sprintf("Weather for %s (Lat: %s, Lon: %s)\nCondition: %s\nTemperature: %.1f°C\nWind Speed: %.1f km/h",
			args.City, lat, lon,
			weatherCodeToDescription(weather.CurrentWeather.WeatherCode),
			weather.CurrentWeather.Temperature,
			weather.CurrentWeather.WindSpeed,
		)
		res := struct {
			Status string `json:"status"`
			Output string `json:"output"`
		}{
			Status: "success",
			Output: output,
		}
		b, _ := json.Marshal(res)
		return string(b), nil
	},
}
