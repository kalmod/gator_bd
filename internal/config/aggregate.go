package config

import (
	"blog_agg_2/internal/database"
	"context"
	"database/sql"
	"fmt"
	"html"
	"strings"
	"time"

	"github.com/google/uuid"
)

func scrapeFeeds(s *State) error {
	next_feed, err := s.Db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return err
	}

	feed_to_update := database.MarkFeedFetchedParams{LastFetchedAt: sql.NullTime{Time: time.Now(), Valid: true}, UpdatedAt: time.Now(), ID: next_feed.ID}
	err = s.Db.MarkFeedFetched(context.Background(), feed_to_update)
	if err != nil {
		return err
	}

	fetched_feed, err := fetchFeed(context.Background(), next_feed.Url)
	if err != nil {
		return err
	}

	for _, feed_item := range fetched_feed.Channel.Item {
		pubDate_parseTime, err := time.Parse(time.RFC1123, feed_item.PubDate)
		if err != nil {
			fmt.Printf("Error parsing feed pubdate %v\n", err)
			return err
		}

		new_post := database.CreatePostParams{ID: uuid.New(), CreatedAt: time.Now(), UpdatedAt: time.Now(),
			Title: feed_item.Title, Url: feed_item.Link,
			Description: sql.NullString{String: html.UnescapeString(feed_item.Description), Valid: true},
			PublishedAt: pubDate_parseTime, FeedID: next_feed.ID}

		_, post_err := s.Db.CreatePost(context.Background(), new_post)
		if post_err != nil && !strings.Contains(post_err.Error(), "pq: duplicate key value violates unique constraint") {
			fmt.Println(feed_item.Title, feed_item.Link)
			return post_err
		}
	}

	return nil
}
