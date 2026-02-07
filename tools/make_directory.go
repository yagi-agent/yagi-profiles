package main

import (
	"encoding/json"
	"os"
)

type MakeDirectoryArgs struct {
	Path string `json:"path"`
}

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
}{
	Name:        "make_directory",
	Description: "Create a directory (including parent directories)",
	Parameters: `{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "Path of the directory to create"
			}
		},
		"required": ["path"]
	}`,
	Run: func(argsJSON string) string {
		var args MakeDirectoryArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "Error: " + err.Error()
		}
		if err := os.MkdirAll(args.Path, 0755); err != nil {
			return "Error creating directory: " + err.Error()
		}
		return "Created directory: " + args.Path
	},
}
