
package main

import (
	"encoding/csv"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/PuerkitoBio/goquery"
)
// Website we want to scrape.
var websiteURL string = "https://news.ycombinator.com/"

// Article represents an article collected from a website.
//
// A struct allows us to group related information together.
// Instead of storing titles and links in separate arrays,
// we can keep each article's information in one object.
type Article struct {
	Title string
	Link  string
}

// ScrapeWebsite sends an HTTP request to a website,
// downloads its HTML, and extracts article information.
//
// Parameters:
//   url: The website address we want to scrape.
//
// Returns:
//   []Article: A slice containing the scraped articles.
//   error: An error if something goes wrong.
func ScrapeWebsite(url string) ([]Article, error) {

	// Create an HTTP client.
	//
	// The client is responsible for sending HTTP requests
	// and receiving responses from websites.
	//
	// Timeout prevents our program from waiting indefinitely
	// if the website takes too long to respond.
	var client http.Client
	client.Timeout = 15 * time.Second

	// Create an HTTP GET request.
	//
	// GET is an HTTP method used to retrieve information
	// from a server without modifying its data.
	//
	// The third argument is the request body.
	// We pass nil because we are not sending any data.
	var request *http.Request
	var err error

	request, err = http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	// Set the User-Agent HTTP header.
	//
	// Websites can use this header to identify the software
	// making the request.
	//
	// Providing a descriptive User-Agent is good practice
	// when developing web scraping applications.
	request.Header.Set("User-Agent", "GoWebScraper/1.0")

	// Send the HTTP request to the website.
	//
	// The response contains information from the server,
	// including its status code, headers, and HTML body.
	var response *http.Response

	response, err = client.Do(request)
	if err != nil {
		return nil, err
	}

	// Close the response body when this function finishes.
	//
	// HTTP response bodies should always be closed to release
	// resources associated with the network connection.
	defer response.Body.Close()

	// Check the HTTP status code.
	//
	// HTTP 200 means the request was successful.
	// Other status codes may indicate errors such as:
	// 404: Page not found.
	// 403: Access forbidden.
	// 500: Internal server error.
	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"unexpected HTTP status: %s",
			response.Status,
		)
	}

	// Parse the HTML response using goquery.
	//
	// response.Body contains the raw HTML received from
	// the website.
	//
	// goquery converts that HTML into a document that we
	// can search using CSS selectors.
	var document *goquery.Document

	document, err = goquery.NewDocumentFromReader(response.Body)
	if err != nil {
		return nil, err
	}

	// Create a slice to store the articles we discover.
	//
	// A slice is a dynamically sized collection.
	// Every time we discover another article, we can append
	// it to this slice.
	var articles []Article

	// Find all HTML elements matching this CSS selector:
	//
	// span.titleline > a
	//
	// This selector finds anchor elements (<a>) that are
	// direct children of spans with the class "titleline".
	//
	// Hacker News uses this structure for its article titles.
	//
	// Each() executes the provided function once for every
	// matching HTML element.
	document.Find("span.titleline > a").Each(
		func(index int, selection *goquery.Selection) {

			// Extract the visible text from the HTML element.
			// This will be the article title.
			var title string

			title = selection.Text()

			// Extract the href attribute.
			// The href attribute contains the article's URL.
			//
			// At this point, the link may be relative rather
			// than a complete URL.
			var link string
			var exists bool

			link, exists = selection.Attr("href")

			// Skip this element if it does not contain a link.
			if !exists {
				return
			}

			// Create an Article struct and append it
			// to our collection of scraped articles.
			var article Article

			article.Title = title
			article.Link = link

			articles = append(articles, article)
		},
	)

	// Return all collected articles.
	//
	// The second return value is nil because no error occurred.
	return articles, nil
}

// SaveToCSV writes the scraped articles into a CSV file.
//
// Parameters:
//   articles: The collection of articles to save.
//   filename: The name and location of the output file.
//
// Returns:
//   error: An error if creating or writing the file fails.
func SaveToCSV(articles []Article, filename string) error {

	// Create a new file.
	//
	// If the file already exists, os.Create overwrites it.
	// If the directory does not exist, this operation fails.
	var file *os.File
	var err error

	file, err = os.Create(filename)
	if err != nil {
		return err
	}

	// Close the file when this function finishes.
	defer file.Close()

	// Create a CSV writer.
	//
	// The CSV writer converts our data into comma-separated
	// values and writes them to the file.
	var writer *csv.Writer

	writer = csv.NewWriter(file)

	// Flush any buffered data before the function returns.
	//
	// CSV writers buffer data to reduce the number of
	// individual write operations.
	//
	// Flush ensures the buffered information is written.
	defer writer.Flush()

	// Write the column names as the first row.
	//
	// This makes the CSV file easier to understand
	// when opened in Excel or another spreadsheet program.
	var header []string

	header = []string{"Title", "Link"}

	err = writer.Write(header)
	if err != nil {
		return err
	}

	// Loop through every article in the slice.
	for _, article := range articles {

		// Create a row containing the article's information.
		var row []string

		row = []string{
			article.Title,
			article.Link,
		}

		// Write the row into the CSV file.
		err = writer.Write(row)
		if err != nil {
			return err
		}
	}

	// Check whether the CSV writer encountered an error
	// while writing buffered data.
	return writer.Error()
}

// main is the entry point of this application.
//
// Execution of the program begins here.
func main() {

	// Website we want to scrape: websiteURL

	fmt.Println("Starting web scraper...")

	// Call ScrapeWebsite to retrieve and extract articles.
	//
	// The function returns two values:
	// 1. A slice containing the scraped articles.
	// 2. An error value.
	//
	// We check the error before continuing.
	var articles []Article
	var err error

	articles, err = ScrapeWebsite(websiteURL)
	if err != nil {
		log.Fatal("Scraping failed:", err)
	}

	// Display the collected articles in the terminal.
	//
	// The range keyword allows us to iterate through
	// every element in a slice.
	//
	// index represents the current position.
	// article represents the current Article struct.
	for index, article := range articles {

		fmt.Printf("%d. %s\n", index+1, article.Title)
		fmt.Printf("   %s\n\n", article.Link)
	}

	// Create the output directory if it does not exist.
	//
	// MkdirAll creates all necessary parent directories.
	// The permission value 0755 gives the owner full access
	// and allows other users to read and execute.
	err = os.MkdirAll("output", 0755)
	if err != nil {
		log.Fatal("Could not create output directory:", err)
	}

	// Save the collected articles into a CSV file.
	var outputFilename string

	outputFilename = "output/articles.csv"

	err = SaveToCSV(articles, outputFilename)
	if err != nil {
		log.Fatal("CSV export failed:", err)
	}

	// Display a completion message.
	fmt.Printf(
		"Successfully scraped %d articles!\n",
		len(articles),
	)

	fmt.Println("Results saved to", outputFilename)
}
