package config

import (
	"context"
	"encoding/xml"
	"html"
	"io"
	"net/http"
)

type RSSFeed struct {
	Channel struct {
		Title       string    `xml:"title"`
		Link        string    `xml:"link"`
		Description string    `xml:"description"`
		Item        []RSSItem `xml:"item"`
	} `xml:"channel"`
}

type RSSItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
}

func (r *RSSFeed) CleanStrings() {
	r.Channel.Title = html.UnescapeString(r.Channel.Title)
	r.Channel.Description = html.UnescapeString(r.Channel.Description)
}

func (r *RSSItem) CleanStrings() {
	r.Title = html.UnescapeString(r.Title)
	r.Description = html.UnescapeString(r.Description)
}

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	// setup client
	// make a request with context
	// add gator as a user-agent in the header
	// io.ReadAll
	// xml.UnMarshal

	client := &http.Client{}
	req, err := http.NewRequestWithContext(ctx, "GET", feedURL, nil)
	if err != nil {
		return &RSSFeed{}, err
	}
	req.Header.Add("User-Agent", "gator")

	resp, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return &RSSFeed{}, err
	}

	var feedData RSSFeed
	err = xml.Unmarshal(data, &feedData)
	if err != nil {
		return &RSSFeed{}, err
	}

	feedData.CleanStrings()
	for _, item := range feedData.Channel.Item {
		item.CleanStrings()
	}

	return &feedData, nil
}
