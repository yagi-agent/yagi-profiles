package main

import (
	"encoding/json"
	"fmt"
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
	Run         func(string) (string, error)
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
	Run: func(argsJSON string) (string, error) {
		var args DeleteFileArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}

		info, err := os.Stat(args.Path)
		if err != nil {
			return "", err
		}

		if info.IsDir() {
			entries, err := os.ReadDir(args.Path)
			if err != nil {
				return "", err
			}
			if len(entries) > 0 && !args.Recursive {
				return "", fmt.Errorf("directory not empty, use recursive to delete")
			}
			if err := os.RemoveAll(args.Path); err != nil {
				return "", err
			}
			return "Deleted directory: " + args.Path, nil
		}
		if err := os.Remove(args.Path); err != nil {
			return "", err
		}
		return "Deleted: " + args.Path, nil
	},
}
