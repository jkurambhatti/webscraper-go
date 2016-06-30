// webscraper using recursion

package main

import (
	"fmt"
	"golang.org/x/net/html"
	// "io/ioutil"
	"net/http"
	"os"
)

func main() {
	// seedurls := os.Args[1:]
	v := os.Args[1]
	fmt.Println(v)
	res, err := http.Get(v)
	if res.StatusCode != http.StatusOK {
		fmt.Println(err)
	}

	nodes, _ := html.Parse(res.Body)

	for _, link := range do(nil, nodes) {
		fmt.Println(link)
	}
}

func do(links []string, n *html.Node) []string {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, at := range n.Attr {
			if at.Key == "href" {
				links = append(links, at.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		links = do(links, c)
	}
	return links
}
