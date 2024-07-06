package main

import (
	"fmt"
	"sort"
	"strings"
)

var maleLetters = []string{
	"a", "e", "r", "n", "l", "o", "t", "i", "h", "s", "y", "d", "j", "c", "b", "m", "u", "g", "k", "w", "p", "v", "f", "x", "z", "q",
}

var femaleLetters = []string{
	"a", "e", "i", "n", "r", "l", "h", "t", "s", "c", "y", "m", "o", "d", "b", "j", "u", "g", "k", "v", "f", "p", "q", "x", "z", "w",
}

var whiteLetters = []string{
	"e", "r", "n", "a", "o", "l", "s", "i", "t", "h", "c", "d", "m", "b", "g", "u", "y", "k", "w", "p", "f", "v", "z", "j", "x", "q",
}

var hispanicLetters = []string{
	"a", "e", "o", "r", "n", "l", "i", "s", "c", "d", "t", "m", "u", "g", "z", "b", "v", "p", "h", "y", "j", "f", "w", "q", "k", "x",
}

var blackLetters = []string{
	"e", "r", "a", "n", "o", "l", "s", "i", "t", "d", "c", "m", "h", "b", "y", "g", "u", "w", "k", "p", "f", "v", "j", "z", "x", "q",
}

var decadeLettersMale = map[string][]string{
	"2010": {"a", "n", "e", "i", "o", "r", "l", "h", "s", "c", "d", "t", "j", "m", "y", "b", "u", "v", "k", "w", "g", "p", "x", "z", "f", "q"},
	"2000": {"a", "n", "e", "r", "i", "o", "s", "t", "l", "h", "c", "j", "d", "m", "u", "b", "k", "y", "v", "g", "p", "w", "x", "z", "f", "q"},
	"1990": {"a", "e", "n", "r", "t", "i", "o", "s", "l", "h", "d", "c", "j", "y", "m", "b", "u", "k", "v", "p", "g", "w", "f", "x", "z", "q"},
	"1980": {"a", "e", "r", "n", "s", "i", "l", "o", "t", "h", "d", "y", "c", "j", "m", "u", "b", "k", "p", "g", "w", "f", "v", "x", "z", "q"},
	"1970": {"r", "a", "e", "n", "o", "l", "t", "i", "d", "y", "s", "h", "c", "j", "m", "b", "p", "g", "k", "u", "w", "f", "v", "q", "x", "z"},
	"1960": {"r", "e", "a", "n", "l", "i", "t", "y", "o", "d", "h", "c", "m", "s", "j", "b", "g", "p", "k", "w", "f", "u", "v", "q", "x", "z"},
	"1950": {"e", "r", "a", "n", "l", "i", "o", "d", "t", "h", "y", "s", "c", "m", "g", "p", "b", "f", "j", "u", "k", "w", "v", "q", "x", "z"},
	"1940": {"e", "r", "a", "l", "n", "o", "i", "d", "h", "y", "m", "t", "s", "c", "b", "j", "g", "p", "u", "w", "f", "k", "v", "q", "x", "z"},
}

var decadeLettersFemale = map[string][]string{
	"2010": {"a", "l", "e", "n", "i", "y", "r", "o", "m", "s", "h", "b", "c", "t", "d", "k", "g", "v", "u", "j", "p", "x", "z", "f", "q", "w"},
	"2000": {"a", "e", "i", "n", "l", "r", "s", "y", "h", "m", "t", "c", "b", "o", "d", "k", "g", "j", "u", "v", "x", "p", "z", "f", "q", "w"},
	"1990": {"a", "e", "i", "n", "l", "r", "s", "t", "h", "y", "c", "m", "k", "o", "b", "d", "j", "g", "u", "p", "v", "f", "x", "q", "w", "z"},
	"1980": {"a", "e", "i", "n", "r", "l", "s", "t", "c", "y", "h", "m", "d", "k", "o", "u", "b", "j", "p", "f", "g", "v", "q", "w", "z", "x"},
	"1970": {"a", "e", "n", "i", "r", "l", "t", "c", "h", "s", "y", "d", "m", "o", "b", "k", "j", "u", "p", "g", "f", "v", "w", "z", "q", "x"},
	"1960": {"a", "e", "n", "i", "r", "l", "t", "h", "y", "c", "d", "o", "s", "b", "j", "m", "k", "u", "g", "p", "v", "w", "z", "f", "q", "x"},
	"1950": {"a", "e", "n", "r", "i", "l", "t", "h", "s", "c", "o", "y", "d", "j", "b", "m", "u", "g", "k", "v", "p", "z", "f", "q", "w", "x"},
	"1940": {"a", "e", "n", "r", "l", "i", "o", "t", "s", "y", "d", "h", "j", "c", "m", "u", "b", "g", "p", "k", "v", "f", "w", "q", "z", "x"},
}

func calculateScore(name string, letters []string) int {
	score := 0
	for _, char := range strings.ToLower(name) {
		index := indexOf(string(char), letters)
		if index != -1 {
			score += 26 - index
		}
	}
	return score
}

func indexOf(str string, list []string) int {
	for i, v := range list {
		if v == str {
			return i
		}
	}
	return -1
}

func estimateGender(name string) string {
	maleScore := calculateScore(name, maleLetters)
	femaleScore := calculateScore(name, femaleLetters)
	if maleScore > femaleScore {
		return "male"
	}
	return "female"
}

func estimateRace(name string) string {
	blackScore := calculateScore(name, blackLetters)
	hispanicScore := calculateScore(name, hispanicLetters)
	whiteScore := calculateScore(name, whiteLetters)

	scores := map[string]int{
		"black":    blackScore,
		"hispanic": hispanicScore,
		"white":    whiteScore,
	}

	sortedRaces := sortMapByValue(scores)
	return sortedRaces[0]
}

func sortMapByValue(m map[string]int) []string {
	type kv struct {
		Key   string
		Value int
	}
	var ss []kv
	for k, v := range m {
		ss = append(ss, kv{k, v})
	}
	sort.Slice(ss, func(i, j int) bool {
		return ss[i].Value > ss[j].Value
	})
	var sortedKeys []string
	for _, kv := range ss {
		sortedKeys = append(sortedKeys, kv.Key)
	}
	return sortedKeys
}

func estimateAge(name, gender string) []string {
	res := make(map[string]int)
	if gender == "male" {
		for decade, letters := range decadeLettersMale {
			res[decade] = calculateScore(name, letters)
		}
	} else {
		for decade, letters := range decadeLettersFemale {
			res[decade] = calculateScore(name, letters)
		}
	}

	sortedDecades := sortMapByValue(res)
	return sortedDecades[:2]
}

func Profile(name string) map[string]string {
	parts := strings.Fields(strings.ToLower(name))
	first := parts[0]
	last := parts[len(parts)-1]

	gender := estimateGender(first)
	race := estimateRace(last)
	ageRange := estimateAge(first, gender)
	age := fmt.Sprintf("%s to %s", ageRange[0], ageRange[1])

	return map[string]string{
		"gender": gender,
		"age":    age,
		"race":   race,
	}
}
