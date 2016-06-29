// the program will take the url from command line
// the program prints all the urls' on the given link

package main

import (
	"fmt"
	// "golang.org/x/net/html"
	"io/ioutil"
	"net/http"
	"os"
)

func main() {
	// var err error
	fmt.Println(len(os.Args))
	if len(os.Args) < 2 {
		fmt.Println("usage : ./webscraper link1 link2 ...")
		return
	}

	for _, link := range os.Args[1:] {
		res, err := http.Get(link)
		if err != nil {
			fmt.Println("error fetching url : %s", err)
			return
		}

		buf, err := ioutil.ReadAll(res.Body)
		fmt.Println(string(buf))

	}
}
