package app

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
)

// quoteFetcher handles querying quotes from external source.
type quoteFetcher struct {
	quotes chan quote
	error  chan error
	ctx    context.Context
	stop   context.CancelFunc
}

// start starts a goroutine that fetches quotes perpetually until
// it is explicitly stopped. Fetched quotes are then enqueued inside
// the quotes channel.
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

// the quote model
type quote struct {
	Text   string `json:"content"`
	length int
}

// getRandomQuote queries a random quote from the API.
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

	quote.Text, quote.length = processText(quote.Text)
	return quote, nil
}

// processText processes the quote by substituting unicode characters
// with its equivalent ASCII representation, and removes all unicode
// characters, tabs, newlines from the string.
func processText(text string) (string, int) {
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

	// remove redundant whitespaces, tabs, newlines
	filtered = strings.Join(strings.Fields(filtered), " ")
	return filtered, len(filtered)
}

// unicodeSubstitute maps unicode character to its equivalent/similar
// ASCII character.
var unicodeSubstitute = map[rune]rune{
	'‘': '\'',
	'’': '\'',
}

// isASCII checks if the given rune belongs to the ASCII charset.
// taken from: https://github.com/scott-ainsworth/go-ascii/blob/e2eb5175fb10/ascii.go#L103
func isASCII(c rune) bool { return c <= 0x7F }
