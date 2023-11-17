package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

//https://api.quotable.io/random?minLength=200
//http://www.randompassages.com/

// quoteFetcher is used to fetch quotes continuously from external source.
// Quotes are then queued inside the quotes buffer.
type quoteFetcher struct {
	quotes chan quote
	error  chan error
	ctx    context.Context
	stop   context.CancelFunc
}

// start kickstarts the continuous fetch process in a goroutine.
func (q *quoteFetcher) start(buffer int) {
	go func() {
		for {
			quote, err := getRandomQuote()
			select {
			case q.quotes <- quote:
			case q.error <- err:
			case <-q.ctx.Done():
				// stop
				close(q.quotes)
				return
			}

		}
	}()
}

// // stop stops the continuous fetch by terminating the underlying goroutine.
// func (q *quoteFetcher) stop() {
// 	if q.cancel != nil {
// 		q.cancel()
// 	}
// }

// newQuoteFetcher returns a new instance of quoteFetcher.
func newQuoteFetcher(ctx context.Context) *quoteFetcher {
	cancelCtx, cancel := context.WithCancel(ctx)

	return &quoteFetcher{
		quotes: make(chan quote, quoteBufferSize),
		error:  make(chan error, 1),
		ctx:    cancelCtx,
		stop:   cancel,
	}
}

// quote stores the quote data.
type quote struct {
	Text   string `json:"content"`
	length int
}

// getRandomQuote reaches to the external API and retrieves a random quote.
// The quote will be returned as an instance of quote object.
func getRandomQuote() (quote, error) {
	url := "https://api.quotable.io/random?minLength=100"
	var quote quote

	resp, err := http.Get(url)
	if err != nil {
		return quote, err
	}

	if resp.StatusCode != 200 {
		return quote, fmt.Errorf("API returns code %v: %s", resp.StatusCode, resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return quote, err
	}
	defer resp.Body.Close()

	err = json.Unmarshal(body, &quote)
	if err != nil {
		return quote, err
	}

	quote.Text = processText(quote.Text)
	quote.length = len(quote.Text)

	return quote, nil
}

// processText sanitizes the given text string by substituting any unicode
// characters with its equivalent ASCII representation, and removes any
// non-ASCII and newline characters from the string.
func processText(text string) string {
	filtered := ""
	for _, rune := range text {
		// replace non-ASCII letter
		if replacement, ok := unicodeSubstitute[rune]; ok {
			rune = replacement
		}

		// remove non-ASCII letter and newline character
		if isASCII(rune) && rune != '\n' {
			filtered += string(rune)
		}
	}
	return filtered
}

// unicodeSubstitute is a table that maps unicode character to its
// equivalent/similar ASCII character.
var unicodeSubstitute = map[rune]rune{
	'‘': '\'',
	'’': '\'',
}

// isASCII checks if the given rune falls under the ASCII charset.
// taken from: https://github.com/scott-ainsworth/go-ascii/blob/e2eb5175fb10/ascii.go#L103
func isASCII(c rune) bool { return c <= 0x7F }
