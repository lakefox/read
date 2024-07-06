package main

import (
	"time"

	"github.com/mmcdole/gofeed"
)

// FeedItem represents a parsed feed item with the necessary fields
type FeedItem struct {
	Title       string
	Link        string
	Description string
	Published   time.Time
	Author      string
}

// ParseRSSFeed parses the RSS feed from the given URL and returns a slice of FeedItems
func ParseRSSFeed(url string) ([]FeedItem, error) {
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(url)
	if err != nil {
		return nil, err
	}

	var items []FeedItem
	for _, item := range feed.Items {
		author := ""
		if item.Author != nil {
			author = item.Author.Name
		} else if len(item.Authors) > 0 {
			author = item.Authors[0].Name
		}

		pubDate := time.Time{}
		if item.PublishedParsed != nil {
			pubDate = *item.PublishedParsed
		}

		feedItem := FeedItem{
			Title:       item.Title,
			Link:        item.Link,
			Description: item.Description,
			Published:   pubDate,
			Author:      author,
		}
		items = append(items, feedItem)
	}

	return items, nil
}

// func main() {
// 	url := "https://www.reddit.com/r/longreads.rss" // Replace with your RSS feed URL
// 	// url := "https://rss.nytimes.com/services/xml/rss/nyt/World.xml" // Replace with your RSS feed URL
// 	items, err := ParseRSSFeed(url)
// 	if err != nil {
// 		log.Fatalf("Error parsing RSS feed: %v", err)
// 	}

// 	for _, item := range items {
// 		fmt.Printf("Title: %s\nLink: %s\nDescription: %s\nPublished: %s\nAuthor: %s\n\n",
// 			item.Title, item.Link, item.Description, item.Published, item.Author)
// 	}
// }
