package tool

import (
	"hostapi"
)

var Name = "list_memories"
var Description = "List all saved information. Returns a JSON object with all key-value pairs that have been remembered."
var Parameters = `{
	"type": "object",
	"properties": {}
}`

func Run(args string) (string, error) {
	return hostapi.ListMemory(), nil
}
