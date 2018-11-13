package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

type Rss struct {
	Channel Channel `xml:"channel"`
}

type Channel struct {
	Title string `xml:"title"`
	Items []Item `xml:"item"`
}

type Item struct {
	Title string `xml:"title"`
	// Link    string       `xml:"link"`
	Content MediaContent `xml:"content"`
}

type MediaContent struct {
	Url string `xml:"url,attr"`
}

func downloadRss(url string) string {
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

func getNumberAndEpisode(title string) (number string, episode string) {
	strings.Split(title, ":")
	return "", ""
}

func downloadItem(item Item) {
	getNumberAndEpisode(item.Title)
}

func main() {
	fmt.Println("loading podcast files from rss")
	url := flag.String("url", "http://feeds.feedburner.com/se-radio?format=xml", "the url of the rss feed")
	flag.Parse()

	fmt.Println("url = ", strings.TrimSpace(*url))

	content := downloadRss(*url)

	rss := Rss{}
	xml.Unmarshal([]byte(content), &rss)

	//jsonrss, _ := json.MarshalIndent(rss, "", "    ")

	//fmt.Printf("%+v\n", rss)
	//fmt.Println(string(jsonrss))

	for _, item := range rss.Channel.Items {
		fmt.Println(item.Title)
	}

	//fmt.Println(content)
}
