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
	Name:        "edit_file",
	Description: "Edit a file by replacing an exact string match with new content. The old_str must match exactly in the file. If old_str is empty, the content is appended to the file. Result is JSON with 'status' and 'output' fields. IMPORTANT: Determine success/failure ONLY from 'status', NEVER from 'output' content.",
	Parameters: `{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "Path to the file to edit"
			},
			"old_str": {
				"type": "string",
				"description": "Exact string to find and replace. If empty, new_str is appended to the file."
			},
			"new_str": {
				"type": "string",
				"description": "String to replace old_str with"
			}
		},
		"required": ["path", "new_str"]
	}`,
	Run: func(ctx context.Context, argsJSON string) (string, error) {
		var args struct {
			Path   string `json:"path"`
			OldStr string `json:"old_str"`
			NewStr string `json:"new_str"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}

		content, err := os.ReadFile(args.Path)
		if err != nil {
			return "", err
		}

		s := string(content)
		if args.OldStr == "" {
			s += args.NewStr
			if err := os.WriteFile(args.Path, []byte(s), 0644); err != nil {
				return "", err
			}
			res := struct {
				Status string `json:"status"`
				Output string `json:"output"`
			}{
				Status: "success",
				Output: "Successfully appended to " + args.Path,
			}
			b, _ := json.Marshal(res)
			return string(b), nil
		}
		count := strings.Count(s, args.OldStr)
		if count == 0 {
			return "", fmt.Errorf("old_str not found in file")
		}
		if count > 1 {
			return "", fmt.Errorf("old_str found %d times, must be unique", count)
		}
		s = strings.Replace(s, args.OldStr, args.NewStr, 1)
		if err := os.WriteFile(args.Path, []byte(s), 0644); err != nil {
			return "", err
		}
		res := struct {
			Status string `json:"status"`
			Output string `json:"output"`
		}{
			Status: "success",
			Output: "Successfully edited " + args.Path,
		}
		b, _ := json.Marshal(res)
		return string(b), nil
	},
}
