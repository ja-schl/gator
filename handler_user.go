package main

import (
	"context"
	"fmt"
	"gator/internal/database"
	"time"

	"github.com/google/uuid"
)

func handlerLogin(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s <name>", cmd.Name)
	}
	username := cmd.Args[0]
	_, err := s.db.GetUser(context.Background(), username)
	if err != nil {
		return fmt.Errorf("could not get user from db: %w", err)
	}

	if err := s.cfg.SetUser(username); err != nil {
		return fmt.Errorf("could not set user name: %w", err)
	}
	fmt.Println("user name has been set to", s.cfg.CurrentUser)
	return nil
}

func handlerRegister(s *state, cmd command) error {
	if len(cmd.Args) != 1 {
		return fmt.Errorf("usage %s <name>", cmd.Name)
	}

	name := cmd.Args[0]

	newUser := database.CreateUserParams{
		ID:        uuid.New(),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Name:      name}

	usr, err := s.db.CreateUser(context.Background(), newUser)
	if err != nil {
		return err
	}

	if err := s.cfg.SetUser(usr.Name); err != nil {
		return fmt.Errorf("could not set user name: %w", err)
	}
	fmt.Println("user created", usr)

	return nil
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return fmt.Errorf("could not get users: %w", err)
	}

	for _, usr := range users {
		result := usr.Name
		if usr.Name == s.cfg.CurrentUser {
			result = result + " (current)"
		}
		fmt.Println("*", result)
	}
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.DeleteUsers(context.Background())
	if err != nil {
		return fmt.Errorf("error deleting users: %w", err)
	}
	fmt.Println("all users deleted")
	return nil
}
