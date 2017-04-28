package smawk

import (
	"database/sql"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strconv"
	"strings"
)

// ExecuteBlessCommand blesses a user for 3 points
func (bot *SmawkBot) ExecuteBlessCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {
	msg_string := "No blessing. All are equal in SMÄWK's eyes."
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
	return bot.API.Send(msg)
}

// ExecuteCurseCommand curses a user for 3 points
func (bot *SmawkBot) ExecuteCurseCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {
	msg_string := "No cursing. SMÄWK does not condone cursing."
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
	return bot.API.Send(msg)
}

// ExecuteScoreCommand returns the current point count for each user in the chat
func (bot *SmawkBot) ExecuteScoreCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {
	// Connect to our database
	db, err := ConnectDB(bot.dbPass)
	if err != nil {
		log.Fatal(err)
		return tgbotapi.Message{}, nil
	}
	defer db.Close()

	if len(cmd) == 1 {
		// Create our query
		users, err := db.Query("SELECT u.username, SUM(s.point) as `points` FROM scores s JOIN users u on u.id = s.user_id WHERE s.chat_id = "+strconv.FormatInt(update.Message.Chat.ID,10)+" GROUP BY s.user_id ORDER BY `points` DESC")
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}
		defer users.Close()

		// Get our scores
		msg_string := ""

		for users.Next() {
			var username string
			var points string
			if err := users.Scan(&username, &points); err != nil {
				log.Fatal(err)
			}
			msg_string += "\n"+username[1:]+": "+points
		}
		if err := users.Err(); err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)
	} else if len(cmd) == 2 {
		var total_points sql.NullString
		err = db.QueryRow("SELECT SUM(s.point) FROM scores s JOIN users u ON s.user_id = u.id WHERE u.username=? AND s.chat_id=?", cmd[1],strconv.FormatInt(update.Message.Chat.ID,10)).Scan(&total_points)
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		} else if err == sql.ErrNoRows || !total_points.Valid {
			msg_string := "User "+cmd[1]+"  does not exist.\n"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
			return bot.API.Send(msg)
		}

		msg_string := cmd[1]+" has "+total_points.String+" points, of which:\n"

		users, err := db.Query("SELECT SUM(s.point) as points, s.reason FROM scores s JOIN users u ON s.user_id = u.id WHERE s.chat_id = "+strconv.FormatInt(update.Message.Chat.ID,10)+" AND u.username = '"+cmd[1]+"' GROUP BY s.reason")
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}
		defer users.Close()
		for users.Next() {
			var points string
			var reason string
			if err := users.Scan(&points, &reason); err != nil {
				log.Fatal(err)
				return tgbotapi.Message{}, nil
			}
			msg_string += points + " is for " + reason + "\n"
		}
		if err := users.Err(); err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)
	}

	return tgbotapi.Message{}, nil
}

// ExecuteUpvoteCommand is responsibe for adding a point to a user
func (bot *SmawkBot) ExecuteUpvoteCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {

	if len(cmd) == 1 {
		// Wrong Usage
		msg_string := "Correct Usage: /upvote @username [reason]"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)

	} else if cmd[1] == "@"+update.Message.From.UserName {
		// Someone commited the cardinal sin
		err := bot.EnterScore(update.Message.Chat.ID, cmd[1], "Self Upvote", "-3")
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}
		msg_string := cmd[1]+" has been docked 3 points for self upvoting!"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)

	} else if len(cmd) == 2 {
		// Upvote User
		err := bot.EnterScore(update.Message.Chat.ID, cmd[1], "no reason", "1")
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}

		msg_string := cmd[1]+" has been upvoted 1 point"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)

	} else if len(cmd) >= 3 {
		// Create our reason
		var reason string
		if cmd[2] == "for" && len(cmd) > 3 {
			reason = strings.Join(cmd[3:]," ")
		} else {
			reason = strings.Join(cmd[2:]," ")
		}

		// Upvote User Reason
		err := bot.EnterScore(update.Message.Chat.ID, cmd[1], reason, "1")
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}

		msg_string := cmd[1]+" has been upvoted 1 point for "+reason
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)
	}

	return tgbotapi.Message{}, nil
}

// ExecuteDownvoteCommand docs points from a user
func (bot *SmawkBot) ExecuteDownvoteCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {

	if len(cmd) == 1 {
		// Wrong Usage
		msg_string := "Correct Usage: /downvote @username [reason]"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)

	} else if len(cmd) == 2 {
		// Downvote User
		err := bot.EnterScore(update.Message.Chat.ID, cmd[1], "no reason", "-1")
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}

		msg_string := cmd[1]+" has been downvoted 1 point"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)

	} else if len(cmd) >= 3 {
		// Create our reason
		var reason string
		if cmd[2] == "for" && len(cmd) > 3 {
			reason = strings.Join(cmd[3:]," ")
		} else {
			reason = strings.Join(cmd[2:]," ")
		}

		// Downvote User Reason
		err := bot.EnterScore(update.Message.Chat.ID, cmd[1], reason, "-1")
		if err != nil {
			log.Fatal(err)
			return tgbotapi.Message{}, nil
		}

		msg_string := cmd[1]+" has been downvoted 1 point for "+reason
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		return bot.API.Send(msg)
	}

	return tgbotapi.Message{}, nil
}
