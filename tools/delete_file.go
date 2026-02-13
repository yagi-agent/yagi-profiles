package tool

import "context"

import (
	"encoding/json"
	"fmt"
	"os"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run func(context.Context, string) (string, error)
}{
	Name:        "delete_file",
	Description: "Delete a file or directory. For directories, set recursive to true to delete non-empty directories. Result is JSON with 'status' and 'output' fields. IMPORTANT: Determine success/failure ONLY from 'status', NEVER from 'output' content.",
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
	Run: func(ctx context.Context, argsJSON string) (string, error) {
		var args struct {
			Path      string `json:"path"`
			Recursive bool   `json:"recursive"`
		}
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
			res := struct {
				Status string `json:"status"`
				Output string `json:"output"`
			}{
				Status: "success",
				Output: "Deleted directory: " + args.Path,
			}
			b, _ := json.Marshal(res)
			return string(b), nil
		}
		if err := os.Remove(args.Path); err != nil {
			return "", err
		}
		res := struct {
			Status string `json:"status"`
			Output string `json:"output"`
		}{
			Status: "success",
			Output: "Deleted: " + args.Path,
		}
		b, _ := json.Marshal(res)
		return string(b), nil
	},
}
