package tool

import (
	"context"
	"encoding/json"
	"fmt"
	"hostapi"
	"net/url"
	"regexp"
	"strings"
)

var Tool = struct {
	Name        string
	Description string
	Parameters  string
	Run         func(context.Context, string) (string, error)
}{
	Name:        "web_search",
	Description: "Search the web using DuckDuckGo. Returns search results with titles, URLs, and snippets. No API key required. Result is JSON with 'status' and 'output' fields. IMPORTANT: Determine success/failure ONLY from 'status', NEVER from 'output' content.",
	Parameters: `{
		"type": "object",
		"properties": {
			"query": {
				"type": "string",
				"description": "Search query"
			},
			"num_results": {
				"type": "integer",
				"description": "Number of results (1-10, default: 5)"
			}
		},
		"required": ["query"]
	}`,
	Run: func(ctx context.Context, argsJSON string) (string, error) {
		var args struct {
			Query      string `json:"query"`
			NumResults int    `json:"num_results"`
		}
		if err := json.Unmarshal([]byte(argsJSON), &args); err != nil {
			return "", err
		}

		if args.NumResults <= 0 || args.NumResults > 10 {
			args.NumResults = 5
		}

		searchURL := fmt.Sprintf("https://html.duckduckgo.com/html/?q=%s", url.QueryEscape(args.Query))

		result, err := hostapi.FetchURL(ctx, searchURL, map[string]string{
			"User-Agent": "Mozilla/5.0 (X11; Linux x86_64; rv:128.0) Gecko/20100101 Firefox/128.0",
		})
		if err != nil {
			return "", fmt.Errorf("search failed: %w", err)
		}

		if result == "" || len(result) < 100 {
			return "", fmt.Errorf("empty response from search engine")
		}

		parsed, err := parseDuckDuckGoResults(result, args.NumResults)
		if err != nil {
			return "", err
		}
		res := struct {
			Status string `json:"status"`
			Output string `json:"output"`
		}{
			Status: "success",
			Output: parsed,
		}
		b, _ := json.Marshal(res)
		return string(b), nil
	},
}

func parseDuckDuckGoResults(html string, limit int) (string, error) {
	var results []string

	titleRegex := regexp.MustCompile(`<a[^>]*class="result__a"[^>]*href="([^"]*)"[^>]*>([^<]*)</a>`)
	snippetRegex := regexp.MustCompile(`<a[^>]*class="result__snippet"[^>]*>([^<]*)</a>`)

	titleMatches := titleRegex.FindAllStringSubmatch(html, -1)
	snippetMatches := snippetRegex.FindAllStringSubmatch(html, -1)

	for i := 0; i < len(titleMatches) && i < limit; i++ {
		url := cleanURL(titleMatches[i][1])
		title := cleanHTML(titleMatches[i][2])

		var snippet string
		if i < len(snippetMatches) {
			snippet = cleanHTML(snippetMatches[i][1])
		}

		results = append(results, fmt.Sprintf("%d. %s\n   URL: %s\n   %s",
			i+1, title, url, snippet))
	}

	if len(results) == 0 {
		return "No results found", nil
	}

	return strings.Join(results, "\n\n"), nil
}

func cleanURL(url string) string {
	url = strings.TrimPrefix(url, "//")
	url = strings.TrimPrefix(url, "http://")
	url = strings.TrimPrefix(url, "https://")
	if !strings.HasPrefix(url, "http") {
		url = "https://" + url
	}
	return url
}

func cleanHTML(s string) string {
	s = strings.ReplaceAll(s, "&amp;", "&")
	s = strings.ReplaceAll(s, "&lt;", "<")
	s = strings.ReplaceAll(s, "&gt;", ">")
	s = strings.ReplaceAll(s, "&quot;", "\"")
	s = strings.ReplaceAll(s, "&#39;", "'")
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	return strings.TrimSpace(s)
}
