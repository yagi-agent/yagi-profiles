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
	Run         func(string) (string, error)
}{
	Name:        "make_directory",
	Description: "Create a directory (including parent directories if needed)",
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
	Run: func(argsJSON string) (string, error) {
		var args MakeDirectoryArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}
		if err := os.MkdirAll(args.Path, 0755); err != nil {
			return "", err
		}
		return "Created directory: " + args.Path, nil
	},
}
