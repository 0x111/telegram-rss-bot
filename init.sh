#!/bin/bash

set -ue

keys="TELEGRAM_AUTH_KEY MIGRATIONS TELEGRAM_API_DEBUG DB_PATH LOG_LEVEL FEED_PARSE_AMOUNT FEED_POST_AMOUNT FEED_UPDATES_INTERVAL FEED_POST_INTERVAL"

TELEGRAM_AUTH_KEY=${TELEGRAM_AUTH_KEY:=token}
MIGRATIONS=${MIGRATIONS:=v1}
TELEGRAM_API_DEBUG=${TELEGRAM_API_DEBUG:=false}
DB_PATH=${DB_PATH:=./bot.db}
LOG_LEVEL=${LOG_LEVEL:=info}
FEED_PARSE_AMOUNT=${FEED_PARSE_AMOUNT:=5}
FEED_POST_AMOUNT=${FEED_POST_AMOUNT:=2}
FEED_UPDATES_INTERVAL=${FEED_UPDATES_INTERVAL:=600}
FEED_POST_INTERVAL=${FEED_POST_INTERVAL:=400}

mkdir -p /code/.telegram-rss-bot
cp ./bot-config.template.json /code/.telegram-rss-bot/bot-config.json

for key in $keys; do
  sed -i "s|$key|$(eval echo "\$${key}")|g" /code/.telegram-rss-bot/bot-config.json
done

echo Executing $*
$*
