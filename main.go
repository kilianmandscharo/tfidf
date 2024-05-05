package main

import (
	"cmp"
	"encoding/json"
	"fmt"
	"math"
	"os"
	"path"
	"regexp"
	"slices"
	"strings"
)

const stopWordsPath = "./assets/stop_words.json"
const dataPath = "./data"

func main() {
	index := New(IndexOptions{withStopWords: true, stopWords: loadStopWords()})

	files := unwrap(os.ReadDir(dataPath))

	for _, file := range files {
		fileName := file.Name()
		path := path.Join(dataPath, fileName)
		file := unwrap(os.ReadFile(path))
		document := NewDocument(fileName, string(file))
		index.AddDocument(document)
	}

	fmt.Println("zimt", index.Search("zimt"))
}

type document struct {
	title, text string
}

func NewDocument(title, text string) document {
	return document{title: title, text: text}
}

type index struct {
	index         map[string]map[string]float64
	withStopWords bool
	stopWords     map[string]struct{}
	nDocuments    int
}

type IndexOptions struct {
	withStopWords bool
	stopWords     []string
}

func New(options IndexOptions) index {
	stopWords := make(map[string]struct{})
	for _, word := range options.stopWords {
		stopWords[word] = struct{}{}
	}

	return index{
		index:         make(map[string]map[string]float64),
		withStopWords: options.withStopWords,
		stopWords:     stopWords,
	}
}

type searchResultEntry struct {
	Name  string
	Tfidf float64
}

func (idx *index) Search(word string) []searchResultEntry {
	result := []searchResultEntry{}

	if entry, ok := idx.index[word]; ok {
		for doc, tfidf := range entry {
			result = append(result, searchResultEntry{Name: doc, Tfidf: tfidf})
		}
	}

	slices.SortStableFunc(result, func(a, b searchResultEntry) int {
		return cmp.Compare(b.Tfidf, a.Tfidf)
	})

	return result[:5]
}

func (idx *index) SaveToDisk(path string) error {
	json, err := json.Marshal(idx.index)
	if err != nil {
		return err
	}

	err = os.WriteFile(path, json, 0664)
	if err != nil {
		return err
	}

	return nil
}

func (idx *index) AddDocument(doc document) {
	idx.nDocuments++

	wordFrequencies := idx.processDocumentText(doc.text)
	for word, tf := range wordFrequencies {
		if _, ok := idx.index[word]; !ok {
			idx.index[word] = make(map[string]float64)
		}
		idf := math.Log10(float64(idx.nDocuments) / float64(len(idx.index[word])+1))
		idx.index[word][doc.title] = tf * idf
	}
}

func (idx *index) processDocumentText(text string) map[string]float64 {
	totalWords, wordCounts := idx.getWordCountsFromText(text)

	wordFrequencies := make(map[string]float64)
	for word, count := range wordCounts {
		wordFrequencies[word] = float64(count) / float64(totalWords)
	}

	return wordFrequencies
}

func (idx *index) getWordCountsFromText(text string) (int, map[string]int) {
	wordCounts := make(map[string]int)
	totalWords := 0

	for _, line := range strings.Split(text, "\n") {
		for _, word := range strings.Split(line, " ") {
			if cleanedWord := idx.cleanWord(word); !idx.shouldDiscardWord(cleanedWord) {
				wordCounts[cleanedWord]++
				totalWords++
			}
		}
	}

	return totalWords, wordCounts
}

func (idx *index) cleanWord(word string) string {
	word = strings.TrimSpace(word)
	word = strings.ToLower(word)
	word = string(
		regexp.MustCompile(`[^a-zA-ZÜüÖöÄäß]`).ReplaceAll([]byte(word), []byte("")),
	)
	return word
}

func (idx *index) shouldDiscardWord(word string) bool {
	if len(word) == 0 {
		return true
	}
	if idx.withStopWords {
		if _, ok := idx.stopWords[word]; ok {
			return true
		}
	}
	return false
}

func loadStopWords() []string {
	stopWordsJson := unwrap(os.ReadFile(stopWordsPath))
	stopWords := []string{}
	unwrapErr(json.Unmarshal(stopWordsJson, &stopWords))
	return stopWords
}

func unwrap[T any](val T, err error) T {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
	return val
}

func unwrapErr(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "%v", err)
		os.Exit(1)
	}
}
