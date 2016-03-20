package main

import (
	"testing"
	"net/http/httptest"
	"net/http"
	"fmt"
	"net/url"
	"github.com/stretchr/testify/assert"
	"errors"
)

type expectedRequest struct {
	method string
	uri    string
}

type fakeResponse struct {
	code        int
	body        string
	contentType string
	err         error
}

func ok(body string) *fakeResponse {
	return &fakeResponse{
		code: http.StatusOK,
		contentType: "application/json",
		body: body,
	}
}

func notFound(body string) *fakeResponse {
	return &fakeResponse{
		code: http.StatusNotFound,
		contentType: "plain/text",
		body: body,
	}
}

func apiError(text string) *fakeResponse {
	return &fakeResponse{
		err: errors.New(text),
	}
}

func testJsonClient(t *testing.T, expected *expectedRequest, response *fakeResponse) (*httptest.Server, *Client) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(response.code)
		w.Header().Set("Content-Type", response.contentType)
		fmt.Fprintln(w, response.body)
	}))

	transport := &http.Transport{
		Proxy: func(req *http.Request) (*url.URL, error) {
			if (response.err != nil) {
				return nil, response.err
			}

			assert.Equal(t, expected.method, req.Method)
			assert.Equal(t, expected.uri, req.URL.RequestURI())
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
	println(err.Error())
	assert.Contains(t, err.Error(), panicErrorContains)
}

func Test_GetTopStories_Success(t *testing.T) {
	server, client := testJsonClient(t, &expectedRequest{http.MethodGet, "/topstories.json"}, ok("[1, 2]"))
	defer server.Close()

	result := GetTopStories(client)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, []int{1, 2 }, result)
}

func Test_GetTopStories_UnmarshalFailure(t *testing.T) {

	server, client := testJsonClient(t, &expectedRequest{"GET", "/topstories.json"}, notFound("{}"))
	defer server.Close()
	defer expectPanic(t, "json: cannot unmarshal object")

	result := GetTopStories(client)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, []int{1, 2 }, result)
}

func Test_GetTopStories_ApiFailure(t *testing.T) {

	server, client := testJsonClient(t, &expectedRequest{"GET", "/topstories.json"}, apiError("failed!"))
	defer server.Close()
	defer expectPanic(t, "failed!")

	result := GetTopStories(client)

	assert.Equal(t, 2, len(result))
	assert.Equal(t, []int{1, 2 }, result)
}

func Test_GetItem_Success(t *testing.T) {
	server, client := testJsonClient(t, &expectedRequest{"GET", "/item/666.json"}, ok(`{"id": 666, "title": "Good news!"}`))
	defer server.Close()

	result := client.GetItem(666)

	assert.Equal(t, &HNItem{
		Id: 666,
		Title: "Good news!",
	}, result)
}
