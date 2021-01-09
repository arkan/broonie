package broonie

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

// Client ...
type Client struct {
	config *Config
	bot    *tgbotapi.BotAPI
}

// ConfigRule contains a channel configuration.
type ConfigRule struct {
	Group               string
	SkipTelegramChannel bool
	Directory           string
	SubDirectory        string
}

// Config contains the bot configuration.
type Config struct {
	Token string
	Debug bool
	Rules []*ConfigRule
}

// NewClient ...
func NewClient(configPath string) (*Client, error) {
	c := &Client{config: &Config{}}
	b, err := ioutil.ReadFile("config.json")
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(b, c.config); err != nil {
		return nil, err
	}

	c.bot, err = tgbotapi.NewBotAPI(c.config.Token)
	if err != nil {
		return nil, err
	}

	c.bot.Debug = c.config.Debug

	log.Printf("Authorized on account %s", c.bot.Self.UserName)
	log.Printf("%#v\n", c.config)

	return c, nil
}

// HandleNewUploads ...
func (c *Client) HandleNewUploads(fn func(filename string, url string, rule *ConfigRule) error) error {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates, err := c.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}

	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		channelID := update.Message.Chat.Title

		rule, ok := c.matchingRule(channelID)
		if !ok {
			log.Printf("[Skip] channel %q doesn't match any rule", channelID)
			continue
		}

		fmt.Printf("Matching rule: %#v\n", rule)
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

		cfg := tgbotapi.FileConfig{
			FileID: fileID,
		}
		file, err := c.bot.GetFile(cfg)
		if err != nil {
			log.Printf("[failed] Upload GetFile failed: %s", err.Error())
			continue
		}

		url := file.Link(c.config.Token)
		filename := path.Base(file.FilePath)
		if err := fn(filename, url, rule); err != nil {
			log.Printf("[failed] Upload processing failed: %s", err.Error())
			continue
		}

		// msg := tgbotapi.NewMessage(update.Message.Chat.ID, update.Message.Text)
		// msg.ReplyToMessageID = update.Message.MessageID

		// c.bot.Send(msg)
	}

	return nil
}

func (c *Client) matchingRule(match string) (*ConfigRule, bool) {
	for _, r := range c.config.Rules {
		if r.Group == match {
			return r, true
		}
	}

	return nil, false
}
