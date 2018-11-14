package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"path"
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
	parts := strings.Split(title, ":")
	if len(parts) == 2 {

	}
	return "", ""
}

func downloadAndSaveItem(episodeNumber int, item Item) {
	filename := createFileName(item, episodeNumber)

	resp, err := http.Get(item.Content.Url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	file, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	written, err := io.Copy(file, resp.Body)
	if err != nil {
		panic(err)
	}

	fmt.Printf("%d bytes written\n", written)
}

func createFileName(item Item, episodeNumber int) string {
	var u url.URL
	u.Path = item.Content.Url
	filename := path.Base(u.Path)
	filename = fmt.Sprintf("SE-%07d_%s", episodeNumber, filename)
	return filename
}

type WorkItem struct {
	I Item
	N int
}

var queue = make(chan WorkItem)
var done = make(chan bool)

func produce(q chan WorkItem, url string) {
	content := downloadRss(url)

	rss := Rss{}
	xml.Unmarshal([]byte(content), &rss)

	for i := len(rss.Channel.Items) - 1; i >= 0; i-- {
		item := rss.Channel.Items[i]
		if item.Content.Url != "" {
			q <- WorkItem{I: item, N: i}
		}
	}

	done <- true
}

func consume(q chan WorkItem) {
	for {
		item := <-q
		downloadAndSaveItem(item.N, item.I)
	}
}

func main() {
	fmt.Println("loading podcast files from rss")
	url := flag.String("url", "http://feeds.feedburner.com/se-radio?format=xml", "the url of the rss feed")
	flag.Parse()

	fmt.Println("url = ", strings.TrimSpace(*url))

	go produce(queue, *url)

	go consume(queue)
	go consume(queue)
	go consume(queue)
	go consume(queue)

	<-done

	//for _, item := range rss.Channel.Items {
	//	fmt.Println(item.Title)
	//}

	//fmt.Println(content)
}
