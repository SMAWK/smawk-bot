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
	users, err := db.Query("SELECT username FROM users WHERE chat_id=? AND (flag_muted IS NULL OR flag_muted = '0') ORDER BY `username` ASC",update.Message.Chat.ID)
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

// ExecuteMuteCommand will mute the current user, so they are not targeted by the /all command
func (bot *SmawkBot) ExecuteMuteCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {
	db, err := ConnectDB(bot.dbPass)
	if err != nil {
		log.Fatal(err)
		return tgbotapi.Message{}, nil
	}
	defer db.Close()

	if len(cmd) > 1 {
		msg_string := "Correct Usage: /mute"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)
	} else {
		_, err := db.Query("UPDATE users SET flag_muted=1 WHERE chat_id=? AND username=? ",update.Message.Chat.ID,"@"+update.Message.From.UserName)
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}

		msg_string := "@"+update.Message.From.UserName+" has been muted."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)
	}

	return tgbotapi.Message{}, nil
}

// ExecuteUnmuteCommand will unmute the current user, so they are notified by the /all command
func (bot *SmawkBot) ExecuteUnmuteCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {
	db, err := ConnectDB(bot.dbPass)
	if err != nil {
		log.Fatal(err)
		return tgbotapi.Message{}, nil
	}
	defer db.Close()

	if len(cmd) > 1 {
		msg_string := "Correct Usage: /unmute"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)
	} else {
		_, err := db.Query("UPDATE users SET flag_muted=0 WHERE chat_id=? AND username=? ",update.Message.Chat.ID,"@"+update.Message.From.UserName)
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}

		msg_string := "@"+update.Message.From.UserName+" has been unmuted."
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)
	}

	return tgbotapi.Message{}, nil
}
