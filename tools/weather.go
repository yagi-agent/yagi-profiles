package tool

import "context"

import (
	"encoding/json"
	"fmt"
	"hostapi"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(context.Context, string) (string, error)
}{
	Name:        "weather",
	Description: "Get the current weather information for a specified city. Returns weather conditions, temperature, humidity, wind speed, etc.",
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

		url := fmt.Sprintf("https://wttr.in/%s?format=j1", args.City)
		return hostapi.FetchURL(ctx, url, nil)
	},
}
