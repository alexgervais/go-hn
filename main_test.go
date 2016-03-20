package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"fmt"
	"net/url"
	"github.com/stretchr/testify/assert"
)

func testJsonClient(t *testing.T, expectedMethod string, expectedUri string, code int, body string) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(code)
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			assert.Equal(t, expectedMethod, req.Method)
			assert.Equal(t, expectedUri, req.URL.RequestURI())
			return url.Parse(server.URL)
		},
	}

	httpClient := &http.Client{Transport: transport}
	client := &Client{server.URL, httpClient}

	return server, client
}

func expectPanic(t *testing.T, panicErrorContains string) {
	r := recover()
	if r == nil {
		t.Errorf("The code did not panic")
		return
	}
	err := r.(error)
	assert.Contains(t, err.Error(), panicErrorContains)
}

func Test_GetTopStories_Success(t *testing.T) {
	server, client := testJsonClient(t, "GET", "/topstories.json", 200, `[1, 2]`)
	defer server.Close()

	result := GetTopStories(client)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, []int{1, 2 }, result)
}

func Test_GetTopStories_UnmarshalFailure(t *testing.T) {

	server, client := testJsonClient(t, "GET", "/topstories.json", 200, "{}")
	defer server.Close()
	defer expectPanic(t, "json: cannot unmarshal object")

	result := GetTopStories(client)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, []int{1, 2 }, result)
}

func Test_GetItem_Success(t *testing.T) {
	server, client := testJsonClient(t, "GET", "/item/666.json", 200, `{"id": 666, "title": "Good news!"}`)
	defer server.Close()

	result := client.GetItem(666)

	assert.Equal(t, &HNItem{
		Id: 666,
		Title: "Good news!",
	}, result)
}
