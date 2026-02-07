package tool

import (
"encoding/json"
"os"
)

var Tool = struct {
Name        string
Description string
Parameters  string
Run         func(string) (string, error)
}{
Name:        "make_directory",
Description: "Create a directory (including parent directories if needed)",
Parameters: `{
"type": "object",
"properties": {
"path": {
"type": "string",
"description": "Path of the directory to create"
}
},
"required": ["path"]
}`,
Run: func(argsJSON string) (string, error) {
var args struct {
Path string `json:"path"`
}
if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
return "", err
}
if err := os.MkdirAll(args.Path, 0755); err != nil {
return "", err
}
return "Created directory: " + args.Path, nil
},
}
