package main

import (
	"encoding/json"
	"fmt"
	"os"
)

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
}{
	Name:        "delete_file",
	Description: "Delete a file or directory. For directories, set recursive to true to delete non-empty directories.",
	Parameters: `{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "Path to the file or directory to delete"
			},
			"recursive": {
				"type": "boolean",
				"description": "If true, recursively delete directory and its contents (default: false)"
			}
		},
		"required": ["path"]
	}`,
	Run: func(argsJSON string) string {
		var params struct {
			Path      string `json:"path"`
			Recursive bool   `json:"recursive"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &params); err != nil {
			return "Error: " + err.Error()
		}

		info, err := os.Stat(params.Path)
		if err != nil {
			return "Error: " + err.Error()
		}

		if info.IsDir() && !params.Recursive {
			entries, _ := os.ReadDir(params.Path)
			if len(entries) > 0 {
				return fmt.Sprintf("Error: directory %s is not empty, set recursive to true", params.Path)
			}
		}

		if params.Recursive {
			err = os.RemoveAll(params.Path)
		} else {
			err = os.Remove(params.Path)
		}
		if err != nil {
			return "Error: " + err.Error()
		}
		return fmt.Sprintf("Deleted %s", params.Path)
	},
}
