package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"strconv"
)

func handlerBrowse(s *state, cmd command, usr database.User) error {
	limit := 2
	if len(cmd.Args) == 1 {
		var err error
		limit, err = strconv.Atoi(cmd.Args[0])
		if err != nil {
			return fmt.Errorf("browse failed: %w", err)
		}
	}
	posts, err := s.db.GetPostsForUser(context.Background(), database.GetPostsForUserParams{UserID: usr.ID, Limit: int32(limit)})
	if err != nil {
		return fmt.Errorf("browse failed: %w", err)
	}
	for _, p := range posts {
		fmt.Printf("Title: %s\t URL: %s\n", p.Title, p.Url)
	}
	return nil
}
