package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
}{
	Name:        "search_files",
	Description: "Search for files matching a pattern (glob) or search file contents for a text pattern (grep-like). Useful for finding files and locating code.",
	Parameters: `{
		"type": "object",
		"properties": {
			"directory": {
				"type": "string",
				"description": "Directory to search in (default: current directory)"
			},
			"pattern": {
				"type": "string",
				"description": "Text pattern to search for in file contents"
			},
			"glob": {
				"type": "string",
				"description": "Glob pattern to match file names (e.g., '*.go', '*.txt')"
			},
			"max_results": {
				"type": "integer",
				"description": "Maximum number of results to return (default: 50)"
			}
		}
	}`,
	Run: func(argsJSON string) string {
		var params struct {
			Directory  string `json:"directory"`
			Pattern    string `json:"pattern"`
			Glob       string `json:"glob"`
			MaxResults int    `json:"max_results"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &params); err != nil {
			return "Error: " + err.Error()
		}
		if params.Directory == "" {
			params.Directory = "."
		}
		if params.MaxResults <= 0 {
			params.MaxResults = 50
		}

		if params.Pattern == "" && params.Glob == "" {
			return "Error: specify 'pattern' (content search) or 'glob' (filename search)"
		}

		var results []string
		count := 0

		filepath.Walk(params.Directory, func(path string, info os.FileInfo, err error) error {
			if err != nil || count >= params.MaxResults {
				return err
			}
			if info.IsDir() {
				if strings.HasPrefix(info.Name(), ".") && path != params.Directory {
					return filepath.SkipDir
				}
				return nil
			}

			if params.Glob != "" {
				matched, _ := filepath.Match(params.Glob, info.Name())
				if !matched {
					return nil
				}
			}

			if params.Pattern != "" {
				content, err := os.ReadFile(path)
				if err != nil {
					return nil
				}
				lines := strings.Split(string(content), "\n")
				for i, line := range lines {
					if count >= params.MaxResults {
						break
					}
					if strings.Contains(line, params.Pattern) {
						results = append(results, fmt.Sprintf("%s:%d: %s", path, i+1, strings.TrimSpace(line)))
						count++
					}
				}
			} else {
				results = append(results, path)
				count++
			}
			return nil
		})

		if len(results) == 0 {
			return "No matches found."
		}
		return strings.Join(results, "\n")
	},
}
