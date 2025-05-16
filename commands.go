package main

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/tepidmilk/gator/internal/config"
	"github.com/tepidmilk/gator/internal/database"
)

type state struct {
	db  *database.Queries
	cfg *config.Config
}

type command struct {
	name string
	args []string
}

type commands struct {
	cmd map[string]func(*state, command) error
}

func (c *commands) run(s *state, cmd command) error {
	val, ok := c.cmd[cmd.name]
	if !ok {
		return errors.New("invalid command name")
	}
	err := val(s, cmd)
	if err != nil {
		return err
	}
	return err
}

func (c *commands) register(name string, f func(*state, command) error) {
	c.cmd[name] = f
}

func handlerLogin(s *state, cmd command) error {

	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == sql.ErrNoRows {
		return errors.New("unable to login. user does not exist")
	} else if err != nil {
		return err
	}

	err = s.cfg.SetUser(cmd.args[0])
	if err != nil {
		return fmt.Errorf("couldn't set current user: %W", err)
	}

	fmt.Printf("User has been set to %s\n", cmd.args[0])
	return err
}

func handlerRegister(s *state, cmd command) error {
	fmt.Println(cmd.args[0])
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %s <name>", cmd.name)
	}

	_, err := s.db.GetUser(context.Background(), cmd.args[0])
	if err == nil {
		return errors.New("user already exists in database")
	} else if err != sql.ErrNoRows {
		return err
	}

	s.db.CreateUser(context.Background(), database.CreateUserParams{ID: uuid.New(), CreatedAt: time.Now().UTC(), UpdatedAt: time.Now().UTC(), Name: cmd.args[0]})
	s.cfg.SetUser(cmd.args[0])

	fmt.Println("User successfully registered in Database")
	return nil
}

func handlerReset(s *state, cmd command) error {
	err := s.db.Reset(context.Background())
	return err
}

func handlerUsers(s *state, cmd command) error {
	users, err := s.db.GetUsers(context.Background())
	if err != nil {
		return err
	}
	for _, user := range users {
		if user == s.cfg.CurrentUserName {
			fmt.Printf("%s (current)\n", user)
		} else {
			fmt.Println(user)
		}
	}
	return err
}

func handlerAgg(s *state, cmd command) error {
	rss, err := fetchFeed(context.Background(), "https://www.wagslane.dev/index.xml")
	if err != nil {
		return err
	}
	fmt.Println(rss)
	return err
}

func handlerAddFeed(s *state, cmd command) error {
	if len(cmd.args) < 2 {
		return fmt.Errorf("usage: %s <Feed Name> <URL>", cmd.name)
	}
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}
	name := cmd.args[0]
	URL := cmd.args[1]
	feed, err := s.db.CreateFeed(context.Background(), database.CreateFeedParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		Name:      name,
		Url:       URL,
		UserID:    user.ID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed: %v", err)
	}
	fmt.Println("Feed Created:", feed.Name)
	err = handlerFollow(s, command{args: []string{URL}})
	if err != nil {
		return err
	}

	return err
}

func handlerFeeds(s *state, cmd command) error {
	feeds, err := s.db.GetFeeds(context.Background())
	if err != nil {
		return fmt.Errorf("error getting feeds: %v", err)
	}
	for _, item := range feeds {
		fmt.Printf(" - '%s' {%s} (user: %s)\n", item.Name, item.Url, item.Name_2)
	}
	return err
}

func handlerFollow(s *state, cmd command) error {
	if len(cmd.args) < 1 {
		return fmt.Errorf("usage: %s <URL>", cmd.name)
	}
	URL := cmd.args[0]
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}
	feedID, err := s.db.GetFeedByURL(context.Background(), URL)
	if err != nil {
		return errors.New("no feed exists at given URL: Try 'addfeed <name> <URL>")
	}
	feedFollow, err := s.db.CreateFeedFollow(context.Background(), database.CreateFeedFollowParams{
		ID:        uuid.New(),
		CreatedAt: time.Now().UTC(),
		UpdatedAt: time.Now().UTC(),
		UserID:    user.ID,
		FeedID:    feedID,
	})
	if err != nil {
		return fmt.Errorf("error creating feed follow record: %v", err)
	}

	fmt.Printf("%s is now following '%s'\n", feedFollow.UserName, feedFollow.FeedName)

	return err
}

func handlerFollowing(s *state, cmd command) error {
	user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
	if err != nil {
		return fmt.Errorf("error getting user: %v", err)
	}
	feedsFollowing, err := s.db.GetFeedFollowsForUser(context.Background(), user.Name)
	if err != nil {
		return fmt.Errorf("error getting followed feeds: %v", err)
	}
	for _, item := range feedsFollowing {
		fmt.Println(item.FeedName)
	}
	return err
}
