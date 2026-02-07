package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
)

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
}{
	Name:        "edit_file",
	Description: "Edit a file by replacing an exact string match with new content. The old_str must match exactly in the file. If old_str is empty, the content is appended to the file.",
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
	Run: func(argsJSON string) string {
		var params struct {
			Path   string `json:"path"`
			OldStr string `json:"old_str"`
			NewStr string `json:"new_str"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &params); err != nil {
			return "Error: " + err.Error()
		}

		content, err := os.ReadFile(params.Path)
		if err != nil {
			return "Error reading file: " + err.Error()
		}

		text := string(content)

		if params.OldStr == "" {
			text = text + params.NewStr
		} else {
			count := strings.Count(text, params.OldStr)
			if count == 0 {
				return fmt.Sprintf("Error: old_str not found in %s", params.Path)
			}
			if count > 1 {
				return fmt.Sprintf("Error: old_str found %d times in %s, must be unique", count, params.Path)
			}
			text = strings.Replace(text, params.OldStr, params.NewStr, 1)
		}

		if err := os.WriteFile(params.Path, []byte(text), 0644); err != nil {
			return "Error writing file: " + err.Error()
		}

		return fmt.Sprintf("Successfully edited %s", params.Path)
	},
}
