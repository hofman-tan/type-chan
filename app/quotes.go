package app

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

// quoteFetcher is used to fetch quotes continuously from external source.
// Quotes are then queued inside the quotes buffer.
type quoteFetcher struct {
	quotes chan quote
	cancel context.CancelFunc
}

// start kickstarts the continuous fetch process in a goroutine.
func (q *quoteFetcher) start(buffer int) {
	q.quotes = make(chan quote, buffer)

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

// stop stops the continuous fetch by terminating the underlying goroutine.
func (q *quoteFetcher) stop() {
	if q.cancel != nil {
		q.cancel()
	}
}

// newQuoteFetcher returns a new instance of quoteFetcher.
func newQuoteFetcher() *quoteFetcher {
	return &quoteFetcher{}
}

// quote stores the quote data.
type quote struct {
	Text   string `json:"content"`
	lines  []string
	words  []string
	length int
}

// getRandomQuote reaches to the external API and retrieves a random quote.
// The quote will be returned as an instance of quote object.
func getRandomQuote() quote {
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

	var quote quote
	err = json.Unmarshal(body, &quote)
	if err != nil {
		panic(err)
	}

	quote.Text = processText(quote.Text)
	quote.lines = splitTextIntoLines(quote.Text, textareaWidth)
	quote.words = strings.Split(quote.Text, " ")
	quote.length = len(quote.Text)

	return quote
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

// splitTextIntoLines splits the given text string into slices of strings,
// while perserving any trailing spaces. The length of strings are bounded
// by the specified limit.
func splitTextIntoLines(text string, limit int) []string {
	if len(text) == 0 {
		return []string{}
	}

	// minus 1 from the limit to offset the space before adding it later on.
	wrapped := wordwrap.String(text, limit-1)
	wrapped = strings.ReplaceAll(wrapped, "\n", " \n")
	textSlices := strings.Split(wrapped, "\n")
	return textSlices
}
