package tavily

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const DefaultTavilyBaseURL = "https://api.tavily.com"

type Client struct {
	BaseURL    string
	HttpClient *http.Client
	APIKey     string
}

func NewClient(APIKey string) *Client {
	return &Client{
		BaseURL:    DefaultTavilyBaseURL,
		HttpClient: http.DefaultClient,
		APIKey:     APIKey,
	}
}

type SearchRequest struct {
	Query             string   `json:"query,omitempty"`
	APIKey            string   `json:"api_key,omitempty"`
	SearchDepth       string   `json:"search_depth,omitempty"`
	Topic             string   `json:"topic,omitempty"`
	IncludeAnswer     bool     `json:"include_answer,omitempty"`
	IncludeRawContent bool     `json:"include_raw_content,omitempty"`
	IncludeImages     bool     `json:"include_images,omitempty"`
	IncludeDomains    []string `json:"include_domains,omitempty"`
	MaxResults        int      `json:"max_results,omitempty"`
}

type SearchResult struct {
	Title      string  `json:"title"`
	URL        string  `json:"url"`
	Content    string  `json:"content"`
	Score      float64 `json:"score"`
	RawContent any     `json:"raw_content"`
}
type SearchResponse struct {
	Query             string         `json:"query"`
	FollowUpQuestions []string       `json:"follow_up_questions"`
	Answer            string         `json:"answer"`
	Images            []string       `json:"images"`
	Results           []SearchResult `json:"results"`
	ResponseTime      float64        `json:"response_time"`
}

// SearchOption is a function that modifies SearchRequest
type SearchOption func(*SearchRequest)

// WithIncludeAnswer sets the IncludeAnswer field
func WithIncludeAnswer(include bool) SearchOption {
	return func(sr *SearchRequest) {
		sr.IncludeAnswer = include
	}
}

// WithMaxResults sets the MaxResults field
func WithMaxResults(max int) SearchOption {
	return func(sr *SearchRequest) {
		sr.MaxResults = max
	}
}

// WithSearchDepth sets the SearchDepth field
func WithSearchDepth(depth string) SearchOption {
	return func(sr *SearchRequest) {
		sr.SearchDepth = depth
	}
}

// WithTopic sets the Topic field
func WithTopic(topic string) SearchOption {
	return func(sr *SearchRequest) {
		sr.Topic = topic
	}
}

// WithIncludeRawContent sets the IncludeRawContent field
func WithIncludeRawContent(include bool) SearchOption {
	return func(sr *SearchRequest) {
		sr.IncludeRawContent = include
	}
}

// WithIncludeImages sets the IncludeImages field
func WithIncludeImages(include bool) SearchOption {
	return func(sr *SearchRequest) {
		sr.IncludeImages = include
	}
}

// WithIncludeDomains sets the IncludeDomains field
func WithIncludeDomains(domains []string) SearchOption {
	return func(sr *SearchRequest) {
		sr.IncludeDomains = domains
	}
}

func (c *Client) SearchWithOptions(ctx context.Context, query string, opts ...SearchOption) (*SearchResponse, error) {
	searchRequest := SearchRequest{
		Query:  query,
		APIKey: c.APIKey,
	}

	for _, opt := range opts {
		opt(&searchRequest)
	}

	searchRequestJSON, err := json.Marshal(searchRequest)
	if err != nil {
		return nil, fmt.Errorf("tavily client search, marshal search request: %w", err)
	}

	request, err := http.NewRequestWithContext(ctx, http.MethodPost, c.BaseURL+"/search", bytes.NewReader(searchRequestJSON))
	if err != nil {
		return nil, fmt.Errorf("tavily client search, build request: %w", err)
	}
	response, err := c.HttpClient.Do(request)
	if err != nil {
		return nil, fmt.Errorf("tavily client search, call /search api: %w", err)
	}
	defer response.Body.Close()
	responseBody, err := io.ReadAll(response.Body)
	if err != nil {
		return nil, fmt.Errorf("tavily client search, read response body: %w", err)
	}

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("tavily client search, response status code: %d, response body: %s", response.StatusCode, string(responseBody))
	}

	result := SearchResponse{}
	err = json.Unmarshal(responseBody, &result)
	if err != nil {
		return nil, fmt.Errorf("tavily client search, parse response: %w", err)
	}

	return &result, nil
}

func (c *Client) Search(ctx context.Context, query string) (*SearchResponse, error) {
	return c.SearchWithOptions(ctx, query,
		WithIncludeAnswer(true),
		WithMaxResults(5),
	)
}
