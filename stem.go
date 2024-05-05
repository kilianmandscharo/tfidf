package tfidf

import (
	"fmt"
	"regexp"
	"strings"
	"unicode"
)

var (
	fullLower        = regexp.MustCompile(`(?:em|er|nd|e|s|n|t)\b`)
	fullUpper        = regexp.MustCompile(`(?:em|er|nd|e|s|n)\b`)
	smallerSixLower  = regexp.MustCompile(`(?:em|er|e|s|n|t)\b`)
	smallerSixUpper  = regexp.MustCompile(`(?:em|er|e|s|n)\b`)
	smallerFiveLower = regexp.MustCompile(`(?:e|s|n|t)\b`)
	smallerFiveUpper = regexp.MustCompile(`(?:e|s|n)\b`)
)

var specialCharacters = map[string]string{
	"ä": "a",
	"Ä": "A",
	"ö": "o",
	"Ö": "O",
	"ü": "u",
	"Ü": "U",
	"ß": "ss",
}

func getSpecialcharactersRegex() string {
	keys := make([]string, 0, len(specialCharacters))
	for key := range specialCharacters {
		keys = append(keys, key)
	}
	return fmt.Sprintf("[%s]", strings.Join(keys, ""))
}

var specialCharactersRegex = regexp.MustCompile(getSpecialcharactersRegex())

func stem(word string) string {
	isNoun := false
	if unicode.IsUpper(rune(word[0])) {
		isNoun = true
	}
	trimmedWord := []byte(specialCharactersRegex.ReplaceAllStringFunc(word, func(match string) string {
		return specialCharacters[match]
	}))
	for {
		match := []int{}
		switch l := len(trimmedWord); {
		case l < 4:
			break
		case l < 5:
			if isNoun {
				match = smallerFiveUpper.FindIndex(trimmedWord)
			} else {
				match = smallerFiveLower.FindIndex(trimmedWord)
			}
		case l < 6:
			if isNoun {
				match = smallerSixUpper.FindIndex(trimmedWord)
			} else {
				match = smallerSixLower.FindIndex(trimmedWord)
			}
		default:
			if isNoun {
				match = fullUpper.FindIndex(trimmedWord)
			} else {
				match = fullLower.FindIndex(trimmedWord)
			}
		}
		if match == nil {
			break
		}
		trimmedWord = trimmedWord[:match[0]]
	}
	return string(trimmedWord)
}
