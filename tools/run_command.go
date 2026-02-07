package main

import (
	"encoding/json"
	"os/exec"
	"strings"
)

type RunCommandArgs struct {
	Command          string `json:"command"`
	WorkingDirectory string `json:"working_directory"`
}

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
}{
	Name:        "run_command",
	Description: "Run a shell command and return its output",
	Parameters: `{
		"type": "object",
		"properties": {
			"command": {
				"type": "string",
				"description": "The command to run"
			},
			"working_directory": {
				"type": "string",
				"description": "Working directory for the command"
			}
		},
		"required": ["command"]
	}`,
	Run: func(argsJSON string) string {
		var args RunCommandArgs
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "Error: " + err.Error()
		}
		cmd := exec.Command("sh", "-c", args.Command)
		if args.WorkingDirectory != "" {
			cmd.Dir = args.WorkingDirectory
		}
		output, err := cmd.CombinedOutput()
		if err != nil {
			return "Exit error: " + err.Error() + "\n" + string(output)
		}
		return strings.TrimRight(string(output), "\n")
	},
}
