package smawk

import (
	"bytes"
	"database/sql"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// ExecuteStartCommand is launched when the bot is started, and sends a message to the chat that started it
func (bot *SmawkBot) ExecuteStartCommand(update tgbotapi.Update) (tgbotapi.Message, error) {
	// Create our message and send it to the chat that started the bot.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Lo, the official SMÄWKBot rises!")

	// Send the message
	return bot.API.Send(msg)
}

// ExecuteIDCommand returns the ID of a chat to the person that called this command.
// This command is only available on a private chat
func (bot *SmawkBot) ExecuteIDCommand(update tgbotapi.Update) (tgbotapi.Message, error) {
	// Check to see if the chat is a private chat
	if update.Message.Chat.Type == "private" {
		// Generate a message that contains the ID of our chat
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your chat ID is: "+strconv.FormatInt(update.Message.Chat.ID,10))

		// Send the message
		return bot.API.Send(msg)
	}

	return tgbotapi.Message{}, nil
}

// ExecuteSMAWKCommand is the command that is used to have 'third person' conversations inside
// of the SMAWK group chat. It is reserved for members of SMAWK only
func (bot *SmawkBot) ExecuteSMAWKCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {
	// Check to see if the user that called this command is actually a member of SMAWK
	// If so, go ahead and send the message
	if bot.isUser(update.Message.From.UserName) {
		// Look to see if the command was executed successfully
		switch len(cmd) {
			// Wrong Usage
			case 1:
				// Create our message with the instructions
				msg_string := "Correct Usage: /smawk <phrase>"
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)

				// Send the message
				return bot.API.Send(msg)

			// Correct Usage (len(cmd) >= 2)
			default:
				// Look to see if we are in test mode. If so we don't want to send the command to actual smawk
				if bot.Testing {
					phrase := strings.Join(cmd[1:]," ")
					msg_string := update.Message.From.UserName+" "+phrase

					// Create our message and prepare to send it to the SMAWK chat
					msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)

					// Send off our message
					return bot.API.Send(msg)
				}

				// Take what we said and turn it into a phrase
				phrase := strings.Join(cmd[1:]," ")
				msg_string := update.Message.From.UserName+" "+phrase

				// Create our message and prepare to send it to the SMAWK chat
				msg := tgbotapi.NewMessage(-9125034, msg_string)

				// Send off our message
				return bot.API.Send(msg)
		}
	}

	return tgbotapi.Message{}, nil
}



// To Do Below This
// ====================

func (bot *SmawkBot) ExecuteHypeCommand(update tgbotapi.Update) {
	// Make sure that we have the hype command in our working directory
	if _, err := os.Stat("hype.gif"); os.IsNotExist(err) {
		// NOOOO!!!! WE DON'T HAVE THE GIF!!!!!
		// Fetch it from the SMAWK source
		cmdname := "curl"
		cmdargs := []string{"-O","http://www.benjaminrmatthews.com/img/smawk-bot/hype.gif"}

		cmd := exec.Command(cmdname,cmdargs...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		}
	}

	doc := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, "hype.gif")
	bot.API.Send(doc)
}

