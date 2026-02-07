package tool

import (
	"encoding/json"
	"hostapi"
)

var Name = "recall"
var Description = "Recall previously saved information. Use this to retrieve facts that were saved using the 'remember' tool."
var Parameters = `{
	"type": "object",
	"properties": {
		"key": {
			"type": "string",
			"description": "The identifier of the information to recall (e.g., 'user_name', 'favorite_language')"
		}
	},
	"required": ["key"]
}`

func Run(args string) (string, error) {
	var params struct {
		Key string `json:"key"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "", err
	}
	value := hostapi.GetMemory(params.Key)
	if value == "" {
		return "No information found for key: " + params.Key, nil
	}
	return value, nil
}
