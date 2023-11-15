package app

import (
	"fmt"
	"strings"
)

type wordHolder struct {
	word string
}

func (w *wordHolder) View() string {
	return fmt.Sprintf("%s> %s", strings.Repeat(" ", paddingX), w.word)
}

func (w *wordHolder) add(letter string) {
	w.word += letter
}

func (w *wordHolder) pop() string {
	word := []rune(w.word)

	if len(word) > 0 {
		lastLetter := word[len(word)-1]
		word = word[:len(word)-1] // remove the last letter
		w.word = string(word)
		return string(lastLetter)
	}
	return ""
}

func (w *wordHolder) clear() {
	w.word = ""
}

func newWordHolder() *wordHolder {
	return &wordHolder{}
}
