package duckduckgo

import (
	"encoding/json"
	"fmt"
	"github.com/r4stl1n/mbop/pkg/tools"
	"go.uber.org/zap"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

type DuckDuckGo struct {
}

func (d DuckDuckGo) Name() string {
	return "duckduckgo"
}

func (d DuckDuckGo) Example() string {
	return "duckduckgo: golang tutorial\n   or\nduckduckgo: query=\"golang tutorial\" format=json"
}

func (d DuckDuckGo) Description() string {
	return "Searches DuckDuckGo and returns the top search results. Parameters: query (required), format (optional, default: json)"
}

func (d DuckDuckGo) Run(values ...string) (string, error) {
	if len(values) != 1 {
		return "", fmt.Errorf("expected one argument (search query or parameters)")
	}

	// Parse parameters using the new parameter parsing mechanism
	params, err := tools.ParseParameters(values[0], "query")
	if err != nil {
		return "", fmt.Errorf("failed to parse parameters: %v", err)
	}

	// Get the query parameter, which is required
	query, err := tools.GetParameterValue(params, "query", "", true)
	if err != nil {
		return "", err
	}

	// Get optional parameters with default values
	format, _ := tools.GetParameterValue(params, "format", "json", false)

	searchURL := fmt.Sprintf("https://api.duckduckgo.com/?q=%s&format=%s", url.QueryEscape(query), url.QueryEscape(format))

	// Create a new request
	client := &http.Client{}
	req, err := http.NewRequest("GET", searchURL, nil)
	if err != nil {
		zap.L().Error("http request creation failed", 
			zap.String("query", query),
			zap.String("url", searchURL), 
			zap.Error(err))
		return "", err
	}

	// Set headers to mimic a browser request
	req.Header.Set("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/91.0.4472.124 Safari/537.36")

	// Send the request
	resp, err := client.Do(req)
	if err != nil {
		zap.L().Error("http request failed", 
			zap.String("query", query),
			zap.String("url", searchURL), 
			zap.Error(err))
		return "", err
	}

	defer func(Body io.ReadCloser) {
		bodyCloseError := Body.Close()
		if bodyCloseError != nil {
			zap.L().Error("response body close failed", zap.Error(bodyCloseError))
		}
	}(resp.Body)

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		zap.L().Error("response body read failed", 
			zap.String("query", query),
			zap.Error(err))
		return "", err
	}

	// Parse the JSON response
	var jsonResponse DuckDuckGoJsonResponse
	err = json.Unmarshal(body, &jsonResponse)
	if err != nil {
		zap.L().Error("JSON response parsing failed", 
			zap.String("query", query),
			zap.Error(err))
		return "", err
	}

	// Convert the JSON response to our standard response format
	response := DuckDuckGoResponse{
		Results: []DuckDuckGoResult{},
	}

	// Add the abstract if available
	if jsonResponse.AbstractText != "" {
		result := DuckDuckGoResult{
			Title:       jsonResponse.Heading,
			URL:         jsonResponse.AbstractURL,
			Description: jsonResponse.AbstractText,
		}
		response.Results = append(response.Results, result)
	}

	// Add related topics
	for _, topic := range jsonResponse.RelatedTopics {
		if topic.Text != "" {
			result := DuckDuckGoResult{
				Title:       topic.Text,
				URL:         topic.FirstURL,
				Description: topic.Text,
			}
			response.Results = append(response.Results, result)
		}

		// Also add nested topics if available
		for _, nestedTopic := range topic.Topics {
			result := DuckDuckGoResult{
				Title:       nestedTopic.Text,
				URL:         nestedTopic.FirstURL,
				Description: nestedTopic.Text,
			}
			response.Results = append(response.Results, result)
		}
	}

	if len(response.Results) == 0 {
		return "No results found for: " + query, nil
	}

	// Format the results as a string
	var formattedResults strings.Builder
	formattedResults.WriteString(fmt.Sprintf("Search results for: %s\n\n", query))

	for i, result := range response.Results {
		formattedResults.WriteString(fmt.Sprintf("%d. %s\n", i+1, result.Title))
		if result.URL != "" {
			formattedResults.WriteString(fmt.Sprintf("   URL: %s\n", result.URL))
		}
		formattedResults.WriteString(fmt.Sprintf("   %s\n\n", result.Description))
	}

	return formattedResults.String(), nil
}

// cleanHTML removes HTML tags and decodes HTML entities
func cleanHTML(html string) string {
	// Remove HTML tags
	tagRegex := regexp.MustCompile(`<[^>]*>`)
	text := tagRegex.ReplaceAllString(html, "")

	// Replace common HTML entities
	text = strings.ReplaceAll(text, "&amp;", "&")
	text = strings.ReplaceAll(text, "&lt;", "<")
	text = strings.ReplaceAll(text, "&gt;", ">")
	text = strings.ReplaceAll(text, "&quot;", "\"")
	text = strings.ReplaceAll(text, "&#39;", "'")

	// Trim whitespace
	return strings.TrimSpace(text)
}
