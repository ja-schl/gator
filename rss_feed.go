package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"gator/internal/database"
	"html"
	"io"
	"net/http"
	"time"

	"github.com/google/uuid"
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

func fetchFeed(ctx context.Context, feedURl string) (*RSSFeed, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, feedURl, nil)
	if err != nil {
		return nil, err
	}
	req.Header.Set("user-agent", "gator")

	client := http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	dat, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var feed RSSFeed
	if err = xml.Unmarshal(dat, &feed); err != nil {
		return nil, err
	}

	feed.Channel.Title = html.UnescapeString(feed.Channel.Title)
	feed.Channel.Description = html.UnescapeString(feed.Channel.Description)
	for i, itm := range feed.Channel.Item {
		itm.Title = html.UnescapeString(itm.Title)
		itm.Description = html.UnescapeString(itm.Description)
		feed.Channel.Item[i] = itm
	}
	return &feed, nil
}

func scrapeFeeds(s *state) error {
	nxtFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		return fmt.Errorf("could scrape feed: %w", err)
	}
	err = s.db.MarkFeedFetched(context.Background(), nxtFeed.ID)
	if err != nil {
		return fmt.Errorf("could scrape feed: %w", err)
	}
	rssFeed, err := fetchFeed(context.Background(), nxtFeed.Url)
	if err != nil {
		return fmt.Errorf("could scrape feed: %w", err)
	}

	for _, itm := range rssFeed.Channel.Item {
		createPost := database.CreatePostParams{
			ID:          uuid.New(),
			CreatedAt:   time.Now(),
			UpdatedAt:   time.Now(),
			Title:       itm.Title,
			Url:         itm.Link,
			Description: itm.Description,
			PublishedAt: parseTime(itm.PubDate),
			FeedID:      nxtFeed.ID,
		}
		_, err := s.db.CreatePost(context.Background(), createPost)
		if err != nil {
			fmt.Println(err)
		}

	}
	return nil

}

func parseTime(s string) time.Time {
	format := "Mon, 2 Jan 2006 15:04:05 -0700"
	t, err := time.Parse(format, s)
	if err != nil {
		return time.Time{}
	}
	return t
}
