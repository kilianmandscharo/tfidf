package stemmer

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
		{word: "geliebt", stem: "lieb"},
	}

	for _, test := range testCases {
		stem := Stem(test.word)
		fmt.Println(stem)
	}
}
