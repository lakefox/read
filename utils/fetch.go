package utils

import (
	"fmt"
	"main/models"
	"main/plugins/category"
	fixtext "main/plugins/date"
	"main/plugins/gender"
	"main/plugins/keywords"
	"time"

	readability "github.com/go-shiori/go-readability"
)

func FetchDoc(url string) (models.Article, error) {
	article, err := readability.FromURL(url, 30*time.Second)
	if err != nil {
		return models.Article{}, fmt.Errorf("failed to parse %s, %v", url, err)
	}
	text := article.TextContent
	profile := gender.Profile(text)
	return models.Article{
		Keywords:    keywords.ExtractKeywords(text, 3),
		Voice:       profile["gender"],
		Category:    category.GetCategory(text),
		Title:       article.Title,
		Author:      article.Byline,
		Image:       article.Image,
		Description: article.Excerpt,
		Site:        article.SiteName,
		Url:         url,
		Date:        fixtext.FormatTime(article.PublishedTime),
		Text:        fixtext.ReplaceDates(text),
	}, nil
}
