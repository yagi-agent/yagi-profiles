package tool

import "context"

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run func(context.Context, string) (string, error)
}{
	Name:        "list_files",
	Description: "List files and directories in a given path. Result is JSON with 'status' and 'output' fields. Determine success/failure ONLY from 'status', NEVER from 'output' content.",
	Parameters: `{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "Directory path to list (default: current directory)"
			}
		}
	}`,
	Run: func(ctx context.Context, argsJSON string) (string, error) {
		var args struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}

		if args.Path == "" {
			args.Path = "."
		}

		entries, err := os.ReadDir(args.Path)
		if err != nil {
			return "", err
		}

		var result strings.Builder
		for _, entry := range entries {
			info, err := entry.Info()
			if err != nil {
				continue
			}

			prefix := "F"
			if entry.IsDir() {
				prefix = "D"
			}

			result.WriteString(prefix)
			result.WriteString(" ")
			result.WriteString(entry.Name())
			result.WriteString(" (")
			result.WriteString(formatSize(info.Size()))
			result.WriteString(")\n")
		}

		res := struct {
			Status string `json:"status"`
			Output string `json:"output"`
		}{
			Status: "success",
			Output: result.String(),
		}
		b, _ := json.Marshal(res)
		return string(b), nil
	},
}

func formatSize(size int64) string {
	if size < 1024 {
		return fmt.Sprintf("%dB", size)
	} else if size < 1024*1024 {
		return fmt.Sprintf("%dKB", size/1024)
	} else {
		return fmt.Sprintf("%dMB", size/(1024*1024))
	}
}
