package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

type client struct {
	token        string
	outDirectory string
	debug        bool
	bot          *tgbotapi.BotAPI
}

func (c *client) downloadFile(fileID string, groupName string) error {
	cfg := tgbotapi.FileConfig{
		FileID: fileID,
	}
	file, err := c.bot.GetFile(cfg)
	if err != nil {
		return err
	}

	now := time.Now().UTC()
	p := path.Join(
		c.outDirectory,
		fmt.Sprintf("%d", now.Year()),
		"Telegram",
		groupName,
	)

	if err := os.MkdirAll(p, 0755); err != nil {
		return err
	}

	destination := path.Join(p, path.Base(file.FilePath))

	// Get the data
	resp, err := http.Get(file.Link(c.token))
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(destination)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	if _, err := io.Copy(out, resp.Body); err != nil {
		return err
	}

	return nil
}

func main() {
	c := &client{
		token:        os.Getenv("TOKEN"),
		outDirectory: "telegram",
		debug:        os.Getenv("DEBUG") == "true",
	}

	var err error
	c.bot, err = tgbotapi.NewBotAPI(c.token)
	if err != nil {
		log.Panic(err)
	}

	c.bot.Debug = c.debug

	log.Printf("Authorized on account %s - Debug: %t", c.bot.Self.UserName, c.debug)

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := c.bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		log.Printf("\n\nNew update [%s] %s", update.Message.From.UserName, update.Message.Text)

		channelID := update.Message.Chat.Title
		fileID := ""
		if update.Message.Photo != nil {
			log.Printf("[Processing Photo] %#v\n", update.Message.Photo)
			fileID = (*update.Message.Photo)[len(*update.Message.Photo)-1].FileID
		} else if update.Message.Video != nil {
			log.Printf("[Processing Video] %#v\n", update.Message.Video)
			fileID = update.Message.Video.FileID
		} else {
			log.Printf("[Processing not supported] %#v\n", update.Message)
			continue
		}

		if err := c.downloadFile(fileID, channelID); err != nil {
			panic(err)
		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// c.bot.Send(msg)
	}
}
