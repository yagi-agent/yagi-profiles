package tool

import "context"

import (
	"encoding/json"
	"os"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run func(context.Context, string) (string, error)
}{
	Name:        "read_file",
	Description: "Read the contents of a file",
	Parameters: `{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "Path to the file to read"
			}
		},
		"required": ["path"]
	}`,
	Run: func(ctx context.Context, argsJSON string) (string, error) {
		var args struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}

		content, err := os.ReadFile(args.Path)
		if err != nil {
			return "", err
		}

		return string(content), nil
	},
}
