package main

import(
	"fmt"
	"context"
	"io"
	"encoding/xml"
	"time"
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

func fetchFeed(ctx context.Context, feedURL string) (*RSSFeed, error) {
	
	client := &http.Client{
		Timeout: time.Second * 5,
	}

	req, err := http.NewRequestWithContext(
		ctx, "GET", feedURL,
		nil)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error: Creating request %w ", err)
	}
	req.Header.Set("User-Agent", "gator")
	res, err := client.Do(req)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error: making the request %w ", err)
	}
	body, err := io.ReadAll(res.Body)
	defer res.Body.Close()
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error: reading response %w ", err)
	}

	var feed RSSFeed
	err = xml.Unmarshal(body, &feed)
	if err != nil {
		return &RSSFeed{}, fmt.Errorf("Error: Unmarshaling %w ", err)
	}

	return &feed, nil
	
}
