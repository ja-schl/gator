package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerAggregation(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %v <time_between_reqs>", cmd.Name)
	}
	timeBetweenReqs, err := time.ParseDuration(cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not scrape: %w", err)
	}

	ticker := time.NewTicker(timeBetweenReqs)
	for ; ; <-ticker.C {
		scrapeFeeds(s)
	}
}

func handlerAddFeed(s *state, cmd command, usr database.User) error {
	if len(cmd.Args) != 2 {
		return fmt.Errorf("usage %v <name> <url>", cmd.Name)
	}

	feed := database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      cmd.Args[0],
		Url:       cmd.Args[1],
		UserID:    usr.ID,
	}
	newFeed, err := s.db.CreateFeed(context.Background(), feed)
	if err != nil {
		return fmt.Errorf("could not create feed: %w", err)
	}
	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    usr.ID,
	}
	fmt.Printf("%+v", newFeed)

	_, err = s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	return nil
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error fetching feeds: %w", err)
	}

	for _, feed := range feeds {
		usr, err := s.db.GetUserById(context.Background(), feed.UserID)
		if err != nil {
			return fmt.Errorf("error fetching feeds: %w", err)
		}
		fmt.Println("###############")
		fmt.Printf("Name: %v\nURL: %v\nUserName: %v\n", feed.Name, feed.Url, usr.Name)
	}
	return nil
}

func handlerFollow(s *state, cmd command, usr database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %v <url>", cmd.Name)
	}

	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	params := database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		FeedID:    feed.ID,
		UserID:    usr.ID,
	}
	result, err := s.db.CreateFeedFollow(context.Background(), params)
	if err != nil {
		return fmt.Errorf("could not follow feed: %w", err)
	}

	fmt.Printf("new follow created: %v follows %v", result.UserName, result.FeedName)
	return nil
}

func handlerFollowing(s *state, cmd command, usr database.User) error {
	result, err := s.db.GetFeedFollowsForUser(context.Background(), usr.ID)
	if err != nil {
		return fmt.Errorf("could not get your feeds: %w", err)
	}

	for _, ffollow := range result {
		fmt.Println(ffollow.FeedName)
	}

	return nil
}

func handlerUnfollow(s *state, cmd command, usr database.User) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %v <url>", cmd.Name)
	}
	feed, err := s.db.GetFeedByUrl(context.Background(), cmd.Args[0])
	if err != nil {
		return fmt.Errorf("could not delete follow: %w", err)
	}

	deleteParams := database.DeleteFollowParams{UserID: usr.ID, FeedID: feed.ID}
	err = s.db.DeleteFollow(context.Background(), deleteParams)
	if err != nil {
		return fmt.Errorf("could not delete follow: %w", err)
	}

	return nil
}
