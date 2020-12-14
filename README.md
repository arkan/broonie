# Telegram Memories Bot

This Telegram bots automatically downloads all the photos and videos in real-time from conversations where the bot has been added to.
This makes it very convenient to save all the memories from family groups for example.

## Usage

You first need to create [a bot on Telegram](https://t.me/botfather) to retrieve a token.

Then, you need to create your configuration file:

```bash
$> cp config.json.sample config.json
```

Configure it with your own values. You can only run one instance of the bot in parallel.
So if you want to handle multiple Telegram groups, you can define several entries in the config file.

Then run the application:
```
$> go run cmd/bot/main.go
```

You should expect to see no error. Enjoy!
## Copyright

See the [LICENSE](./LICENSE) (MIT) file for more details.
