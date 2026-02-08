package tool

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) (string, error)
}{
	Name:        "glob",
	Description: "Find files by glob pattern across a directory tree. Supports '**' for recursive matching (e.g., '**/*.go', 'src/**/*.ts', '**/*test*'). Returns matching file paths sorted by modification time (newest first).",
	Parameters: `{
		"type": "object",
		"properties": {
			"pattern": {
				"type": "string",
				"description": "Glob pattern to match files. Use ** for recursive directory matching. Examples: **/*.go, src/**/*.ts, **/*test*, *.json"
			},
			"directory": {
				"type": "string",
				"description": "Root directory to search in (default: current directory)"
			},
			"limit": {
				"type": "integer",
				"description": "Maximum number of results to return (default: 100)"
			}
		},
		"required": ["pattern"]
	}`,
	Run: func(argsJSON string) (string, error) {
		var args struct {
			Pattern   string `json:"pattern"`
			Directory string `json:"directory"`
			Limit     int    `json:"limit"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}

		if args.Pattern == "" {
			return "", fmt.Errorf("pattern is required")
		}
		if args.Directory == "" {
			args.Directory = "."
		}
		if args.Limit <= 0 {
			args.Limit = 100
		}

		type fileEntry struct {
			path    string
			modTime int64
		}
		var matches []fileEntry

		segments := splitPattern(args.Pattern)

		filepath.Walk(args.Directory, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil
			}
			if info.IsDir() {
				name := info.Name()
				if strings.HasPrefix(name, ".") && path != args.Directory {
					return filepath.SkipDir
				}
				if name == "node_modules" || name == "vendor" {
					return filepath.SkipDir
				}
				return nil
			}

			rel, err := filepath.Rel(args.Directory, path)
			if err != nil {
				return nil
			}
			rel = filepath.ToSlash(rel)

			if matchGlob(segments, strings.Split(rel, "/")) {
				matches = append(matches, fileEntry{
					path:    rel,
					modTime: info.ModTime().UnixNano(),
				})
			}
			return nil
		})

		sort.Slice(matches, func(i, j int) bool {
			return matches[i].modTime > matches[j].modTime
		})

		if len(matches) > args.Limit {
			matches = matches[:args.Limit]
		}

		if len(matches) == 0 {
			return "No files found", nil
		}

		var sb strings.Builder
		for i, m := range matches {
			if i > 0 {
				sb.WriteByte('\n')
			}
			sb.WriteString(m.path)
		}
		return sb.String(), nil
	},
}

func splitPattern(pattern string) []string {
	return strings.Split(filepath.ToSlash(pattern), "/")
}

func matchGlob(patternSegments, pathSegments []string) bool {
	return doMatch(patternSegments, pathSegments)
}

func doMatch(pat, path []string) bool {
	for len(pat) > 0 {
		if pat[0] == "**" {
			pat = pat[1:]
			if len(pat) == 0 {
				return true
			}
			for i := 0; i <= len(path); i++ {
				if doMatch(pat, path[i:]) {
					return true
				}
			}
			return false
		}

		if len(path) == 0 {
			return false
		}

		matched, _ := filepath.Match(pat[0], path[0])
		if !matched {
			return false
		}
		pat = pat[1:]
		path = path[1:]
	}
	return len(path) == 0
}
