package recommend

import (
	"sort"
)

type Article struct {
	Keywords    []string
	Voice       string
	Category    string
	Title       string
	Author      string
	Image       string
	Description string
	Site        string
	Url         string
	Text        string
	Date        string
}

func Rank(newArticles, oldArticles []Article) []Article {
	// Filter out articles with duplicate URLs
	filteredArticles := filterDuplicateUrls(newArticles, oldArticles)

	// Calculate similarity scores
	type ScoredArticle struct {
		Article Article
		Score   int
	}
	scoredArticles := []ScoredArticle{}
	for _, newArticle := range filteredArticles {
		score := calculateMaxSimilarity(newArticle, oldArticles)
		scoredArticles = append(scoredArticles, ScoredArticle{Article: newArticle, Score: score})
	}

	// Sort articles by score in descending order
	sort.Slice(scoredArticles, func(i, j int) bool {
		return scoredArticles[i].Score > scoredArticles[j].Score
	})

	// Extract sorted articles
	sortedArticles := []Article{}
	for _, scoredArticle := range scoredArticles {
		sortedArticles = append(sortedArticles, scoredArticle.Article)
	}

	return sortedArticles
}

func filterDuplicateUrls(newArticles, oldArticles []Article) []Article {
	oldUrls := make(map[string]bool)
	for _, article := range oldArticles {
		oldUrls[article.Url] = true
	}

	filteredArticles := []Article{}
	for _, article := range newArticles {
		if !oldUrls[article.Url] {
			filteredArticles = append(filteredArticles, article)
		}
	}

	return filteredArticles
}

func calculateMaxSimilarity(newArticle Article, oldArticles []Article) int {
	maxScore := 0
	for _, oldArticle := range oldArticles {
		score := calculateSimilarity(newArticle, oldArticle)
		if score > maxScore {
			maxScore = score
		}
	}
	return maxScore
}

func calculateSimilarity(a, b Article) int {
	score := 0

	// Compare Keywords
	keywordSet := make(map[string]bool)
	for _, keyword := range a.Keywords {
		keywordSet[keyword] = true
	}
	for _, keyword := range b.Keywords {
		if keywordSet[keyword] {
			score += 1
		}
	}

	// Compare Category
	if a.Category == b.Category {
		score += 3
	}

	// Compare Author
	if a.Author == b.Author {
		score += 2
	}

	// Compare Site
	if a.Site == b.Site {
		score += 1
	}

	// Additional criteria can be added here

	return score
}
