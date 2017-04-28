package smawk

import (
	"gopkg.in/telegram-bot-api.v4"
	"strconv"
	"strings"
)

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
	if bot.isSmawkUser(update.Message.From.UserName) {
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

// ExecuteStartCommand is launched when the bot is started, and sends a message to the chat that started it
func (bot *SmawkBot) ExecuteStartCommand(update tgbotapi.Update) (tgbotapi.Message, error) {
	// Create our message and send it to the chat that started the bot.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Lo, the official SMÄWKBot rises!")

	// Send the message
	return bot.API.Send(msg)
}

// ExecuteVersionCommand returns the current version of the bot back to the chat that called it
func (bot *SmawkBot) ExecuteVersionCommand(update tgbotapi.Update) (tgbotapi.Message, error) {
	// Create our message and send it to the chat that started the bot.
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Current Bot Version: " + bot.Version)

	// Send the message
	return bot.API.Send(msg)
}
