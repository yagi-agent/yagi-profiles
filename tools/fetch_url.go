package tool

import (
	"encoding/json"
	"hostapi"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) (string, error)
}{
	Name:        "fetch_url",
	Description: "Fetch the content of a URL and return it as text. HTML pages are converted to plain text with links preserved.",
	Parameters: `{
		"type": "object",
		"properties": {
			"url": {
				"type": "string",
				"description": "The URL to fetch"
			}
		},
		"required": ["url"]
	}`,
	Run: func(args string) (string, error) {
		var params struct {
			URL string `json:"url"`
		}
		if err := json.Unmarshal([]byte(args), &params); err != nil {
			return "", err
		}
		return hostapi.FetchURL(params.URL), nil
	},
}
