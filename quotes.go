package main

import (
	"encoding/json"
	"io"
	"net/http"
)

//https://api.quotable.io/random?minLength=200
//http://www.randompassages.com/

type Quote struct {
	Content string `json:"content"`
	Author  string `json:"author"`
}

func getRandomQuote() Quote {
	url := "https://api.quotable.io/random?minLength=150"

	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		panic("Error fetching random quote from API")
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}

	var quote Quote
	err = json.Unmarshal(body, &quote)
	if err != nil {
		panic(err)
	}

	return quote
}
