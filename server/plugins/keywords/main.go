package main

import (
	"regexp"
	"sort"
	"strings"
)

var stopWords = map[string]struct{}{
	"i": {}, "me": {}, "my": {}, "myself": {}, "we": {}, "our": {}, "ours": {}, "ourselves": {}, "you": {}, "your": {}, "yours": {}, "yourself": {}, "yourselves": {}, "he": {}, "him": {}, "his": {}, "himself": {}, "she": {}, "her": {}, "hers": {}, "herself": {}, "it": {}, "its": {}, "itself": {}, "they": {}, "them": {}, "their": {}, "theirs": {}, "themselves": {}, "what": {}, "which": {}, "who": {}, "whom": {}, "this": {}, "that": {}, "these": {}, "those": {}, "am": {}, "is": {}, "are": {}, "was": {}, "were": {}, "be": {}, "been": {}, "being": {}, "have": {}, "has": {}, "had": {}, "having": {}, "do": {}, "does": {}, "did": {}, "doing": {}, "a": {}, "an": {}, "the": {}, "and": {}, "but": {}, "if": {}, "or": {}, "because": {}, "as": {}, "until": {}, "while": {}, "of": {}, "at": {}, "by": {}, "for": {}, "with": {}, "about": {}, "against": {}, "between": {}, "into": {}, "through": {}, "during": {}, "before": {}, "after": {}, "above": {}, "below": {}, "to": {}, "from": {}, "up": {}, "down": {}, "in": {}, "out": {}, "on": {}, "off": {}, "over": {}, "under": {}, "again": {}, "further": {}, "then": {}, "once": {}, "here": {}, "there": {}, "when": {}, "where": {}, "why": {}, "how": {}, "all": {}, "any": {}, "both": {}, "each": {}, "few": {}, "more": {}, "most": {}, "other": {}, "some": {}, "such": {}, "no": {}, "nor": {}, "not": {}, "only": {}, "own": {}, "same": {}, "so": {}, "than": {}, "too": {}, "very": {}, "s": {}, "t": {}, "can": {}, "will": {}, "just": {}, "don": {}, "should": {}, "now": {},
}

func ExtractKeywords(text string, numKeywords int) []string {
	// Clean the text by removing punctuation and converting to lower case
	re := regexp.MustCompile(`[^\w\s]`)
	cleanText := re.ReplaceAllString(strings.ToLower(text), "")

	// Split the text into an array of words
	words := strings.Fields(cleanText)

	// Filter out the stop words
	filteredWords := make([]string, 0)
	for _, word := range words {
		if _, found := stopWords[word]; !found {
			filteredWords = append(filteredWords, word)
		}
	}

	// Count the frequency of each word
	wordFrequency := make(map[string]int)
	for _, word := range filteredWords {
		wordFrequency[word]++
	}

	// Convert the word frequency map to a slice of pairs
	type kv struct {
		Key   string
		Value int
	}
	var wordFrequencyArray []kv
	for k, v := range wordFrequency {
		wordFrequencyArray = append(wordFrequencyArray, kv{k, v})
	}

	// Sort the slice by frequency in descending order
	sort.Slice(wordFrequencyArray, func(i, j int) bool {
		return wordFrequencyArray[i].Value > wordFrequencyArray[j].Value
	})

	// Extract the top N keywords
	keywords := make([]string, 0, numKeywords)
	for i, kv := range wordFrequencyArray {
		if i >= numKeywords {
			break
		}
		keywords = append(keywords, kv.Key)
	}

	return keywords
}

// func main() {
// 	text := "This is an example text. It is meant to demonstrate the extraction of keywords from text."
// 	keywords := extractKeywords(text, 3)
// 	fmt.Println(keywords) // Output: [example text meant demonstrate extraction]
// }
