package smawk

import (
	"database/sql"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
)

// ExecuteLabelCommand assigns a label to the specified user, in the specified channel
func (bot *SmawkBot) ExecuteLabelCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {
	// Connect to our database
	db, err := ConnectDB(bot.dbPass)
	if err != nil {
		log.Fatal(err)
		return tgbotapi.Message{}, nil
	}
	defer db.Close()

	if len(cmd) == 1 || len(cmd) == 2 {
		// Wrong Usage
		msg_string := "Correct Usage: /label @username <name>"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)

	} else if cmd[1] == "@"+update.Message.From.UserName {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "One must not label themself.")
		return bot.API.Send(msg)

	} else if len(cmd) >= 3 {
		label := strings.Join(cmd[2:]," ")

		votes, err := db.Query("UPDATE users SET label=? WHERE username=? ",label,cmd[1])
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}
		defer votes.Close()

		msg_string := cmd[1]+" is now "+label
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)
	}

	return tgbotapi.Message{}, nil
}

// ExecuteWhoisCommand is used to get the label for a specified user
func (bot *SmawkBot) ExecuteWhoisCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {
	// Connect to our database
	db, err := ConnectDB(bot.dbPass)
	if err != nil {
		log.Fatal(err)
		return tgbotapi.Message{}, nil
	}
	defer db.Close()

	if len(cmd) == 1 {
		// Wrong Usage
		msg_string := "Correct Usage: /whois @username"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)

	} else if len(cmd) >= 2 {
		var label sql.NullString
		err = db.QueryRow("SELECT label FROM users WHERE username=?", cmd[1]).Scan(&label)
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		} else if err == sql.ErrNoRows || !label.Valid {
			msg_string := cmd[1]+" has not been labeled.\n"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
			return bot.API.Send(msg)
		}

		msg_string := cmd[1]+" is known as "+label.String
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)
	}

	return tgbotapi.Message{}, nil
}
