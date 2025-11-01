package main

import (
	"os"
	"strings"
	"testing"
	"time"

	"github.com/gocolly/colly/v2"
)

// TestValidateURLs tests URL validation functionality
func TestValidateURLs(t *testing.T) {
	tests := []struct {
		name    string
		urls    []string
		wantErr bool
	}{
		{
			name:    "valid URLs",
			urls:    []string{"https://en.wikipedia.org/wiki/Robotics", "https://example.com"},
			wantErr: false,
		},
		{
			name:    "empty URL list",
			urls:    []string{},
			wantErr: true,
		},
		{
			name:    "invalid URL",
			urls:    []string{"invalid-url"},
			wantErr: true,
		},
		{
			name:    "invalid scheme",
			urls:    []string{"ftp://example.com"},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validateURLs(tt.urls)
			if (err != nil) != tt.wantErr {
				t.Errorf("validateURLs() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// TestValidateOutputFile tests output file validation
func TestValidateOutputFile(t *testing.T) {
	// Create a temporary file for testing
	tmpFile, err := os.CreateTemp("", "test*.jl")
	if err != nil {
		t.Fatalf("Cannot create temp file: %v", err)
	}
	defer os.Remove(tmpFile.Name())

	// Test with non-existent file
	if err := validateOutputFile("nonexistent.jl"); err == nil {
		t.Error("Expected error for non-existent file, got nil")
	}

	// Test with empty file
	if err := validateOutputFile(tmpFile.Name()); err == nil {
		t.Error("Expected error for empty file, got nil")
	}

	// Test with valid file
	content := `{"url":"https://example.com","text":"test content"}`
	if _, err := tmpFile.WriteString(content); err != nil {
		t.Fatalf("Cannot write to temp file: %v", err)
	}
	tmpFile.Close()

	if err := validateOutputFile(tmpFile.Name()); err != nil {
		t.Errorf("Unexpected error for valid file: %v", err)
	}
}

// TestValidateJSONLStructure tests JSONL file structure validation
func TestValidateJSONLStructure(t *testing.T) {
	// Create valid JSONL file
	validFile, err := os.CreateTemp("", "valid*.jl")
	if err != nil {
		t.Fatalf("Cannot create temp file: %v", err)
	}
	defer os.Remove(validFile.Name())

	validContent := `{"url":"https://example.com","text":"test content"}
{"url":"https://example.org","text":"more content"}`
	if _, err := validFile.WriteString(validContent); err != nil {
		t.Fatalf("Cannot write to temp file: %v", err)
	}
	validFile.Close()

	if err := validateJSONLStructure(validFile.Name()); err != nil {
		t.Errorf("Unexpected error for valid JSONL: %v", err)
	}

	// Create invalid JSONL file
	invalidFile, err := os.CreateTemp("", "invalid*.jl")
	if err != nil {
		t.Fatalf("Cannot create temp file: %v", err)
	}
	defer os.Remove(invalidFile.Name())

	invalidContent := `{"url":"https://example.com","text":"test content"}
{"url":"","text":"missing url"}
invalid json`
	if _, err := invalidFile.WriteString(invalidContent); err != nil {
		t.Fatalf("Cannot write to temp file: %v", err)
	}
	invalidFile.Close()

	if err := validateJSONLStructure(invalidFile.Name()); err == nil {
		t.Error("Expected error for invalid JSONL, got nil")
	}
}

// BenchmarkCrawlerPerformance benchmarks the crawling performance
func BenchmarkCrawlerPerformance(b *testing.B) {
	urls := []string{
		"https://en.wikipedia.org/wiki/Robotics",
		"https://en.wikipedia.org/wiki/Robot",
	}

	for i := 0; i < b.N; i++ {
		start := time.Now()

		c := colly.NewCollector(colly.Async(true))
		c.Limit(&colly.LimitRule{
			DomainGlob:  "*",
			Parallelism: 2,
			RandomDelay: 500 * time.Millisecond,
		})

		var scrapedData []PageData

		c.OnHTML("p", func(e *colly.HTMLElement) {
			text := strings.TrimSpace(e.Text)
			if text != "" {
				currentURL := e.Request.URL.String()
				var found bool
				for i, data := range scrapedData {
					if data.URL == currentURL {
						scrapedData[i].Text += " " + text
						found = true
						break
					}
				}
				if !found {
					scrapedData = append(scrapedData, PageData{
						URL:  currentURL,
						Text: text,
					})
				}
			}
		})

		for _, u := range urls {
			c.Visit(u)
		}
		c.Wait()

		elapsed := time.Since(start)
		b.Logf("Run %d completed in %v nanoseconds", i+1, elapsed.Nanoseconds())
	}
}
