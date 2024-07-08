package fixtext

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"
	"unicode"
)

var months = []string{
	"january", "february", "march", "april", "may", "june", "july", "august", "september", "october", "november", "december",
}

var monthCombos = append(months, []string{
	"jan", "feb", "mar", "apr", "jun", "jul", "aug", "sep", "oct", "nov", "dec",
}...)

var ordinals = []string{
	"zeroth", "first", "second", "third", "fourth", "fifth", "sixth", "seventh", "eighth", "ninth", "tenth",
	"eleventh", "twelfth", "thirteenth", "fourteenth", "fifteenth", "sixteenth", "seventeenth", "eighteenth",
	"nineteenth", "twentieth", "twenty-first", "twenty-second", "twenty-third", "twenty-fourth", "twenty-fifth",
	"twenty-sixth", "twenty-seventh", "twenty-eighth", "twenty-ninth", "thirtieth", "thirty-first",
}

var numberWords = []string{
	"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine", "ten", "eleven", "twelve",
	"thirteen", "fourteen", "fifteen", "sixteen", "seventeen", "eighteen", "nineteen", "twenty", "twenty-one",
	"twenty-two", "twenty-three", "twenty-four", "twenty-five", "twenty-six", "twenty-seven", "twenty-eight",
	"twenty-nine", "thirty", "thirty-one", "thirty-two", "thirty-three", "thirty-four", "thirty-five",
	"thirty-six", "thirty-seven", "thirty-eight", "thirty-nine", "forty", "forty-one", "forty-two", "forty-three",
	"forty-four", "forty-five", "forty-six", "forty-seven", "forty-eight", "forty-nine", "fifty", "fifty-one",
	"fifty-two", "fifty-three", "fifty-four", "fifty-five", "fifty-six", "fifty-seven", "fifty-eight",
	"fifty-nine", "sixty", "sixty-one", "sixty-two", "sixty-three", "sixty-four", "sixty-five", "sixty-six",
	"sixty-seven", "sixty-eight", "sixty-nine", "seventy", "seventy-one", "seventy-two", "seventy-three",
	"seventy-four", "seventy-five", "seventy-six", "seventy-seven", "seventy-eight", "seventy-nine", "eighty",
	"eighty-one", "eighty-two", "eighty-three", "eighty-four", "eighty-five", "eighty-six", "eighty-seven",
	"eighty-eight", "eighty-nine", "ninety", "ninety-one", "ninety-two", "ninety-three", "ninety-four",
	"ninety-five", "ninety-six", "ninety-seven", "ninety-eight", "ninety-nine",
}

var yearNumberWords = []string{
	"one hundred", "two hundred", "three hundred", "four hundred", "five hundred", "six hundred", "seven hundred",
	"eight hundred", "nine hundred", "one thousand", "eleven", "twelve", "thirteen", "fourteen", "fifteen",
	"sixteen", "seventeen", "eighteen", "nineteen", "two thousand",
}

func numberToWords(number int) string {
	if number <= 99 {
		return numberWords[number]
	}
	words := ""
	for _, digit := range strconv.Itoa(number) {
		digitInt, _ := strconv.Atoi(string(digit))
		words += numberWords[digitInt] + " "
	}
	return strings.TrimSpace(words)
}

func dayToOrdinal(day int) string {
	if day <= 31 {
		return ordinals[day]
	}
	return numberToWords(day)
}

func yearToWords(year int) string {
	yearStr := strconv.Itoa(year)
	if len(yearStr) == 4 {
		firstPart, _ := strconv.Atoi(yearStr[:2])
		secondPart, _ := strconv.Atoi(yearStr[2:])
		return yearNumberWords[firstPart-1] + " " + numberWords[secondPart]
	}
	words := ""
	for _, digit := range yearStr {
		digitInt, _ := strconv.Atoi(string(digit))
		words += numberWords[digitInt] + " "
	}
	return strings.TrimSpace(words)
}

