package tool

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(string) (string, error)
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
	Run: func(argsJSON string) (string, error) {
		var args struct {
			Command          string `json:"command"`
			WorkingDirectory string `json:"working_directory"`
			TimeoutSeconds   int    `json:"timeout_seconds"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}

		if args.TimeoutSeconds <= 0 {
			args.TimeoutSeconds = 30
		}
		if args.TimeoutSeconds > 120 {
			args.TimeoutSeconds = 120
		}

		cmd := exec.Command("sh", "-c", args.Command)
		if args.WorkingDirectory != "" {
			cmd.Dir = args.WorkingDirectory
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
		case <-time.After(time.Duration(args.TimeoutSeconds) * time.Second):
			cmd.Process.Kill()
			return "", fmt.Errorf("command timed out after %d seconds", args.TimeoutSeconds)
		}

		result := strings.TrimRight(string(output), "\n")
		if cmdErr != nil {
			return result, cmdErr
		}
		return result, nil
	},
}
