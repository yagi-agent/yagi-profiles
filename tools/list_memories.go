package tool

import (
	"hostapi"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) (string, error)
}{
	Name:        "list_memories",
	Description: "List all saved information. Returns a JSON object with all key-value pairs that have been remembered.",
	Parameters: `{
		"type": "object",
		"properties": {}
	}`,
	Run: func(args string) (string, error) {
		return hostapi.ListMemory(), nil
	},
}
