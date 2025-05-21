package main

import (
	"context"
	"database/sql"
	"log"
	"os"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/google/uuid"
	_ "github.com/lib/pq"
	"github.com/tepidmilk/gator/internal/config"
	"github.com/tepidmilk/gator/internal/database"
)

func main() {
	cfg, err := config.Read()
	if err != nil {
		log.Fatalf("error reading config: %v", err)
	}
	db, err := sql.Open("postgres", cfg.DbURL)
	if err != nil {
		log.Fatalf("error opening database: %v", err)
	}
	dbQueries := database.New(db)
	s := state{
		db:  dbQueries,
		cfg: &cfg,
	}
	c := commands{
		cmd: make(map[string]func(*state, command) error),
	}
	c.register("login", handlerLogin)
	c.register("register", handlerRegister)
	c.register("reset", handlerReset)
	c.register("users", handlerUsers)
	c.register("agg", handlerAgg)
	c.register("addfeed", middlewareLoggedIn(handlerAddFeed))
	c.register("feeds", handlerFeeds)
	c.register("follow", middlewareLoggedIn(handlerFollow))
	c.register("following", middlewareLoggedIn(handlerFollowing))
	c.register("unfollow", middlewareLoggedIn(handlerUnfollow))
	c.register("browse", middlewareLoggedIn(handlerBrowse))
	args := os.Args
	if len(args) < 2 {
		log.Fatal("Usage: cli <command> [args...]")
		return
	}
	cmd := command{
		name: args[1],
		args: args[2:],
	}
	err = c.run(&s, cmd)
	if err != nil {
		log.Fatal(err)
	}
}

// takes a handler that requires a logged in user and returns a normal handler to register
func middlewareLoggedIn(handler func(s *state, cmd command, user database.User) error) func(*state, command) error {
	return func(s *state, cmd command) error {
		user, err := s.db.GetUser(context.Background(), s.cfg.CurrentUserName)
		if err != nil {
			return err
		}

		return handler(s, cmd, user)
	}
}

func scrapeFeeds(s *state) {
	nextFeed, err := s.db.GetNextFeedToFetch(context.Background())
	if err != nil {
		log.Println("error getting next feed to fetch:", err)
		return
	}

	err = s.db.MarkFeedFetched(context.Background(), database.MarkFeedFetchedParams{
		LastFethcedAt: sql.NullTime{Time: time.Now().UTC(), Valid: true},
		ID:            nextFeed.ID,
	})
	if err != nil {
		log.Println("error marking feed fetched:", err)
		return
	}

	RSSFeed, err := fetchFeed(context.Background(), nextFeed.Url)
	if err != nil {
		log.Println("error fetching feed:", err)
		return
	}

	for _, item := range RSSFeed.Channel.Item {
		publishTime, err := dateparse.ParseAny(item.PubDate)
		if err != nil {
			log.Println("error parsing item's time:", err)
			return
		}
		err = s.db.CreatePost(context.Background(), database.CreatePostParams{
			ID:        uuid.New(),
			CreatedAt: time.Now().UTC(),
			UpdatedAt: time.Now().UTC(),
			Title:     item.Title,
			Url:       item.Link,
			Description: sql.NullString{
				String: item.Description,
				Valid:  true,
			},
			PublishedAt: sql.NullTime{
				Time:  publishTime,
				Valid: true,
			},
			FeedID: nextFeed.ID,
		})
		if err != nil {
			if strings.Contains(err.Error(), "duplicate key value violates unique constraint") {
				continue
			}
		}
		log.Printf("couldn't create post: %v", err)
		continue
	}
	log.Printf("Feed %s collectedm %v posts found", nextFeed.Name, len(RSSFeed.Channel.Item))
}
