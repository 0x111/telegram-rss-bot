package feeds

import (
	"database/sql"
	"errors"
	"github.com/0x111/telegram-rss-bot/conf"
	"github.com/0x111/telegram-rss-bot/db"
	"github.com/0x111/telegram-rss-bot/models"
	"github.com/0x111/telegram-rss-bot/replies"
	"github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/mmcdole/gofeed"
	log "github.com/sirupsen/logrus"
	"time"
)

// Add a new feed to the database
func AddFeed(Bot *tgbotapi.BotAPI, name string, url string, chatid int64, userid int) error {
	DB := db.GetDB()
	exists, err := Exists(url, chatid)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		replies.SimpleMessage(Bot, chatid, 0, "There was an error while adding your feed\\! Please try again\\!")
		return err
	}

	// Check if user is providing a valid feed URL
	if err := isValidFeed(url); err != nil {
		log.WithFields(log.Fields{"error": err}).Debug("Invalid feed!")
		replies.SimpleMessage(Bot, chatid, 0, "The feed you are trying to add is not a valid feed URL\\!")
		return errors.New("invalid_feed_url")
	}

	if exists {
		log.WithFields(log.Fields{"exists": exists}).Debug("Feed exists!")
		replies.SimpleMessage(Bot, chatid, 0, "The feed you are trying to add already exists\\!")
		return errors.New("feed_exists")
	}

	stmt, err := DB.Prepare("INSERT INTO feeds(name, url, chatid, userid) VALUES(?,?,?,?)")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while preparing the query!")
		return err
	}

	_, err = stmt.Exec(name, url, chatid, userid)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while executing the query!")
		return err
	}

	log.Debug("Feed added successfully!")

	return nil
}

// Delete a feed by a feedID parameter
func DeleteFeedByID(feedid int, chatid int64, userid int) error {
	if err := FeedExistsByID(feedid, chatid, userid); err != nil {
		return err
	}

	DB := db.GetDB()

	stmt, err := DB.Prepare("DELETE FROM feeds WHERE id=? AND userid=? AND chatid=?")

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while preparing the query!")
		return err
	}

	_, err = stmt.Exec(feedid, userid, chatid)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while executing the query!")
		return err
	}

	return nil
}

// Get a list of all Feeds added to the database
func ListAllFeeds() (*[]models.Feed, error) {
	DB := db.GetDB()
	rows, err := DB.Query("SELECT id, name, url, userid, chatid FROM feeds")
	defer rows.Close()

	if err != nil {
		if err != sql.ErrNoRows {
			log.WithFields(log.Fields{"error": err}).Error("There was an error with the query!")
			return nil, err
		}
	}

	var feed []models.Feed

	for rows.Next() {
		f := models.Feed{}
		err := rows.Scan(&f.ID, &f.Name, &f.Url, &f.UserID, &f.ChatID)
		feed = append(feed, f)

		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while iterating through the results!")
			return nil, err
		}
	}

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Debug("No results in resultset!")
		return nil, err
	}

	return &feed, nil
}

// List all feeds added to the database, filter by userID and chatID
func ListFeeds(userid int, chatid int64) (*[]models.Feed, error) {
	DB := db.GetDB()
	rows, err := DB.Query("SELECT id, name, url FROM feeds WHERE userid=? AND chatid=?", userid, chatid)
	defer rows.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return nil, err
	}

	var feed []models.Feed

	for rows.Next() {
		f := models.Feed{}
		err := rows.Scan(&f.ID, &f.Name, &f.Url)
		feed = append(feed, f)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while iterating through the results!")
			return nil, err
		}
	}

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Debug("No results in resultset!")
		return nil, err
	}

	return &feed, nil
}

// Get all unpublished feed data from the database, this will contain all feeds
// which were not yet posted to the specific channels
func GetAllUnPublishedFeedData() (*[]models.FeedData, error) {
	DB := db.GetDB()
	config := conf.GetConfig()
	rows, err := DB.Query("SELECT feedData.id, feedData.feedid, feedData.title, feedData.link, feedData.published, feedData.publishedDate, feeds.chatid FROM feedData INNER JOIN feeds on feedData.feedid = feeds.id WHERE feedData.published=? ORDER BY feedData.id ASC LIMIT ?", false, config.GetInt("feed_post_amount"))
	defer rows.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return nil, err
	}

	var feedData []models.FeedData

	for rows.Next() {
		f := models.FeedData{}

		err := rows.Scan(&f.ID, &f.FeedID, &f.Title, &f.Link, &f.Published, &f.PublishedDate, &f.ChatID)
		feedData = append(feedData, f)
		if err != nil {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while iterating through the results!")
			return nil, err
		}
	}

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Debug("No results in resultset!")
		return nil, err
	}

	return &feedData, nil
}

// Check for a feed if it exists by its url
// returns true if exists and false if it does not
func Exists(url string, chatid int64) (bool, error) {
	DB := db.GetDB()
	var err error
	stmt, err := DB.Prepare("SELECT id, name, url, chatid, userid FROM feeds WHERE url=? AND chatid=? LIMIT 1")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
		return false, err
	}

	feed := models.Feed{}

	err = stmt.QueryRow(url, chatid).Scan(&feed.ID, &feed.Name, &feed.Url, &feed.ChatID, &feed.UserID)

	if err != nil {
		log.WithFields(log.Fields{"warn": err}).Warn("The link was not found in the database!")
		if err == sql.ErrNoRows {
			// no rows found, it does not exist
			return false, nil
		} else {
			log.WithFields(log.Fields{"error": err}).Error("There was an error while querying the database!")
			return false, err
		}
	}

	return true, nil
}

