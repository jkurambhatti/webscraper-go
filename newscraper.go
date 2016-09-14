// webscraper using recursion

package main

import (
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"os"
	"strings"
)

var link string
var name int

func main() {
	link = os.Args[1]
	fmt.Println(link)
	parser(link)
	for {
	}
}

func parser(link string) {

	res, err := http.Get(link)
	if res.StatusCode != http.StatusOK {
		fmt.Println(err)
	}
	defer res.Body.Close()

	nodes, _ := html.Parse(res.Body)
	// getlinks(nodes)
	getimageslink(nodes)
}

func getlinks(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, at := range n.Attr {
			if at.Key == "href" {
				if strings.HasSuffix(at.Val, ".php") {
					fmt.Println(at.Val)
					// go getimageslink(os.Args[1] + at.Val)
				}
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getlinks(c)
	}

}

func getimageslink(n *html.Node) {
	if n.Type == html.ElementNode && n.Data == "img" {
		for _, at := range n.Attr {
			if at.Key == "src" {
				go downloadimages(at.Val)
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		getimageslink(c)
	}
}

func downloadimages(imglink string) {
	fmt.Println(imglink)
	r, err := http.Get(imglink)

	if err != nil {
		fmt.Println(err)
		return
	}
	name++
	f, err := os.Create(fmt.Sprint("image", name, ".jpg"))
	if err != nil {
		fmt.Println(err)
		return
	}

	io.Copy(f, r.Body)
	f.Close()
	r.Body.Close()
}
