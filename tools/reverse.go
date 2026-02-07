package tool

import "encoding/json"

var Name = "reverse"
var Description = "Reverse the input string"
var Parameters = `{
	"type": "object",
	"properties": {
		"text": {
			"type": "string",
			"description": "The text to reverse"
		}
	},
	"required": ["text"]
}`

func Run(args string) string {
	var params struct {
		Text string `json:"text"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil {
		return "Error: " + err.Error()
	}
	runes := []rune(params.Text)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}