func (bot *SmawkBot) ExecuteScoreCommand(update tgbotapi.Update, cmd []string) {
	// Connect to our database
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if len(cmd) == 1 {
		// Create our query
		users, err := db.Query("SELECT u.username, SUM(s.point) as `points` FROM scores s JOIN users u on u.id = s.user_id WHERE s.chat_id = "+strconv.FormatInt(update.Message.Chat.ID,10)+" GROUP BY s.user_id ORDER BY `points` DESC")
		if err != nil {
			log.Fatal(err)
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
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	} else if len(cmd) == 2 {
		var total_points sql.NullString
		err = db.QueryRow("SELECT SUM(s.point) FROM scores s JOIN users u ON s.user_id = u.id WHERE u.username=? AND s.chat_id=?", cmd[1],strconv.FormatInt(update.Message.Chat.ID,10)).Scan(&total_points)
		if err != nil {
				log.Fatal(err)
		} else if err == sql.ErrNoRows || !total_points.Valid {
			msg_string := "User "+cmd[1]+"  does not exist.\n"
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
			bot.API.Send(msg)
			return
		}

		msg_string := cmd[1]+" has "+total_points.String+" points, of which:\n"

		users, err := db.Query("SELECT SUM(s.point) as points, s.reason FROM scores s JOIN users u ON s.user_id = u.id WHERE s.chat_id = "+strconv.FormatInt(update.Message.Chat.ID,10)+" AND u.username = '"+cmd[1]+"' GROUP BY s.reason")
		if err != nil {
				log.Fatal(err)
		}
		defer users.Close()
		for users.Next() {
				var points string
				var reason string
				if err := users.Scan(&points, &reason); err != nil {
						log.Fatal(err)
				}
				msg_string += points+" is for "+reason+"\n"
		}
		if err := users.Err(); err != nil {
				log.Fatal(err)
		}

		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	}
}

func (bot *SmawkBot) ExecuteUpvoteCommand(update tgbotapi.Update, cmd []string) {
	// Connect to our database
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if len(cmd) == 1 {
		// Wrong Usage
		msg_string := "Correct Usage: /upvote @username [reason]"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	} else if cmd[1] == "@"+update.Message.From.UserName {
		// Someone commited the cardinal sin
		_, err := db.Query("INSERT INTO scores(user_id,point,chat_id,reason) SELECT id,-3,?,'Self Upvote' FROM users u WHERE u.username=?",update.Message.Chat.ID,cmd[1])
		if err != nil {
			log.Fatal(err)
		}
		msg_string := cmd[1]+" has been docked 3 points for self upvoting!"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	} else if len(cmd) == 2 {
		// Upvote User
		votes, err := db.Query("INSERT INTO scores(user_id,point,chat_id,reason) SELECT id,1,?,'no reason' FROM users u WHERE u.username=?",update.Message.Chat.ID,cmd[1])
		if err != nil {
				log.Fatal(err)
		}
		defer votes.Close()

		msg_string := cmd[1]+" has been upvoted 1 point"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	} else if len(cmd) >= 3 {
		// Create our reason
		var reason string
		if cmd[2] == "for" && len(cmd) > 3 {
			reason = strings.Join(cmd[3:]," ")
		} else {
			reason = strings.Join(cmd[2:]," ")
		}

		// Upvote User Reason
		votes, err := db.Query("INSERT INTO scores(user_id,point,chat_id,reason) SELECT id,1,?,? FROM users u WHERE u.username=?",update.Message.Chat.ID,reason,cmd[1])
		if err != nil {
				log.Fatal(err)
		}
		defer votes.Close()

		msg_string := cmd[1]+" has been upvoted 1 point for "+reason
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	}
}

func (bot *SmawkBot) ExecuteDownvoteCommand(update tgbotapi.Update, cmd []string) {
	// Connect to our database
	db, err := ConnectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	if len(cmd) == 1 {
		// Wrong Usage
		msg_string := "Correct Usage: /downvote @username [reason]"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	} else if len(cmd) == 2 {
		// Downvote User
		votes, err := db.Query("INSERT INTO scores(user_id,point,chat_id,reason) SELECT id,-1,?,'no reason' FROM users u WHERE u.username=?",update.Message.Chat.ID,cmd[1])
		if err != nil {
				log.Fatal(err)
		}
		defer votes.Close()

		msg_string := cmd[1]+" has been downvoted 1 point"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	} else if len(cmd) >= 3 {
		// Create our reason
		var reason string
		if cmd[2] == "for" && len(cmd) > 3 {
			reason = strings.Join(cmd[3:]," ")
		} else {
			reason = strings.Join(cmd[2:]," ")
		}

		// Downvote User Reason
		votes, err := db.Query("INSERT INTO scores(user_id,point,chat_id,reason) SELECT id,-1,?,? FROM users u WHERE u.username=?",update.Message.Chat.ID,reason,cmd[1])
		if err != nil {
				log.Fatal(err)
		}
		defer votes.Close()

		msg_string := cmd[1]+" has been downvoted 1 point for "+reason
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	}
}

func (bot *SmawkBot) ExecuteBlessCommand(update tgbotapi.Update, cmd []string) {
	if update.Message.From.UserName == "bnmtthews" || update.Message.From.UserName == "ReverendRecker" || update.Message.From.UserName == "CMoneys" {
		// Connect to our database
		db, err := ConnectDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		votes, err := db.Query("INSERT INTO scores(user_id,point,chat_id,reason) SELECT id,3,?,'Blessings from Dude' FROM users u WHERE u.username=?",update.Message.Chat.ID,cmd[1])
		if err != nil {
				log.Fatal(err)
		}
		defer votes.Close()

		msg_string := cmd[1]+" has been blessed for 3 points"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	} else {
		msg_string := "The power of blessing has not been bestowed upon you"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	}
}

func (bot *SmawkBot) ExecuteCurseCommand(update tgbotapi.Update, cmd []string) {
	if update.Message.From.UserName == "bnmtthews" || update.Message.From.UserName == "ReverendRecker" || update.Message.From.UserName == "CMoneys" {
		// Connect to our database
		db, err := ConnectDB()
		if err != nil {
			log.Fatal(err)
		}
		defer db.Close()

		votes, err := db.Query("INSERT INTO scores(user_id,point,chat_id,reason) SELECT id,-3,?,'Curses from Dude' FROM users u WHERE u.username=?",update.Message.Chat.ID,cmd[1])
		if err != nil {
				log.Fatal(err)
		}
		defer votes.Close()

		msg_string := cmd[1]+" has been cursed for 3 points"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	} else {
		msg_string := "The power of cursing has not been bestowed upon you"
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
		bot.API.Send(msg)
	}
}
