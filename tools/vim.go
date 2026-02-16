package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"strings"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(context.Context, string) (string, error)
}{
	Name:        "vim",
	Description: "Control vim editor. Only available when running inside vim. Use this tool to read buffer, execute commands, insert text, etc.",
	Parameters: `{
		"type": "object",
		"properties": {
			"action": {
				"type": "string",
				"enum": ["get_buffer", "get_cursor", "execute", "insert", "search", "replace"],
				"description": "Action to perform"
			},
			"args": {
				"type": "object",
				"description": "Arguments for the action"
			}
		},
		"required": ["action"]
	}`,
	Run: func(ctx context.Context, argsJSON string) (string, error) {
		var args struct {
			Action string          `json:"action"`
			Args   json.RawMessage `json:"args"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return `{"status": "error", "message": "invalid arguments"}`, nil
		}

		switch args.Action {
		case "get_buffer":
			return handleGetBuffer(args.Args)
		case "get_cursor":
			return handleGetCursor(args.Args)
		case "execute":
			return handleExecute(args.Args)
		case "insert":
			return handleInsert(args.Args)
		case "search":
			return handleSearch(args.Args)
		case "replace":
			return handleReplace(args.Args)
		default:
			return `{"status": "error", "message": "unknown action: ` + args.Action + `"}`, nil
		}
	},
}

func handleGetBuffer(args json.RawMessage) (string, error) {
	path := ""
	if len(args) > 0 {
		var a struct {
			Path string `json:"path"`
		}
		if err := json.Unmarshal(args, &a); err == nil {
			path = a.Path
		}
	}

	if path == "" {
		path = os.Getenv("YAGI_VIM_BUFFER_PATH")
		if path == "" {
			return `{"status": "error", "message": "no buffer path available"}`, nil
		}
	}

	content, err := os.ReadFile(path)
	if err != nil {
		return `{"status": "error", "message": "` + err.Error() + `"}`, nil
	}

	res := map[string]interface{}{
		"status":  "success",
		"content": string(content),
		"path":    path,
	}
	b, _ := json.Marshal(res)
	return string(b), nil
}

func handleGetCursor(args json.RawMessage) (string, error) {
	line := os.Getenv("YAGI_VIM_CURSOR_LINE")
	col := os.Getenv("YAGI_VIM_CURSOR_COL")

	lineNum, _ := strconv.Atoi(line)
	colNum, _ := strconv.Atoi(col)

	res := map[string]interface{}{
		"status": "success",
		"line":   lineNum,
		"column": colNum,
	}
	b, _ := json.Marshal(res)
	return string(b), nil
}

func handleExecute(args json.RawMessage) (string, error) {
	var a struct {
		Command string `json:"command"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		return `{"status": "error", "message": "missing command argument"}`, nil
	}

	cmd := a.Command
	res := map[string]interface{}{
		"status":      "success",
		"vim_command": cmd,
	}
	b, _ := json.Marshal(res)
	return string(b), nil
}

func handleInsert(args json.RawMessage) (string, error) {
	var a struct {
		Text    string `json:"text"`
		After   bool   `json:"after"`
		NewLine bool   `json:"new_line"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		return `{"status": "error", "message": "missing text argument"}`, nil
	}

	var cmd string
	if a.NewLine {
		if a.After {
			cmd = "o" + a.Text + "\u001b"
		} else {
			cmd = "O" + a.Text + "\u001b"
		}
	} else {
		if a.After {
			cmd = "a" + a.Text + "\u001b"
		} else {
			cmd = "i" + a.Text + "\u001b"
		}
	}

	res := map[string]interface{}{
		"status":      "success",
		"vim_command": cmd,
	}
	b, _ := json.Marshal(res)
	return string(b), nil
}

func handleSearch(args json.RawMessage) (string, error) {
	var a struct {
		Pattern    string `json:"pattern"`
		Forward    bool   `json:"forward"`
		WrapAround bool   `json:"wrap_around"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		return `{"status": "error", "message": "missing pattern argument"}`, nil
	}

	escaped := strings.ReplaceAll(a.Pattern, "/", "\\/")

	cmd := "/"
	if !a.Forward {
		cmd = "?"
	}
	cmd += escaped + "\u000d"

	if a.WrapAround {
		cmd += ":set wrapscan\u000d"
	}

	res := map[string]interface{}{
		"status":      "success",
		"vim_command": cmd,
	}
	b, _ := json.Marshal(res)
	return string(b), nil
}

func handleReplace(args json.RawMessage) (string, error) {
	var a struct {
		Pattern string `json:"pattern"`
		Replace string `json:"replace"`
		Flags   string `json:"flags"`
	}
	if err := json.Unmarshal(args, &a); err != nil {
		return `{"status": "error", "message": "missing pattern or replace argument"}`, nil
	}

	pattern := strings.ReplaceAll(a.Pattern, "/", "\\/")
	replace := strings.ReplaceAll(a.Replace, "/", "\\/")

	flags := a.Flags
	if flags == "" {
		flags = "g"
	}

	cmd := fmt.Sprintf(":%s/%s/%s/%s\u000d", flags, pattern, replace, flags)

	res := map[string]interface{}{
		"status":      "success",
		"vim_command": cmd,
	}
	b, _ := json.Marshal(res)
	return string(b), nil
}
