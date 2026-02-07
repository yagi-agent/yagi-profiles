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
	Name:        "remember",
	Description: "IMPORTANT: Always use this tool when the user shares personal information like their name, preferences, location, or any facts they want you to remember. Save information that should persist across conversations.",
	Parameters: `{
		"type": "object",
		"properties": {
			"key": {
				"type": "string",
				"description": "A short identifier for what to remember (e.g., 'user_name', 'favorite_language', 'project_name')"
			},
			"value": {
				"type": "string",
				"description": "The information to remember"
			}
		},
		"required": ["key", "value"]
	}`,
	Run: func(args string) (string, error) {
		var params struct {
			Key   string `json:"key"`
			Value string `json:"value"`
		}
		if err := json.Unmarshal([]byte(args), &params); err != nil {
			return "", err
		}
		return hostapi.SaveMemory(params.Key, params.Value)
	},
}
