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
}

func getRandomQuote() string {
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

	text := processText(quote.Content)
	return text
}

func processText(text string) string {
	filtered := ""
	for _, rune := range text {
		// replace non-ASCII letter
		if replacement, ok := unicodeSubstitute[rune]; ok {
			rune = replacement
		}

		// remove non-ASCII letter
		if IsASCII(rune) {
			filtered += string(rune)
		}
	}
	return filtered
}

var unicodeSubstitute = map[rune]rune{
	'‘': '\'',
	'’': '\'',
}

// Return true if c is a valid ASCII character; otherwise, return false.
// https://github.com/scott-ainsworth/go-ascii/blob/e2eb5175fb10/ascii.go#L103
func IsASCII(c rune) bool { return c <= 0x7F }
