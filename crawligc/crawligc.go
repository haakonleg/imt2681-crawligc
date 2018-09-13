package crawligc

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"path"
	"regexp"
	"strings"
	"sync"
)

// Found IGC links, need a mutex to synchronize R/W
var foundIGC = make(map[string]bool)
var foundIGCLock sync.Mutex

// HTML links that have been crawled already
var crawledLinks = make(map[string]bool)
var crawledLinksLock sync.Mutex

// Regex to match links
var linkMatcher = regexp.MustCompile("<a\\s+[^>]*href=\"([^\"]+)\"")

// Channel containing urls to crawl
var urlsToCrawl = make(chan string)

// CrawlIGC starts the crawler and prints results
func CrawlIGC(baseURL string, url string) {
	var wg sync.WaitGroup
	wg.Add(1)
	go crawl(baseURL, url, &wg)
	wg.Wait()

	// Print results
	fmt.Println(len(foundIGC))
	for igc := range foundIGC {
		fmt.Println(igc)
	}
}

// Recursively visit webpages
func crawl(baseURL string, url string, wg *sync.WaitGroup) {
	defer wg.Done()

	// Check if link has been visited
	crawledLinksLock.Lock()
	_, visited := crawledLinks[url]
	if visited {
		crawledLinksLock.Unlock()
		return
	}
	crawledLinks[url] = true
	crawledLinksLock.Unlock()

	// Retrieve the page
	body, err := retrievePage(baseURL + url)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Find all IGC and HTML links on the page
	igcLinks, links := findURLs(body)

	// Add igc links to found links
	foundIGCLock.Lock()
	for _, igc := range igcLinks {
		foundIGC[igc] = true
	}
	foundIGCLock.Unlock()

	// Visit found links
	for _, link := range links {
		wg.Add(1)
		go crawl(baseURL, link, wg)
	}
}

// Finds all IGC links and other links
func findURLs(body []byte) (igcLinks, links []string) {
	urls := linkMatcher.FindAllStringSubmatch(string(body), -1)
	igcLinks = make([]string, 0, len(urls))
	links = make([]string, 0, len(urls))

	for _, url := range urls {
		foundURL := url[1]
		if isLinkIGC(foundURL) {
			igcLinks = append(igcLinks, foundURL)
		} else {
			links = append(links, foundURL)
		}
	}

	return igcLinks, links
}

// Checks if link is a link to an igc file
func isLinkIGC(link string) bool {
	ext := path.Ext(link)
	if strings.ToLower(ext) == ".igc" {
		return true
	}
	return false
}

// Retrieve the content of a webpage
func retrievePage(url string) (body []byte, e error) {
	data, err := http.Get(url)
	if err != nil {
		return nil, err
	}

	body, err = ioutil.ReadAll(data.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}
