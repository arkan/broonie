# Telegram Memories Bot

This Telegram bots automatically downloads all the photos and videos from conversations where the bot has been added to.
This makes it very convenient to save all the memories from family groups for example.

## Usage

You first need to create a bot on Telegram to retrieve a token. And then you can run the following command:

```
docker run -e TOKEN="YOUR-TOKEN" -v $(pwd)/my-data:/app/telegram arkan/telegram_memories_bot:0.1.1
```

## Copyright

See the [LICENSE](./LICENSE) (MIT) file for more details.
