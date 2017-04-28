package smawk

import (
	//"database/sql"
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

// ExecuteAllCommand is responsible for notifying each of the users in the channel about a message.
// It will fetch all the users from the database, and build a message string of their usernames
func (bot *SmawkBot) ExecuteAllCommand(update tgbotapi.Update) (tgbotapi.Message, error) {
	// Connect to our database
	db, err := ConnectDB(bot.dbPass)
	if err != nil {
		log.Fatal(err)
		return tgbotapi.Message{}, nil
	}
	defer db.Close()

	// Create our query
	users, err := db.Query("SELECT username FROM users")
	if err != nil {
		log.Fatal(err)
		return tgbotapi.Message{}, nil
	}
	defer users.Close()

	// Get our scores
	msg_string := ""

	for users.Next() {
		var username string
		if err := users.Scan(&username); err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}
		msg_string += " " + username
	}
	if err := users.Err(); err != nil {
		log.Fatal(err)
		return tgbotapi.Message{}, nil
	}

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
	return bot.API.Send(msg)
}
