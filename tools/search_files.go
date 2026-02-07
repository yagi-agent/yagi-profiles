package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type SearchFilesArgs struct {
	Directory string `json:"directory"`
	Pattern   string `json:"pattern"`
	Glob      string `json:"glob"`
}

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
}{
	Name:        "search_files",
	Description: "Search files by glob pattern or text pattern",
	Parameters: `{
		"type": "object",
		"properties": {
			"directory": {
				"type": "string",
				"description": "Directory to search in"
			},
			"pattern": {
				"type": "string",
				"description": "Text pattern to search for in file contents"
			},
			"glob": {
				"type": "string",
				"description": "Glob pattern to match file names"
			}
		},
		"required": ["directory"]
	}`,
	Run: func(argsJSON string) string {
		var args SearchFilesArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "Error: " + err.Error()
		}
		if args.Pattern == "" && args.Glob == "" {
			return "Error: either pattern or glob must be specified"
		}
		if args.Glob != "" {
			var results []string
			filepath.Walk(args.Directory, func(path string, info os.FileInfo, err error) error {
				if err != nil {
					return nil
				}
				if info.IsDir() {
					return nil
				}
				matched, _ := filepath.Match(args.Glob, info.Name())
				if matched {
					results = append(results, path)
				}
				return nil
			})
			return strings.Join(results, "\n")
		}
		var results []string
		filepath.Walk(args.Directory, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			f, err := os.Open(path)
			if err != nil {
				return nil
			}
			defer f.Close()
			scanner := bufio.NewScanner(f)
			lineNum := 0
			for scanner.Scan() {
				lineNum++
				line := scanner.Text()
				if strings.Contains(line, args.Pattern) {
					results = append(results, fmt.Sprintf("%s:%d:%s", path, lineNum, line))
				}
			}
			return nil
		})
		if len(results) == 0 {
			return "No matches found"
		}
		return strings.Join(results, "\n")
	},
}
