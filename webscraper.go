// the program will take the url from command line
// use http package to make the GET request and get the response in the response.body
// use golang.org/x/net/html to parse the html response
// the program prints all the urls' on the given link

package main

import (
	"bytes"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"io/ioutil"
	"net/http"
	"os"
)

// parsing html response recieved
// parse StartTagToken
// search for <a> tag
// search for href attribute in []attr in StartTagToken
func parseHTML(body []byte) ([]string, error) {
	url := make(map[string]int)
	uniqueURLs := make([]string, 0)
	tokens := html.NewTokenizer(bytes.NewReader(body))
	for {
		token := tokens.Next()
		switch token {
		case html.ErrorToken:
			if tokens.Err() == io.EOF {
				fmt.Println("reached end of file")
				return uniqueURLs, nil
			}
			return uniqueURLs, tokens.Err()
		case html.StartTagToken:
			tag := tokens.Token()
			if tag.Data == "a" {
				for _, tt := range tag.Attr {
					if tt.Key == "href" {
						if _, ok := url[tt.Val]; !ok {
							uniqueURLs = append(uniqueURLs, tt.Val)
						}
					}
				}
			}
		}
	}
}

// func getData(links []string) {

// }

func main() {
	// var err error
	fmt.Println(len(os.Args))
	if len(os.Args) < 2 {
		fmt.Println("usage : ./webscraper link1 link2 ...")
		return
	}

	getData(os.Args[1:])

	for _, link := range os.Args[1:] {
		res, err := http.Get(link)
		if res.StatusCode != http.StatusOK {
			fmt.Println("error fetching url : %s", err)
			return
		}

		buf, err := ioutil.ReadAll(res.Body)

		urls, err := parseHTML(buf)
		if err != nil {
			fmt.Println("error parsing html ", err)
			return
		}
		// print all urls
		for _, u := range urls {
			fmt.Println(u)
		}
	}
}
