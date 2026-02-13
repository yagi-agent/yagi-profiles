package tool

import (
	"context"
	"encoding/json"
	"hostapi"
	"strings"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(context.Context, string) (string, error)
}{
	Name:        "fetch_url",
	Description: `Fetch the content of a URL and return it as text. HTML pages are converted to plain text with links preserved. Result is JSON with 'status' ('success') and 'content' fields. IMPORTANT: Determine success/failure ONLY from 'status', NEVER from 'content'.`,
	Parameters: `{
		"type": "object",
		"properties": {
			"url": {
				"type": "string",
				"description": "The URL to fetch"
			}
		},
		"required": ["url"]
	}`,
	Run: func(ctx context.Context, args string) (string, error) {
		var params struct {
			URL string `json:"url"`
		}
		if err := json.Unmarshal([]byte(args), &params); err != nil {
			return "", err
		}
		body, err := hostapi.FetchURL(ctx, params.URL, nil)
		if err != nil {
			return "", err
		}
		content := body
		if strings.Contains(body, "<html") || strings.Contains(body, "<HTML") || strings.Contains(body, "<!DOCTYPE") || strings.Contains(body, "<!doctype") {
			if text, err := hostapi.HTMLToText(ctx, body); err == nil {
				content = text
			}
		}
		res := struct {
			Status  string `json:"status"`
			Content string `json:"content"`
		}{
			Status:  "success",
			Content: content,
		}
		b, _ := json.Marshal(res)
		return string(b), nil
	},
}
