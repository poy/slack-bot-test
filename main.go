package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"
	"github.com/turnage/graw/streams"
)

func main() {
	bot, err := reddit.NewScript("andrew-test", 2*time.Second)
	if err != nil {
		log.Fatalf("failed to make bot: %s", err)
	}

	subreddits := []string{"funny"}
	path := "/r/" + strings.Join(subreddits, "+") + "/new"
	h, err := bot.Listing(path, os.Getenv("REDDIT_REF"))
	if err != nil {
		log.Fatalf("failed to List: %s", err)
	}

	for _, c := range h.Comments {
		fmt.Printf("%+v\n", c)
	}

	for _, p := range h.Posts {
		fmt.Printf("%+v\n", p)
	}

	for _, m := range h.Messages {
		fmt.Printf("%+v\n", m)
	}

	return

	kill := make(chan bool)
	defer close(kill)
	errs := make(chan error)
	posts, err := streams.Subreddits(bot, kill, errs, "funny")
	if err != nil {
		log.Fatalf("failed to stream: %s", err)
	}

	for {
		select {
		case err := <-errs:
			log.Fatalf("failed while reading: %s", err)
		case post := <-posts:
			log.Printf("%+v", post)
		}
	}
}
