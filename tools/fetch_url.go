package tool

import (
	"encoding/json"
	"hostapi"
)

var Name = "fetch_url"
var Description = "Fetch the content of a URL and return it as text. HTML pages are converted to plain text with links preserved."
var Parameters = `{
	"type": "object",
	"properties": {
		"url": {
			"type": "string",
			"description": "The URL to fetch"
		}
	},
	"required": ["url"]
}`

func Run(args string) (string, error) {
	var params struct {
		URL string `json:"url"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", err
	}
	return hostapi.FetchURL(params.URL), nil
}
