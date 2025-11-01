package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/gocolly/colly/v2"
)

// PageData represents the scraped data structure
type PageData struct {
	URL  string `json:"url"`
	Text string `json:"text"`
}

// isValidURL checks if the URL is valid and uses HTTP/HTTPS scheme
func isValidURL(u string) bool {
	parsed, err := url.Parse(u)
	if err != nil {
		return false
	}
	return parsed.Scheme == "http" || parsed.Scheme == "https"
}

// validateURLs checks if the URL list is not empty and contains valid URLs
func validateURLs(urls []string) error {
	if len(urls) == 0 {
		return fmt.Errorf("URL list is empty")
	}

	for _, u := range urls {
		if !isValidURL(u) {
			return fmt.Errorf("invalid URL: %s", u)
		}
	}
	return nil
}

// validateOutputFile checks if the output file exists and is not empty
func validateOutputFile(filename string) error {
	info, err := os.Stat(filename)
	if err != nil {
		return fmt.Errorf("output file does not exist: %v", err)
	}
	if info.Size() == 0 {
		return fmt.Errorf("output file is empty")
	}
	return nil
}

// validateJSONLStructure validates the structure of the .jl file
func validateJSONLStructure(filename string) error {
	file, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	decoder := json.NewDecoder(file)
	lineNumber := 0

	for decoder.More() {
		var data PageData
		if err := decoder.Decode(&data); err != nil {
			return fmt.Errorf("invalid JSON at line %d: %v", lineNumber+1, err)
		}
		if data.URL == "" || data.Text == "" {
			return fmt.Errorf("missing required fields at line %d", lineNumber+1)
		}
		lineNumber++
	}

	return nil
}

func main() {
	// Define Wikipedia URLs to crawl
	urls := []string{
		"https://en.wikipedia.org/wiki/Robotics",
		"https://en.wikipedia.org/wiki/Robot",
		"https://en.wikipedia.org/wiki/Reinforcement_learning",
		"https://en.wikipedia.org/wiki/Robot_Operating_System",
		"https://en.wikipedia.org/wiki/Intelligent_agent",
		"https://en.wikipedia.org/wiki/Software_agent",
		"https://en.wikipedia.org/wiki/Robotic_process_automation",
		"https://en.wikipedia.org/wiki/Chatbot",
		"https://en.wikipedia.org/wiki/Applications_of_artificial_intelligence",
		"https://en.wikipedia.org/wiki/Android_(robot)",
	}

	// Pre-crawl validation
	if err := validateURLs(urls); err != nil {
		log.Fatalf("URL validation failed: %v", err)
	}

	// Output filename
	outputFile := "scraped_data.jl"

	// Remove existing output file
	os.Remove(outputFile)

	// Initialize collector with concurrent settings
	c := colly.NewCollector(
		colly.Async(true),
		colly.UserAgent("Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36"),
	)

	// Set parallelism and add random delay to be respectful
	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 2,
		RandomDelay: 1 * time.Second,
	})

	// Slice to store scraped data
	var scrapedData []PageData

	// Extract text from paragraph tags
	c.OnHTML("p", func(e *colly.HTMLElement) {
		// Get clean text without HTML tags
		text := strings.TrimSpace(e.Text)
		if text != "" {
			currentURL := e.Request.URL.String()

			// Find or create entry for current URL
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

	c.OnError(func(r *colly.Response, err error) {
		log.Printf("Request URL: %s failed with response: %v\n", r.Request.URL, r)
		log.Printf("Error: %v", err)
	})

	// Start timing
	startTime := time.Now()

	// Add URLs to the queue
	for _, u := range urls {
		if err := c.Visit(u); err != nil {
			log.Printf("Failed to visit %s: %v", u, err)
		}
	}

	// Wait for all threads to finish
	c.Wait()

	// Write results to JSON Lines file
	file, err := os.Create(outputFile)
	if err != nil {
		log.Fatalf("Cannot create output file: %v", err)
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	for _, data := range scrapedData {
		if err := encoder.Encode(data); err != nil {
			log.Printf("Error encoding data for URL %s: %v", data.URL, err)
		}
	}

	// Calculate and display runtime
	elapsed := time.Since(startTime)
	fmt.Printf("Scraping completed in %v\n", elapsed)

	// Post-crawl validations
	if err := validateOutputFile(outputFile); err != nil {
		log.Fatalf("Output file validation failed: %v", err)
	}

	if err := validateJSONLStructure(outputFile); err != nil {
		log.Fatalf("JSONL structure validation failed: %v", err)
	}

	fmt.Printf("Successfully scraped %d pages and saved to %s\n", len(scrapedData), outputFile)
}
