package main

import (
	"encoding/json"
	"os"
)

type ReadFileArgs struct {
	Path string `json:"path"`
}

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
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
	Run: func(argsJSON string) string {
		var args ReadFileArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "Error: " + err.Error()
		}

		content, err := os.ReadFile(args.Path)
		if err != nil {
			return "Error reading file: " + err.Error()
		}

		return string(content)
	},
}
