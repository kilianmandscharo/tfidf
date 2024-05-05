package tfidf

import (
	"fmt"
	"testing"
)

func TestStem(t *testing.T) {
	testCases := []struct{ word, stem string }{
		{word: "singt", stem: "sing"},
		{word: "singen", stem: "sing"},
		{word: "beliebt", stem: "belieb"},
		{word: "beliebtester", stem: "belieb"},
		{word: "stören", stem: "stö"},
		{word: "stöhnen", stem: "stöh"},
	}

	for _, test := range testCases {
		stem := stem(test.word)
		fmt.Println(stem)
	}
}
