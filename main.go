package main

import (
	"fmt"
	"flag"
	"strings"
	"net/http"
	"io/ioutil"
)

func downloadFile(url string) string {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	return fmt.Sprintf("%s\n", body)
}

func main() {
	fmt.Println("loading podcast files from rss")
	url := flag.String("url", "http://feeds.feedburner.com/se-radio?format=xml", "the url of the rss feed")
	flag.Parse()

	fmt.Println("url = ", strings.TrimSpace(*url))

	content := downloadFile(*url)

	fmt.Println(content)
}
