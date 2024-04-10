package tavily

import (
	"context"
	_ "github.com/joho/godotenv/autoload"
	"os"
	"testing"
)

func TestTavilyClient_Search(t *testing.T) {
	tavilyClient := NewClient(os.Getenv("TAVILY_API_KEY"))
	response, err := tavilyClient.Search(context.Background(), "What is GitHub?")
	if err != nil {
		t.Fatal("tavilyClient Search", err)
	}
	t.Log(response)
}
