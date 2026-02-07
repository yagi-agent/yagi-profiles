package main

import (
	"encoding/json"
	"os"
)

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
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
	Run: func(argsJSON string) string {
		var params struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &params); err != nil {
			return "Error: " + err.Error()
		}
		if err := os.MkdirAll(params.Path, 0755); err != nil {
			return "Error: " + err.Error()
		}
		return "Created directory: " + params.Path
	},
}
