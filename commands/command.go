package commands

import (
	"fmt"
	"github.com/0x111/telegram-rss-bot/feeds"
	"github.com/0x111/telegram-rss-bot/replies"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	log "github.com/sirupsen/logrus"
	"strconv"
	"strings"
)

// Commands functions which are executed upon receiving a command

func AddCommand(Bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	commandArguments := strings.Split(update.Message.CommandArguments(), " ")

	if len(commandArguments) < 2 {
		log.Debug("Not enough arguments\\. We need \"/add name url\"")
		return
	}

	feedName := commandArguments[0]
	feedUrl := commandArguments[1]
	userid := update.Message.From.ID
	chatid := update.Message.Chat.ID

	err := feeds.AddFeed(Bot, feedName, feedUrl, chatid, userid)
	txt := ""

	if err == nil {
		txt = fmt.Sprintf("The feed with the url [%s] was successfully added to this channel\\!", replies.FilterMessageChars(feedUrl))
		replies.SimpleMessage(Bot, chatid, update.Message.MessageID, txt)
	}
}

func UpdateCommand(Bot *tgbotapi.BotAPI, update *tgbotapi.Update) {

	commandArguments := strings.Split(update.Message.CommandArguments(), " ")

	if len(commandArguments) < 2 {
		log.Debug("Not enough arguments\\. We need \"/update id url\"")
		return
	}

	feedUrl := commandArguments[1]
	chatid := update.Message.Chat.ID
	userid := update.Message.From.ID
	feedid, err := strconv.Atoi(commandArguments[0])
	if err != nil {
		txt := fmt.Sprintf("Fail parse to int id %s", commandArguments[0])
		log.Debug(txt)
		replies.SimpleMessage(Bot, chatid, update.Message.MessageID, txt)
		return
	}

	err = feeds.UpdateFeedByID(Bot, feedid, feedUrl, chatid, userid)
	if err != nil {
		replies.SimpleMessage(Bot, chatid, update.Message.MessageID, err.Error())
		return
	}
	txt := ""

	if err == nil {
		txt = fmt.Sprintf("The feed with the url [%s] was successfully updated to this channel\\!", replies.FilterMessageChars(feedUrl))
		replies.SimpleMessage(Bot, chatid, update.Message.MessageID, txt)
	}
}

func ListCommand(Bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	chatid := update.Message.Chat.ID
	userid := update.Message.From.ID
	feedres, err := feeds.ListFeeds(userid, chatid)

	if err != nil {
		log.Debugf("Fail to ListFeeds, got %v", err)
		return
	}

	replies.ListOfFeeds(Bot, feedres, chatid, update.Message.MessageID)
}

func DeleteCommand(Bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	commandArguments := strings.Split(update.Message.CommandArguments(), " ")

	if len(commandArguments) < 1 {
		log.Debug("Not enough arguments\\. We need \"/delete id\"")
		return
	}

	chatid := update.Message.Chat.ID
	userid := update.Message.From.ID
	feedid, err := strconv.Atoi(commandArguments[0])
	if err != nil {
		txt := fmt.Sprintf("Fail parse to int id %s", commandArguments[0])
		log.Debug(txt)
		replies.SimpleMessage(Bot, chatid, update.Message.MessageID, txt)
		return
	}
	err = feeds.DeleteFeedByID(feedid, chatid, userid)

	if err != nil {
		txt := fmt.Sprintf("There is no feed with the id [%d]\\!", feedid)
		replies.SimpleMessage(Bot, chatid, update.Message.MessageID, txt)
		return
	}

	txt := fmt.Sprintf("The feed with the id [%d] was successfully deleted\\!", feedid)
	replies.SimpleMessage(Bot, chatid, update.Message.MessageID, txt)
}

func HelpCommand(Bot *tgbotapi.BotAPI, update *tgbotapi.Update) {
	txt := `
	Avaliable commands:
/add %FeedName %URL - With this you can add a new feed for the current channel, both the name and the url parameters are required
/list - With this command you are able to list all the existing feeds with their ID numbers
/delete %ID - With this command you are able to delete an added feed if you do not need it anymore. The ID parameter is required and you can get it from the /list command 
/update %ID %URL - With this command you are able to update a feed with a new url. The ID parameter and url are required and you can get id from the /list command 
	`
	replies.SimpleMessage(Bot, update.Message.Chat.ID, update.Message.MessageID, replies.FilterMessageChars(txt))
}
