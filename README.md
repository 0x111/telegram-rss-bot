# telegram-rss-bot

![Build telegram-rss-bot](https://github.com/0x111/telegram-rss-bot/workflows/Build%20telegram-rss-bot/badge.svg)

## Introduction
This is an another telegram bot for usage with RSS feeds.

## First steps
To use this, you first need to register a telegram bot by reading the documentation: https://core.telegram.org/bots#3-how-do-i-create-a-bot

For this bot to work, you will need a token which authorizes you to use the telegram api.

## First run
After you have the token, you should create a copy of the sample file and fill it out accordingly.

## Configuration

Rename `bot-config.sample.json` to `bot-config.json`.

You can put the config file in the current folder on where the binary resides or put it in this folder if created `$HOME/.telegram-rss-bot`, the app should be able to find it here too if you want to have some fixed location for your configuration files.

```json
{
  "telegram_auth_key": "token",
  "migrations": "v1",
  "telegram_api_debug": false,
  "db_path": "./bot.db",
  "log_level": "info",
  "feed_parse_amount": 5,
  "feed_post_amount": 2,
  "feed_updates_interval": 600,
  "feed_posts_interval": 400
}
```

- telegram_auth_key: is the token which you've got from registering the bot
- migrations: this will be used mostly in the future, to define which migration to run (this can change)
- telegram_api_debug: if this is turned on, you will see debug messages from the telegram api on your stdout
- db_path: this contains a path to the db file, if this does not exist, it will be created if the app will have permission to do that
- log_level: with this, you can set the log level to display, the app is using logrus for logging, so this is accepting all the values from this url https://github.com/sirupsen/logrus#level-logging
- feed_parse_amount: this represents the amount of how much of the feed items should be parsed from the provided feed url (e.g. if you provide a url, which has 10 items, then only the 5 latest will be saved to our database, you can alter this value if needed)
- feed_post_amount: this represents the amount of how much of the parsed feed data should be posted to their respective channels (e.g. if you set this to 2, the bot will post every $feed_posts_interval only 2 entries, you can alter value if needed)
- feed_updates_interval: this represents the interval at which rate the feeds saved in the database should be updated in seconds (60 = 60 seconds and so on)
- feed_posts_interval: this represents the interval at which rate the feeds should be posted to their respective channels in seconds (60 = 60 seconds and so on)

## Docker support
You can also run this application as a docker container.

### Docker hub

You can pull the official docker image
```bash
docker pull ruthless/telegram-rss-bot
```

### Build from source
Execute the following steps:
```
git clone https://github.com/0x111/telegram-rss-bot
docker build -t telegram-rss-bot:latest .
docker run --name telegram-rss-bot -e TELEGRAM_AUTH_KEY="MY-TOKEN" -d telegram-rss-bot:latest
```

## Important
Advisory: You should respect the rate limiting of the Telegram API (More info about this: https://core.telegram.org/bots/faq#my-bot-is-hitting-limits-how-do-i-avoid-this)

Feel free to open a PR if you find some bugs or have improvements (I am sure there can be many of those :))

If you find bugs but you have no idea how to fix them, please open an issue with a detailed description on how to reproduce the bug.