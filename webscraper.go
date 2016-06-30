// the program will take the url from command line
// use http package to make the GET request and get the response in the response.body
// use golang.org/x/net/html to parse the html response
// the program prints all the unique urls' on the given link
// uses golang concurrency to improve performance
// sends data across channels to communicate
// so the total time taken to scrap all the links is equal to the time taken by the slowest goroutine

package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

// parsing html response recieved
// parse StartTagToken
// search for <a> tag
// search for href attribute in []attr in StartTagToken, which is a "key" and the value of href is saved in "val"
func parseHTML(body []byte, linkChan chan string, sigChan chan bool) {
	url := make(map[string]int)
	// uniqueURLs := make([]string, 0)
	defer func() {
		sigChan <- true
	}()
	tokens := html.NewTokenizer(bytes.NewReader(body))
	for {
		token := tokens.Next()
		switch token {
		case html.ErrorToken:
			if tokens.Err() == io.EOF {
				fmt.Println("reached end of file")
				return
			}
		case html.StartTagToken:
			tag := tokens.Token()
			if tag.Data == "a" {
				for _, tt := range tag.Attr {
					if tt.Key == "href" {
						if _, ok := url[tt.Val]; !ok {
							if strings.HasPrefix(tt.Val, "http") {
								linkChan <- tt.Val
								// uniqueURLs = append(uniqueURLs, tt.Val)
							}
						}
					}
				}
			}
		}
	}
}

// get the data from the URL using the http.Get
// creates 2 channels 1) to recieve the links  2) to recieve the finished status
// fire individual goroutine for every inputURL
func getLinks(inURLs []string) {
	linksChan := make(chan string)
	signalChan := make(chan bool)

	for _, link := range inURLs {

		res, err := http.Get(link)
		if res.StatusCode != http.StatusOK {
			fmt.Println("error fetching url : %s", err)
			return
		}
		buf, err := ioutil.ReadAll(res.Body)
		go parseHTML(buf, linksChan, signalChan)
	}

	for count := 0; count < len(inURLs); {
		select {
		case u := <-linksChan:
			fmt.Println(u)
		case <-signalChan:
			count++
		}
	}
}

// IO : read command line inputs
func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage : ./webscraper link1 link2 ...")
		return
	}
	inputURLs := os.Args[1:]
	getLinks(inputURLs)

}
