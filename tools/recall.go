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
	Name:        "recall",
	Description: "Recall previously saved information. Use this to retrieve facts that were saved using the 'remember' tool.",
	Parameters: `{
		"type": "object",
		"properties": {
			"key": {
				"type": "string",
				"description": "The identifier of the information to recall (e.g., 'user_name', 'favorite_language')"
			}
		},
		"required": ["key"]
	}`,
	Run: func(args string) (string, error) {
		var params struct {
			Key string `json:"key"`
		}
		if err := json.Unmarshal([]byte(args), &params); err != nil {
			return "", err
		}
		value, err := hostapi.GetMemory(params.Key)
		if err != nil {
			return "", err
		}
		if value == "" {
			return "No information found for key: " + params.Key, nil
		}
		return value, nil
	},
}
