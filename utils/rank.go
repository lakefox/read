package utils

import (
	"sort"

	"main/models"
)

func Rank(newArticles, oldArticles []models.Article) []models.Article {
	filteredArticles := filterDuplicateUrls(newArticles, oldArticles)

	type ScoredArticle struct {
		Article models.Article
		Score   int
	}
	scoredArticles := []ScoredArticle{}
	for _, newArticle := range filteredArticles {
		score := calculateMaxSimilarity(newArticle, oldArticles)
		scoredArticles = append(scoredArticles, ScoredArticle{Article: newArticle, Score: score})
	}

	sort.Slice(scoredArticles, func(i, j int) bool {
		return scoredArticles[i].Score > scoredArticles[j].Score
	})

	sortedArticles := []models.Article{}
	for _, scoredArticle := range scoredArticles {
		sortedArticles = append(sortedArticles, scoredArticle.Article)
	}

	return sortedArticles
}

func filterDuplicateUrls(newArticles, oldArticles []models.Article) []models.Article {
	oldUrls := make(map[string]bool)
	for _, article := range oldArticles {
		oldUrls[article.Url] = true
	}

	filteredArticles := []models.Article{}
	for _, article := range newArticles {
		if !oldUrls[article.Url] {
			filteredArticles = append(filteredArticles, article)
		}
	}

	return filteredArticles
}

func calculateMaxSimilarity(newArticle models.Article, oldArticles []models.Article) int {
	maxScore := 0
	for _, oldArticle := range oldArticles {
		score := calculateSimilarity(newArticle, oldArticle)
		if score > maxScore {
			maxScore = score
		}
	}
	return maxScore
}

func calculateSimilarity(a, b models.Article) int {
	score := 0

	keywordSet := make(map[string]bool)
	for _, keyword := range a.Keywords {
		keywordSet[keyword] = true
	}
	for _, keyword := range b.Keywords {
		if keywordSet[keyword] {
			score += 1
		}
	}

	if a.Category == b.Category {
		score += 3
	}

	if a.Author == b.Author {
		score += 2
	}

	if a.Site == b.Site {
		score += 1
	}

	return score
}
