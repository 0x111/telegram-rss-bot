package replies

import (
	"github.com/0x111/telegram-rss-bot/models"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"strconv"
	"strings"
)

// This function replies the list of feeds to for the command /list
func ListOfFeeds(botAPI *tgbotapi.BotAPI, feeds *[]models.Feed, chatid int64, messageid int) {
	txt := "Here is the list of your added Feeds for this Room: \n"

	if len(*feeds) == 0 {
		txt += "There is currently no feed added to the list for this Room\\!\n"
	}

	for _, feed := range *feeds {
		txt += "[\\#" + strconv.Itoa(feed.ID) + "] *" + FilterMessageChars(feed.Name) + "*: " + FilterMessageChars(feed.Url) + "\n"
	}

	msg := tgbotapi.NewMessage(chatid, txt)
	msg.ReplyToMessageID = messageid

	msg.ParseMode = "markdownv2"
	msg.DisableWebPagePreview = true

	botAPI.Send(msg)
}

func SimpleMessage(botAPI *tgbotapi.BotAPI, chatid int64, messageid int, text string) error {

	msg := tgbotapi.NewMessage(chatid, text)

	if messageid != 0 {
		msg.ReplyToMessageID = messageid
	}

	msg.ParseMode = "markdownv2"
	msg.DisableWebPagePreview = false

	_, err := botAPI.Send(msg)

	if err != nil {
		return err
	}

	return nil
}

func FilterMessageChars(msg string) string {
	var markdownEscaper = strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)

	return markdownEscaper.Replace(msg)
}
