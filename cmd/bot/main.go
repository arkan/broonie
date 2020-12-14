package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	"github.com/arkan/telegram_memories_bot"
)

func main() {
	client, err := telegram_memories_bot.NewClient("config.json")
	if err != nil {
		log.Fatalf("Unable to create new client: %s", err.Error())
	}

	err = client.HandleNewUploads(func(filename string, url string, rule *telegram_memories_bot.ConfigRule) error {
		now := time.Now().UTC()
		p := path.Join(
			rule.Directory,
			fmt.Sprintf("%d", now.Year()),
			"Telegram",
		)

		if rule.SubDirectory != "" {
			p = path.Join(p, rule.SubDirectory)
		}

		if err := os.MkdirAll(p, 0755); err != nil {
			return err
		}

		destination := path.Join(p, filename)

		resp, err := http.Get(url)
		if err != nil {
			return err
		}
		defer resp.Body.Close()

		out, err := os.Create(destination)
		if err != nil {
			return err
		}
		defer out.Close()

		if _, err := io.Copy(out, resp.Body); err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Fatalf("Unable to run HandleNewUploads: %s", err.Error())
	}
}
