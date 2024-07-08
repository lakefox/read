package utils

import (
	"main/models"
	"math/rand"
)

func ShuffleArticles(arr []models.Article) {
	for i := len(arr) - 1; i > 0; i-- {
		j := rand.Intn(i + 1)
		arr[i], arr[j] = arr[j], arr[i]
	}
}

func Min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