func FeedExistsByID(id int, chatid int64, userid int) error {
	DB := db.GetDB()
	feed := models.Feed{}
	stmt, err := DB.Prepare("SELECT id FROM feeds WHERE id=? AND chatid=? AND userid=? LIMIT 1")

	err = stmt.QueryRow(id, chatid, userid).Scan(&feed.ID)

	if err != nil {
		return err
	}

	return nil
}

// This channel is used for getting the feed updates from the database
// Here we send all the feed data to the channel
func GetFeedUpdatesChan() chan models.Feed {
	ch := make(chan models.Feed)
	config := conf.GetConfig()
	ticker := time.NewTicker(config.GetDuration("feed_updates_interval") * time.Second)
	log.Info("Getting feed updates")

	go func() {
		for {
			select {
			case <-ticker.C:
				log.Debug("Ticker triggered!")

				feeds, err := ListAllFeeds()

				if err != nil {
					log.WithFields(log.Fields{"error": err}).Error("There was an error while getting the list of feeds!")
				}

				for _, feed := range *feeds {
					log.WithFields(log.Fields{"feed": feed}).Debug("Writing feed into the channel!")
					ch <- feed
				}
			case <-ch:
				ticker.Stop()
				log.Debug("Stopping ticker")
				return
			}
		}
	}()

	return ch
}

// This channel is use for getting all the unpublished feed data from the database
// We write all the unpublished feed data into the channel
func PostFeedUpdatesChan() chan models.FeedData {
	config := conf.GetConfig()
	ch := make(chan models.FeedData)
	ticker := time.NewTicker(config.GetDuration("feed_posts_interval") * time.Second)
	go func() {
		for {
			select {
			case <-ticker.C:
				feedDatas, err := GetAllUnPublishedFeedData()

				if err != nil {
					log.WithFields(log.Fields{"error": err}).Error("There was an error while getting all unpublished feeds data!")
				}

				for _, feedData := range *feedDatas {
					log.WithFields(log.Fields{"feedData": feedData}).Debug("Writing updates to the channel!")
					ch <- feedData
				}
			case <-ch:
				ticker.Stop()
				log.Debug("Stopping ticker")
				return
			}
		}
	}()

	return ch
}

// This function requests the RSS Feed, parses and processes the data
func GetFeed(feedUrl string, feedID int) error {
	config := conf.GetConfig()
	fp := gofeed.NewParser()
	feed, err := fp.ParseURL(feedUrl)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while parsing the feed!")
		log.Debug("There was an error while parsing the feed!")
		return err
	}

	for index, feedItem := range feed.Items {
		if index > config.GetInt("feed_parse_amount") {
			log.WithFields(log.Fields{"feed_parse_amount": config.GetInt("feed_parse_amount")}).Debug("The amount of feeds to save has been reached, do not parse further!")
			continue
		}

		if !LinkExists(feedItem.Link) {
			WriteFeedData(feedItem, feedID)
		}
	}
	return nil
}

func isValidFeed(url string) error {
	fp := gofeed.NewParser()
	_, err := fp.ParseURL(url)

	if err != nil {
		return err
	}

	return nil
}

// This function writes new feed data to the database
func WriteFeedData(feedT *gofeed.Item, feedID int) (string, error) {
	DB := db.GetDB()
	feedData := &models.FeedData{}
	feedData.Link = feedT.Link
	feedData.Published = false
	feedData.Title = feedT.Title
	feedData.FeedID = feedID
	parsePubDate, _ := time.Parse(time.RFC1123Z, feedT.Published)
	feedData.PublishedDate = parsePubDate

	stmt, err := DB.Prepare("INSERT INTO feedData(feedid, title, link, publishedDate, published) VALUES(?,?,?,?,?)")
	defer stmt.Close()

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while preparing the query!")
		return "db_prepare_error", err
	}

	_, err = stmt.Exec(feedData.FeedID, feedData.Title, feedData.Link, feedData.PublishedDate, feedData.Published)

	if err != nil {
		log.WithFields(log.Fields{"error": err}).Error("There was an error while executing the query!")
		return "db_exec_error", err
	}

	return "db_feed_added", nil
}

// This function updates the Feed data and sets the respective feed data entry to published
func UpdateFeedDataPublished(data *models.FeedData) (bool, error) {
	DB := db.GetDB()
	_, err := DB.Exec("UPDATE feedData SET published=1 WHERE id = ?", data.ID)
	if err != nil {
		log.WithFields(log.Fields{"feedData": data}).Error("There was an error while updating the record of the feedData!")
		return false, err
	}
	return true, nil
}

// This function checks for the URL if it is contained in the feedData table
// We use this function to not write the same feed data into the table twice
func LinkExists(link string) bool {
	DB := db.GetDB()
	log.WithFields(log.Fields{"link": link}).Debug("Checking if the current feed posts URL exists in the database")
	err := DB.QueryRow("SELECT link FROM feedData WHERE link = ?", link).Scan(&link)
	if err != nil {
		if err != sql.ErrNoRows {
			log.WithFields(log.Fields{"link": link}).Debug("The URL was not found!")
			log.WithFields(log.Fields{"error": err}).Error("There was a problem while querying the database!")
			return false
		}
		return false
	}
	log.WithFields(log.Fields{"link": link}).Debug("The URL was found!")
	return true
}
