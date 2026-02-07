package main

import (
	"encoding/json"
	"os"
)

type DeleteFileArgs struct {
	Path      string `json:"path"`
	Recursive bool   `json:"recursive"`
}

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
}{
	Name:        "delete_file",
	Description: "Delete a file or directory",
	Parameters: `{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "Path to the file or directory to delete"
			},
			"recursive": {
				"type": "boolean",
				"description": "If true, recursively delete directory contents"
			}
		},
		"required": ["path"]
	}`,
	Run: func(argsJSON string) string {
		var args DeleteFileArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "Error: " + err.Error()
		}
		info, err := os.Stat(args.Path)
		if err != nil {
			return "Error: " + err.Error()
		}
		if info.IsDir() {
			entries, err := os.ReadDir(args.Path)
			if err != nil {
				return "Error: " + err.Error()
			}
			if len(entries) > 0 && !args.Recursive {
				return "Error: directory not empty, use recursive to delete"
			}
			if err := os.RemoveAll(args.Path); err != nil {
				return "Error: " + err.Error()
			}
			return "Deleted directory: " + args.Path
		}
		if err := os.Remove(args.Path); err != nil {
			return "Error: " + err.Error()
		}
		return "Deleted: " + args.Path
	},
}