func ReplaceDates(text string) string {
	textSplit := strings.Fields(text)
	numbers := []struct {
		text  string
		index int
	}{}

	for i, word := range textSplit {
		if _, err := strconv.Atoi(word); err == nil {
			numbers = append(numbers, struct {
				text  string
				index int
			}{word, i})
		}
	}

	replaceable := []struct {
		text  string
		index int
		line  []string
	}{}

	for i, target := range numbers {
		neighbors := textSplit[max(target.index-2, 0):min(target.index+2, len(textSplit))]

		if i > 0 && target.index-numbers[i-1].index <= 1 {
			continue
		}

		combos := generateCombinations(neighbors)
		sort.Slice(combos, func(a, b int) bool {
			return len(combos[a]) > len(combos[b])
		})

		var foundDate time.Time
		var bestLine string
		filter := time.Now().String()
		for _, combo := range combos {
			line := strings.Join(combo, " ")
			date, err := time.Parse("January 2 2006", keepOnlyMonths(line))
			if err == nil {
				if !foundDate.IsZero() {
					if strings.Contains(line, target.text) {
						if len(strings.Fields(keepOnlyMonths(line))) >= len(strings.Fields(keepOnlyMonths(bestLine))) {
							bestLine = line
						}
					}
				} else if strings.Contains(line, target.text) {
					foundDate = date
					bestLine = line
				}
			}
		}

		if bestLine != filter {
			parts := strings.FieldsFunc(bestLine, func(r rune) bool {
				return !unicode.IsLetter(r) && !unicode.IsNumber(r)
			})
			str := ""
			year := yearToWords(foundDate.Year())
			month := months[foundDate.Month()-1]
			day := dayToOrdinal(foundDate.Day())

			switch len(parts) {
			case 3:
				str = fmt.Sprintf("%s %s %s", month, day, year)
			case 2:
				if containsIgnoreCase(monthCombos, parts[0]) {
					str = month
				}
				if yearInt, err := strconv.Atoi(parts[1]); err == nil && len(parts[1]) == 4 {
					str += fmt.Sprintf(" %s", yearToWords(yearInt))
				} else {
					str += fmt.Sprintf(" %s", day)
				}
			case 1:
				if len(parts[0]) == 4 {
					str = year
				}
			}

			replaceable = append(replaceable, struct {
				text  string
				index int
				line  []string
			}{str, target.index, strings.Fields(bestLine)})
		}
	}

	return strings.Join(replaceLinesWithText(replaceable, textSplit), " ")
}

func generateCombinations(array []string) [][]string {
	result := [][]string{}
	var helper func(start int, combination []string)
	helper = func(start int, combination []string) {
		result = append(result, append([]string(nil), combination...))
		for i := start; i < len(array); i++ {
			combination = append(combination, array[i])
			helper(i+1, combination)
			combination = combination[:len(combination)-1]
		}
	}
	helper(0, []string{})
	return result[1:]
}

func replaceLinesWithText(objects []struct {
	text  string
	index int
	line  []string
}, textArray []string) []string {
	for _, obj := range objects {
		startIndex := findLineIndex(obj.line, textArray[max(0, obj.index-5):min(len(textArray), obj.index+5)])
		if startIndex != -1 {
			textArray = append(textArray[:startIndex], append([]string{obj.text}, textArray[startIndex+len(obj.line):]...)...)
		}
	}
	return textArray
}

func findLineIndex(line, textArray []string) int {
	for i := 0; i <= len(textArray)-len(line); i++ {
		if strings.EqualFold(strings.Join(textArray[i:i+len(line)], " "), strings.Join(line, " ")) {
			return i
		}
	}
	return -1
}

func keepOnlyMonths(input string) string {
	months := []string{
		"january", "jan", "february", "feb", "march", "mar", "april", "apr", "may", "june", "jun",
		"july", "jul", "august", "aug", "september", "sep", "october", "oct", "november", "nov",
		"december", "dec",
	}
	input = strings.ToLower(input)
	result := ""
	parts := strings.FieldsFunc(input, func(r rune) bool {
		return !unicode.IsLetter(r)
	})
	for _, part := range parts {
		if containsIgnoreCase(months, part) {
			result += part + " "
		}
	}
	return strings.TrimSpace(result)
}

func containsIgnoreCase(slice []string, item string) bool {
	for _, v := range slice {
		if strings.EqualFold(v, item) {
			return true
		}
	}
	return false
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func FormatTime(t *time.Time) string {
	if t == nil {
		return ""
	}
	return t.Format("Mon Jan 02 2006")
}
