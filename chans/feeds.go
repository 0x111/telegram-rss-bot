package chans

import (
	"fmt"
	"github.com/0x111/telegram-rss-bot/feeds"
	"github.com/0x111/telegram-rss-bot/replies"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
)

// Get Feed Updates from the feeds
func FeedUpdates() {
	feedUpdates := feeds.GetFeedUpdatesChan()

	for feedUpdate := range feedUpdates {
		log.WithFields(log.Fields{"feedData": feedUpdate}).Debug("Requesting feed data")
		log.WithFields(log.Fields{"feedID": feedUpdate.ID, "feedUrl": feedUpdate.Url}).Info("Updating feeds")
		feeds.GetFeed(feedUpdate.Url, feedUpdate.ID)
	}
}

// Post Feed data to the channel
func FeedPosts(Bot *tgbotapi.BotAPI) {
	feedPosts := feeds.PostFeedUpdatesChan()

	for feedPost := range feedPosts {
		link := replies.FilterMessageChars(feedPost.Link)
		msg := fmt.Sprintf("_%s_ \\- *%s* \\- [%s](%s)", replies.FilterMessageChars(feedPost.Name), replies.FilterMessageChars(feedPost.Title), link, link)
		log.WithFields(log.Fields{"feedPost": feedPost, "chatID": feedPost.ChatID}).Debug("Posting feed update to the Telegram API")
		err := replies.SimpleMessage(Bot, feedPost.ChatID, 0, msg)
		if err == nil {
			log.WithFields(log.Fields{"feedPost": feedPost, "chatID": feedPost.ChatID}).Debug("Setting the Feed Data entry to published!")
			_, err := feeds.UpdateFeedDataPublished(&feedPost)
			if err != nil {
				log.WithFields(log.Fields{"error": err, "feedPost": feedPost, "chatID": feedPost.ChatID}).Error("There was an error while updating the Feed Data entry!")
			}
		} else {
			log.WithFields(log.Fields{"error": err, "feedPost": feedPost, "chatID": feedPost.ChatID}).Error("There was an error while posting the update to the feed!")
		}
	}
}
