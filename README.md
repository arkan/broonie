# Broonie

Broonie([why this name?](https://en.wikipedia.org/wiki/Brownie_(folklore))) is a Telegram bots that automatically downloads photos and videos in real-time from conversations where the bot has been added to.
This makes it very convenient to save all the memories from family Telegram groups and expose them on Plex for example.

## Usage

You first need to create [a bot on Telegram](https://t.me/botfather) to retrieve a token.

Then, you need to create your configuration file:

```bash
$> cp config.json.sample config.json
```

Configure it with your own values. You can only run one instance of the bot in parallel.
So if you want to handle multiple Telegram groups, you have to define several entries in the config file.

Then run the application:
```
$> go run cmd/broonie/main.go
```

You should expect to see no error. Enjoy!

### Docker 
A Docker image is also available:

```
docker run -v config.json:/app/config.json:ro arkan/broonie:latest
```

## Copyright

See the [LICENSE](./LICENSE) (MIT) file for more details.
