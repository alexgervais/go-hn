package main

import (
	"fmt"
	"net/http"
	"encoding/json"
	"time"
	"log"
	"io/ioutil"
	"github.com/alexgervais/go-hn/lib"
)

type Config struct {
	BaseURL string `yaml:"base-url"`
}

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

type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

func main() {

	config := lib.LoadConfig()
	if !config.EnableDebug {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}

	client := &Client{config.BaseURL, http.DefaultClient}

	stories := GetTopStories(client)
	item := client.GetItem(stories[0])

	fmt.Printf("%d) %s (%s) -- %s\n", item.Id, item.Title, item.Url, time.Unix(item.TimeRaw, 0))
}

func GetTopStories(client *Client) []int {

	var data []int
	client.get("/topstories.json", &data)

	return data
}

func (client *Client) GetItem(itemId int) (data *HNItem) {

	itemUrl := fmt.Sprintf("/item/%d.json", itemId)

	data = new(HNItem)
	client.get(itemUrl, data)

	return
}

func (client *Client) get(uri string, data interface{}) {

	url := client.BaseURL + uri
	log.Printf("GET %s\n", url)

	resp, err := client.HTTPClient.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(data); err != nil {
		panic(err)
	}
}
