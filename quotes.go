package main

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/muesli/reflow/wordwrap"
)

//https://api.quotable.io/random?minLength=200
//http://www.randompassages.com/

type QuoteFetcher struct {
	quotes chan Quote
	cancel context.CancelFunc
}

func (q *QuoteFetcher) Start(buffer int) {
	q.quotes = make(chan Quote, buffer)

	ctx, cancel := context.WithCancel(context.Background())
	q.cancel = cancel

	go func() {
		for {
			select {
			case q.quotes <- getRandomQuote():
			case <-ctx.Done():
				// stop
				close(q.quotes)
				return
			}

		}
	}()
}

func (q *QuoteFetcher) Stop() {
	if q.cancel != nil {
		q.cancel()
	}
}

func newQuoteFetcher() *QuoteFetcher {
	return &QuoteFetcher{}
}

type Quote struct {
	Text   string `json:"content"`
	lines  []string
	words  []string
	length int
}

func getRandomQuote() Quote {
	url := "https://api.quotable.io/random?minLength=100"

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

	quote.Text = processText(quote.Text)
	quote.lines = splitTextIntoLines(quote.Text)
	quote.words = strings.Split(quote.Text, " ")
	quote.length = len(quote.Text)

	return quote
}

func processText(text string) string {
	filtered := ""
	for _, rune := range text {
		// replace non-ASCII letter
		if replacement, ok := unicodeSubstitute[rune]; ok {
			rune = replacement
		}

		// remove non-ASCII letter and newline character
		if IsASCII(rune) && rune != '\n' {
			filtered += string(rune)
		}
	}
	return filtered
}

var unicodeSubstitute = map[rune]rune{
	'‘': '\'',
	'’': '\'',
}

// https://github.com/scott-ainsworth/go-ascii/blob/e2eb5175fb10/ascii.go#L103
func IsASCII(c rune) bool { return c <= 0x7F }

func splitTextIntoLines(text string) []string {
	if len(text) == 0 {
		return []string{}
	}
	wrapped := wordwrap.String(text, textareaWidth)
	wrapped = strings.ReplaceAll(wrapped, "\n", " \n")
	textSlices := strings.Split(wrapped, "\n")
	return textSlices
}
