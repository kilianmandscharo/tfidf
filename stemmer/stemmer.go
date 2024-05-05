package stemmer

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

// An implementation of: http://www.inf.fu-berlin.de/lehre/WS98/digBib/projekt/_stemming.html
// Note that as the above text says, all stop words have to be removed before
// applying the algorithm

var specialCharacters = map[string]string{
	"ä":   "a",
	"Ä":   "A",
	"ö":   "o",
	"Ö":   "O",
	"ü":   "u",
	"Ü":   "U",
	"ß":   "ss",
	"sch": "$",
	"ch":  "§",
	"ei":  "%",
	"ie":  "&",
}

var doubleCharacters = map[string]string{
	"bb": "b*",
	"dd": "d*",
	"ff": "f*",
	"gg": "g*",
	"ll": "l*",
	"mm": "m*",
	"nn": "n*",
	"pp": "p*",
	"rr": "r*",
	"ss": "s*",
	"tt": "t*",
}

func constructRegexFromMap(m map[string]string) string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	return fmt.Sprintf("(%s)", strings.Join(keys, "|"))
}

var specialCharactersRegex = regexp.MustCompile(constructRegexFromMap(specialCharacters))
var doubleCharactersRegex = regexp.MustCompile(constructRegexFromMap(doubleCharacters))

func Stem(word string) string {
	word = replaceSpecialCharacters(word)
	word = replaceDoubleCharacters(word)
	stem, ok := word, true
	for ok {
		stem, ok = stripSuffix(stem)
	}
	stem = stripSuffixFinal(stem)
	return stem
}

func replaceSpecialCharacters(word string) string {
	return specialCharactersRegex.ReplaceAllStringFunc(word, func(match string) string {
		return specialCharacters[match]
	})
}

func replaceDoubleCharacters(word string) string {
	return doubleCharactersRegex.ReplaceAllStringFunc(word, func(match string) string {
		return doubleCharacters[match]
	})
}

func stripSuffixFinal(word string) string {
	if len(word) > 3 {
		if strings.HasPrefix(word, "ge") {
			return word[2:]
		}
		if strings.HasSuffix(word, "ge") {
			return word[:len(word)-2]
		}
	}
	return word
}

func stripSuffix(word string) (string, bool) {
	l := len(word)

	if l > 5 {
		if suffix := word[l-2:]; suffix == "nd" {
			return word[:l-2], true
		}
	}

	if l > 4 {
		if suffix := word[l-2:]; suffix == "em" || suffix == "er" {
			return word[:l-2], true
		}
	}

	if l > 3 {
		suffix := word[l-1:]

		if suffix == "e" || suffix == "s" || suffix == "n" {
			return word[:l-1], true
		}

		if !unicode.IsUpper(rune(word[0])) && suffix == "t" {
			return word[:l-1], true
		}
	}

	return word, false
}
