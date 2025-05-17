package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"
)

func Manga(min_chapter string, max_chapter string, url string) {

	// scrape urls
	urls := get_urls(min_chapter, max_chapter, url)
	// scrape urls

	var wg sync.WaitGroup
	for index, url := range urls {
		wg.Add(1)

		go func(u string) {
			defer wg.Done()
			fmt.Println("Sending request for:", u)
			call_lambda(u, index)
		}(url)
	}
	wg.Wait()
}

type Lambda_response struct {
	Download_url string `json:"download_url"`
	Message      string `json:"message"`
}

func call_lambda(url string, index int) {
	API_URL := "https://6pnsxvgd6jrwb2cz2pamoxv4pq0inrjd.lambda-url.eu-north-1.on.aws/"
	data := map[string]string{
		"url": url,
	}
	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON")
		return
	}
	req, err := http.NewRequest("POST", API_URL, bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request")
		return
	}
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making the request", err)
		return
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	var response_json Lambda_response
	json.Unmarshal(body, &response_json)
	// fmt.Println(response_json.Download_url)
	filename := extractFilename(response_json.Download_url)
	err = downloadPDF(response_json.Download_url, filename)
	if err != nil {
		fmt.Println("Error downloading the file:", err)
		return
	}

	fmt.Println("PDF downloaded successfully as output.pdf")
}
func downloadPDF(url string, filename string) error {
	dir := "outputs"
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		err := os.MkdirAll(dir, os.ModePerm) // Create the folder with permissions
		if err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("failed to download file: status code %d", resp.StatusCode)
	}
	filename = dir + "/" + filename
	file, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	return err
}
func extractFilename(url string) string {
	parts := strings.Split(url, "/pdfs/")
	if len(parts) > 1 {
		filename := parts[1]                              // Get everything after "/pdfs/"
		filenameParts := strings.SplitN(filename, "?", 2) // Split at "?" (only once)
		return filenameParts[0]                           // Return part before "?"
	}
	return "" // Return empty if "/pdfs/" not found
}

func get_urls(min_chapter string, max_chapter string, url string) []string {
	// min_chapter := os.Args[1]
	// max_chapter := os.Args[2]
	// // url := "https://mangapark.io/title/12574-en-dragon-ball"
	// url := os.Args[3]
	base_mangapark_url := "https://mangapark.io"

	// Send HTTP request
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal("Error making request:", err)
	}
	defer resp.Body.Close()

	// Parse the HTML content
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		log.Fatal("Error loading HTML:", err)
	}

	var chapters []string
	in_range := false

	doc.Find("#app-wrapper main div:nth-of-type(4) div:nth-of-type(2) div div div").Each(func(index int, element *goquery.Selection) {
		// For each div, find the first child div and the anchor tag <a>
		aTag := element.Find("div:nth-of-type(1) a")
		if aTag.Length() > 0 {
			link, _ := aTag.Attr("href")
			text := aTag.Text()
			if max_chapter == text {
				in_range = true
			}
			if in_range == true {
				chapter_url := base_mangapark_url + link
				chapters = append(chapters, chapter_url)
				if min_chapter == text {
					in_range = false
				}
			}
		}
	})
	fmt.Print(chapters)
	return chapters
}
