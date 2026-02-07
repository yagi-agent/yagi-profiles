package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

type EditFileArgs struct {
	Path   string `json:"path"`
	OldStr string `json:"old_str"`
	NewStr string `json:"new_str"`
}

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
}{
	Name:        "edit_file",
	Description: "Edit a file by replacing old_str with new_str. If old_str is empty, new_str is appended.",
	Parameters: `{
		"type": "object",
		"properties": {
			"path": {
				"type": "string",
				"description": "Path to the file to edit"
			},
			"old_str": {
				"type": "string",
				"description": "Text to find and replace"
			},
			"new_str": {
				"type": "string",
				"description": "Replacement text"
			}
		},
		"required": ["path", "old_str", "new_str"]
	}`,
	Run: func(argsJSON string) string {
		var args EditFileArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "Error: " + err.Error()
		}
		content, err := os.ReadFile(args.Path)
		if err != nil {
			return "Error reading file: " + err.Error()
		}
		s := string(content)
		if args.OldStr == "" {
			s += args.NewStr
			if err := os.WriteFile(args.Path, []byte(s), 0644); err != nil {
				return "Error writing file: " + err.Error()
			}
			return "Successfully appended to " + args.Path
		}
		count := strings.Count(s, args.OldStr)
		if count == 0 {
			return "Error: old_str not found in file"
		}
		if count > 1 {
			return fmt.Sprintf("Error: old_str found %d times, must be unique", count)
		}
		s = strings.Replace(s, args.OldStr, args.NewStr, 1)
		if err := os.WriteFile(args.Path, []byte(s), 0644); err != nil {
			return "Error writing file: " + err.Error()
		}
		return "Successfully edited " + args.Path
	},
}
