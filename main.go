package main

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/turnage/graw/reddit"
	"github.com/turnage/graw/streams"
)

func main() {
	bot, err := reddit.NewScript("andrew-test", 2*time.Second)

	// bot, err := reddit.NewBot(
	// 	reddit.BotConfig{
	// 		Agent: "poy-bot",
	// 		App: reddit.App{
	// 			ID:       "4egNiRc70biolA",
	// 			Secret:   "z_sWcTg5cJWLVFon4H5PJBvtQyI",
	// 			Username: "poy_",
	// 			Password: "asdfasdf",
	// 		},
	// 		Rate: time.Second,
	// 	},
	// )
	if err != nil {
		log.Fatalf("failed to make bot: %s", err)
	}

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
			req, err := http.NewRequest(
				"POST",
				"https://hooks.slack.com/services/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",
				strings.NewReader(fmt.Sprintf(`{"text":"%s: https://reddit.com%s", "attachments":[{"image_url":%q, "title":"post"}]}`, post.Title, post.Permalink, post.Thumbnail)),
			)
			if err != nil {
				log.Fatalf("failed to build request: %s", err)
			}

			resp, err := http.DefaultClient.Do(req)
			if err != nil {
				log.Fatalf("failed to POST request: %s", err)
			}

			if resp.StatusCode != http.StatusOK {
				log.Printf("unexpected status code: %d", resp.StatusCode)
			}
		}
	}
}
