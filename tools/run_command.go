package main

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) string
}{
	Name:        "run_command",
	Description: "Execute a shell command and return its output (stdout and stderr). Use this to run programs, scripts, build tools, git commands, etc.",
	Parameters: `{
		"type": "object",
		"properties": {
			"command": {
				"type": "string",
				"description": "The shell command to execute"
			},
			"working_directory": {
				"type": "string",
				"description": "Working directory for the command (default: current directory)"
			},
			"timeout_seconds": {
				"type": "integer",
				"description": "Timeout in seconds (default: 30, max: 120)"
			}
		},
		"required": ["command"]
	}`,
	Run: func(argsJSON string) string {
		var params struct {
			Command          string `json:"command"`
			WorkingDirectory string `json:"working_directory"`
			TimeoutSeconds   int    `json:"timeout_seconds"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &params); err != nil {
			return "Error: " + err.Error()
		}

		if params.TimeoutSeconds <= 0 {
			params.TimeoutSeconds = 30
		}
		if params.TimeoutSeconds > 120 {
			params.TimeoutSeconds = 120
		}

		cmd := exec.Command("sh", "-c", params.Command)
		if params.WorkingDirectory != "" {
			cmd.Dir = params.WorkingDirectory
		}

		done := make(chan struct{})
		var output []byte
		var cmdErr error
		go func() {
			output, cmdErr = cmd.CombinedOutput()
			close(done)
		}()

		select {
		case <-done:
		case <-time.After(time.Duration(params.TimeoutSeconds) * time.Second):
			cmd.Process.Kill()
			return fmt.Sprintf("Error: command timed out after %d seconds", params.TimeoutSeconds)
		}

		result := strings.TrimRight(string(output), "\n")
		if cmdErr != nil {
			return fmt.Sprintf("Exit error: %v\n%s", cmdErr, result)
		}
		return result
	},
}
