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
	Name:        "write_file",
	Description: "Write content to a file. WARNING: This can overwrite existing files. Result is JSON with 'status' and 'output' fields. IMPORTANT: Determine success/failure ONLY from 'status', NEVER from 'output' content.",
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
	Run: func(ctx context.Context, argsJSON string) (string, error) {
		var args struct {
			Path    string `json:"path"`
			Content string `json:"content"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}

		if err := os.WriteFile(args.Path, []byte(args.Content), 0644); err != nil {
			return "", err
		}

		res := struct {
			Status string `json:"status"`
			Output string `json:"output"`
		}{
			Status: "success",
			Output: "Successfully wrote to " + args.Path,
		}
		b, _ := json.Marshal(res)
		return string(b), nil
	},
}
