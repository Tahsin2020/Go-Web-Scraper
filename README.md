# Go Web Scraper

A simple web scraping application written in **Go** that collects article titles and URLs from [Hacker News](https://news.ycombinator.com/).

The application uses Go's standard `net/http` package to send HTTP requests and the `goquery` library to parse HTML and extract information using CSS selectors. The collected data is displayed in the terminal and exported to a CSV file.

## Features

- **HTTP Requests:** Retrieves website content using HTTP GET requests.
- **HTML Parsing:** Parses HTML documents using the goquery library.
- **CSS Selectors:** Locates and extracts specific HTML elements.
- **Data Extraction:** Collects article titles and URLs.
- **CSV Export:** Saves scraped information into a CSV file.
- **Error Handling:** Handles HTTP request failures, invalid responses, and file operation errors.
- **Request Timeout:** Prevents the application from waiting indefinitely for a website response.

## Technologies Used

| Technology | Purpose |
|---|---|
| Go (Golang) | Main programming language |
| net/http | Sending HTTP requests |
| goquery | Parsing HTML and selecting elements |
| encoding/csv | Exporting scraped data |
| os | File and directory operations |

## Project Structure

```text
go-web-scraper/
│
├── main.go              # Main application
├── go.mod               # Go module configuration
├── go.sum               # Dependency checksums
│
└── output/
    └── articles.csv     # Scraped article data
```

## Requirements

Before running the application, make sure you have:

- [Go](https://go.dev/dl/) installed.
- An active internet connection.
- Access to the target website.

## Installation

**1. Clone the repository**

```bash
git clone https://github.com/your-username/go-web-scraper.git
```

Navigate into the project directory:

```bash
cd go-web-scraper
```

**2. Initialize the Go module**

If you are creating the project from scratch:

```bash
go mod init go-web-scraper
```

**3. Install dependencies**

Install goquery:

```bash
go get github.com/PuerkitoBio/goquery
```

**4. Download dependencies**

```bash
go mod tidy
```

## Running the Application

Run the application using:

```bash
go run main.go
```

Alternatively, compile the application into an executable:

```bash
go build -o scraper
```

Then execute it:

```bash
./scraper
```

On Windows:

```powershell
.\scraper.exe
```

## Example Output

The application retrieves article information from Hacker News and displays it in the terminal.

```text
Starting web scraper...

1. Example Article Title
   https://example.com/article1

2. Another Interesting Article
   https://example.com/article2

3. Go Programming Resources
   https://example.com/article3

Successfully scraped 30 articles!
Results saved to output/articles.csv
```

*Note: Actual article titles, URLs, and the number of results depend on the website's current content.*

## CSV Output

After scraping, the application creates:

```text
output/articles.csv
```

The CSV file contains two columns:

| Title | Link |
|---|---|
| Example Article Title | https://example.com/article1 |
| Another Interesting Article | https://example.com/article2 |
| Go Programming Resources | https://example.com/article3 |

The exported CSV file can be opened using Microsoft Excel, Google Sheets, or other spreadsheet applications.

## How It Works

The application follows a simple scraping workflow:

### 1. Send an HTTP Request

The `net/http` package sends an HTTP GET request to the target website.

```go
response, err = client.Do(request)
```

The server returns an HTTP response containing the requested HTML document.

### 2. Parse the HTML

The goquery library converts the HTML response into a document that can be searched using CSS selectors.

```go
document, err = goquery.NewDocumentFromReader(response.Body)
```

### 3. Extract Information

The application searches for article links using a CSS selector:

```go
document.Find("span.titleline > a").Each(
    func(index int, selection *goquery.Selection) {
        // Extract article information.
    },
)
```

It extracts the article title using:

```go
title = selection.Text()
```

And retrieves the URL using:

```go
link, exists = selection.Attr("href")
```

### 4. Store the Results

Extracted information is stored using an `Article` struct:

```go
type Article struct {
    Title string
    Link  string
}
```

Multiple articles are collected in a slice:

```go
var articles []Article
```

### 5. Export to CSV

Finally, the `encoding/csv` package writes the collected information into a CSV file.

```go
writer.Write([]string{
    article.Title,
    article.Link,
})
```

## Error Handling

The application includes basic error handling to prevent unexpected failures.

Examples include:

- HTTP request failures.
- Network timeouts.
- Unexpected HTTP status codes.
- HTML parsing errors.
- File creation and writing errors.

The application also uses `defer` to close HTTP response bodies and files after they are no longer needed.

## Future Improvements

Potential features to expand this project include:

- [ ] Support multiple websites.
- [ ] Handle relative URLs and convert them into absolute URLs.
- [ ] Implement pagination to scrape multiple pages.
- [ ] Add duplicate article detection.
- [ ] Implement concurrent scraping using goroutines.
- [ ] Add retry mechanisms for failed HTTP requests.
- [ ] Store scraped information in PostgreSQL or MySQL.
- [ ] Expose scraped information through a REST API.
- [ ] Add unit tests.
- [ ] Support configurable scraping targets and CSS selectors.

## Learning Objectives

This project demonstrates several fundamental Go programming concepts:

- Structs and custom data types.
- Slices and arrays.
- Pointers and memory addresses.
- Functions and return values.
- Error handling.
- HTTP clients and requests.
- External package management.
- File handling.
- CSV data processing.
- HTML parsing.

## Responsible Web Scraping

This project is intended for educational purposes.

Before scraping a website:

- Review its Terms of Service.
- Check its `robots.txt` policies.
- Respect rate limits and server resources.
- Avoid bypassing authentication or access restrictions.
- Only collect publicly accessible information you are permitted to use.

---

## License

This project is available for educational and personal use.
