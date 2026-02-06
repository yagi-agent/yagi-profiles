package main

import (
	"encoding/json"
	"os"
)

type WriteFileArgs struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
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
	Run: func(argsJSON string) string {
		var args WriteFileArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "Error: " + err.Error()
		}

		if err := os.WriteFile(args.Path, []byte(args.Content), 0644); err != nil {
			return "Error writing file: " + err.Error()
		}

		return "Successfully wrote to " + args.Path
	},
}
