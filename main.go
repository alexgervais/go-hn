package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"log"
	"flag"
	"io/ioutil"
)

type HNItem struct {
	Id          int
	Deleted     bool
	ItemType    string `json:"type"`
	By          string
	TimeRaw     int64 `json:"time"`
	Time        time.Time
	Text        string
	Dead        bool
	Parent      int
	Kids        []int
	Url         string
	Score       int
	Title       string
	Parts       []int
	Descendants int
}

func main() {

	isDebugPtr := flag.Bool("debug", false, "Display debug output")
	flag.Parse()

	if !*isDebugPtr {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	stories := getTopStories()
	item := getItem(stories[1])

	fmt.Printf("%d) %s (%s) -- %s\n", item.Id, item.Title, item.Url, time.Unix(item.TimeRaw, 0))
}

func getTopStories() []int {

	var data []int
	get("https://hacker-news.firebaseio.com/v0/topstories.json", &data)

	return data
}

func getItem(itemId int) (data *HNItem) {

	itemUrl := fmt.Sprintf("https://hacker-news.firebaseio.com/v0/item/%d.json", itemId)

	data = new(HNItem)
	get(itemUrl, data)

	return
}

func get(url string, data interface{}) {

	log.Printf("Getting %s\n", url)

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		panic(err)
	}
}

func Size(a int) string {
	switch {
	case a < 0:
		return "negative"
	case a == 0:
		return "zero"
	}

	return "small"
}
