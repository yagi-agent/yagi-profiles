package tool

import (
	"encoding/json"
	"os"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) (string, error)
}{
	Name:        "write_file",
	Description: "Write content to a file. WARNING: This can overwrite existing files.",
	Parameters: `{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "Path to the file to write"
			},
			"content": {
				"type": "string",
				"description": "Content to write to the file"
			}
		},
		"required": ["path", "content"]
	}`,
	Run: func(argsJSON string) (string, error) {
		var args struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}

		if err := os.WriteFile(args.Path, []byte(args.Content), 0644); err != nil {
			return "", err
		}

		return "Successfully wrote to " + args.Path, nil
	},
}
