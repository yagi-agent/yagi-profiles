package tool

import (
	"encoding/json"
	"fmt"
	"hostapi"
)

var Name = "nostr_timeline"
var Description = "Fetch the latest posts (kind:1 text notes) from the Nostr relay wss://yabu.me"
var Parameters = `{
	"type": "object",
	"properties": {
		"limit": {
			"type": "integer",
			"description": "Number of posts to fetch (default: 10, max: 50)"
		}
	}
}`

func Run(args string) (string, error) {
	var params struct {
		Limit int `json:"limit"`
	}
	if err := json.Unmarshal([]byte(args), &params); err != nil || params.Limit <= 0 {
		params.Limit = 10
	}
	if params.Limit > 50 {
		params.Limit = 50
	}

	req := []interface{}{
		"REQ",
		"yagi-sub",
		map[string]interface{}{
			"kinds": []int{1},
			"limit": params.Limit,
		},
	}
	reqJSON, _ := json.Marshal(req)

	raw := hostapi.WebSocketSend("wss://yabu.me", string(reqJSON), params.Limit+1, 10)

	var messages []string
	json.Unmarshal([]byte(raw), &messages)

	var result string
	count := 0
	for _, msg := range messages {
		var envelope []json.RawMessage
		if err := json.Unmarshal([]byte(msg), &envelope); err != nil || len(envelope) < 2 {
			continue
		}
		var tag string
		json.Unmarshal(envelope[0], &tag)
		if tag != "EVENT" || len(envelope) < 3 {
			continue
		}
		var ev struct {
			PubKey    string `json:"pubkey"`
			Content   string `json:"content"`
			CreatedAt int64  `json:"created_at"`
		}
		if err := json.Unmarshal(envelope[2], &ev); err != nil {
			continue
		}
		count++
		content := ev.Content
		if len(content) > 200 {
			content = content[:200] + "..."
		}
		result += fmt.Sprintf("[%d] (by %s...)\n%s\n\n", count, ev.PubKey[:12], content)
	}
	if result == "" {
		return "No events found.", nil
	}
	return result, nil
}
