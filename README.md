# 🧠 Go and Python Web Crawling Comparison – Intelligent Systems and Robotics

## 📘 Overview
This project explores the development of a **concurrent web crawler and scraper in Go** using the **Colly** framework to extract text from Wikipedia pages related to **intelligent systems** and **robotics**.  
The assignment replicates the functionality of a **Python/Scrapy-based crawler**, comparing their design efficiency and runtime performance.

---

## 📂 Folder Structure
```
C:.
|   go.mod
|   go.sum
|   go_wiki_reader.go
|   go_wiki_reader_test.go
|   scraped_data.jl
|   structure.txt
|   
\---WebFocusedCrawlWorkV001
    +---WebFocusedCrawlWorkV001
    |   |   items.jl
    |   |   run-articles-spider.py
    |   |   scrapy.cfg
    |   |
    |   +---WebFocusedCrawl
    |   |   |   items.py
    |   |   |   middlewares.py
    |   |   |   pipelines.py
    |   |   |   settings.py
    |   |   |
    |   |   +---spiders
    |   |   |       articles-spider.py
    |   |   |
    |   |   \---__pycache__
    |   \---wikipages
    |           (contains downloaded HTML files)
```

---

## 🚀 Go Portion – `go_wiki_reader.go`
The **Go-based crawler** (`go_wiki_reader.go`) uses the **Colly** framework for asynchronous web scraping.  
It:
- Validates the Wikipedia URL list.
- Concurrently crawls and scrapes text from each page (`<p>` tags only).
- Writes clean text results into a **JSON Lines file (`scraped_data.jl`)**.
- Performs file and data structure validations.
- Logs total crawl runtime.

The **test file (`go_wiki_reader_test.go`)** includes:
- Unit tests for URL, output file, and JSON structure validation.
- A benchmark function measuring crawling performance across runs.

**Result File:**  
`scraped_data.jl` — each line represents one Wikipedia page’s extracted text and URL.

---

## 🐍 Python Portion – `WebFocusedCrawlWorkV001`
The **Python/Scrapy crawler** (originally provided as a class reference) was **slightly modified** to:
- Record total runtime using `time.time()`.
- Validate and report the output JSON Lines file (`items.jl`), including file size and line count.

These enhancements make the Python implementation directly comparable to the Go version in terms of timing and output validation.

---

## ⏱️ Performance Comparison

| Language | Total Crawl Time | Speed Notes |
|-----------|------------------|--------------|
| **Go (Colly)** | **3.79 seconds** | Uses concurrency (goroutines) for parallel page requests. |
| **Python (Scrapy)** | **14.35 seconds** | Sequential processing, slower with larger page sets. |

➡️ **Go’s concurrency results in a crawl nearly 4× faster** than Python’s sequential Scrapy implementation.

---

## 🧩 Summary of Workflow
1. **Go-based Crawler:** Concurrently extracts Wikipedia text → validates output → generates `scraped_data.jl`.
2. **Python Scrapy Crawler:** Sequentially extracts full HTML → writes to `items.jl` → reports runtime and validation results.
3. **Comparison:** Go demonstrates improved runtime efficiency with similar scraping accuracy.

---

## 🧪 How to Run

### Run the Go crawler
```bash
go run go_wiki_reader.go
```
Output:
```
Scraping completed in 3.79s
Successfully scraped 10 pages and saved to scraped_data.jl
```

### Run Go tests and benchmarks
```bash
go test ./...
go test -bench=.
```

### Run the Python Scrapy crawler
From inside the `WebFocusedCrawlWorkV001` folder:
```bash
python run-articles-spider.py
```

---

## 🤖 GenAI Tools
**ChatGPT** – Used for **planning and guidance**, including project structure, input/output validation strategy, and test design.  
**DeepSeek** – Used for **coding assistance**, including actual Go code creation and function implementation for scraping, validation, and benchmarking.  

---

## 📚 Deliverables
- `go_wiki_reader.go` — Go crawler implementation.  
- `go_wiki_reader_test.go` — Tests and benchmarks.  
- `scraped_data.jl` — Output JSON Lines file (Go).  
- `WebFocusedCrawlWorkV001` — Modified Python Scrapy crawler for comparison.  
- `structure.txt` — Folder structure record.
