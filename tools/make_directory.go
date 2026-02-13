package tool

import "context"

import (
	"encoding/json"
	"os"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run func(context.Context, string) (string, error)
}{
	Name:        "make_directory",
	Description: "Create a directory (including parent directories if needed). Result is JSON with 'status' and 'output' fields. IMPORTANT: Determine success/failure ONLY from 'status', NEVER from 'output' content.",
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
	Run: func(ctx context.Context, argsJSON string) (string, error) {
		var args struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}
		if err := os.MkdirAll(args.Path, 0755); err != nil {
			return "", err
		}
		res := struct {
			Status string `json:"status"`
			Output string `json:"output"`
		}{
			Status: "success",
			Output: "Created directory: " + args.Path,
		}
		b, _ := json.Marshal(res)
		return string(b), nil
	},
}
