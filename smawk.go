package smawk

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strings"
)

// BotAPI allows you to interact with the Telegram Bot API.
type SmawkBot struct {
	API *tgbotapi.BotAPI
}

// Connect takes a provided access token, and returns a pointer
// to the Telegram Bot API. This function must be called in order
// to have access to any of the commands
func Connect(tkn string, debug bool) (*SmawkBot, error) {
	// Call the Telegram API wrapper and authenticate our Bot
	bot, err := tgbotapi.NewBotAPI(tkn)

	// Check to see if there were any errors with our bot and fail
	// if there were
	if err != nil {
		log.Fatal(err)
	}

	if (debug) {
		// Print confirmation
		log.Printf("Authorized on account %s", bot.Self.UserName)
	}

	// Set our bot to either be in debug mode (everything gets put out to the console)
	// or non debug mode (everything is silent)
	bot.Debug = debug

	// Create the SmawkBot
	sbot := &SmawkBot {
		API: bot,
	}

	// Return our bot back to the caller
	return sbot, err
}

// OpenWebhookWithCert is the wrapper function that calls Telegram's Bot API
// and listens to a self-signed https webhook for commands
func (bot *SmawkBot) OpenWebhookWithCert(url string, cert string) {
	_, err := bot.API.SetWebhook(tgbotapi.NewWebhookWithCert(url, cert))
	if err != nil {
		log.Fatal(err)
	}
}

// Listen opens a connection on the specified url and waits for a command
// to come in. After it receives a command from the API, it returns the update
// channel to the caller
func (bot *SmawkBot) Listen(token string) <-chan tgbotapi.Update {
	updates := bot.API.ListenForWebhook(token)
	return updates
}

// ParseAndExecuteUpdate takes in the Update struct from the API,
// and isolates the command and arguments, then passes the information
// on to the proper method
func (bot *SmawkBot) ParseAndExecuteUpdate(update tgbotapi.Update) {
	if update.Message.Text != "" {
		// Get the command and remove the trailing '@smawk_bot' (if it exists)
		switch cmd := strings.Split(update.Message.Text, " "); strings.Replace(cmd[0],"@smawk_bot","",-1) {
			case "/start":
				bot.ExecuteStartCommand(update)
			case "/id":
				bot.ExecuteIDCommand(update)
			case "/smawk", "/me":
				bot.ExecuteSMAWKCommand(update, cmd)
			case "/hype":
				bot.ExecuteHypeCommand(update)
			case "/score":
				bot.ExecuteScoreCommand(update, cmd)
			case "/upvote":
				bot.ExecuteUpvoteCommand(update, cmd)
			case "/downvote":
				bot.ExecuteDownvoteCommand(update, cmd)
			case "/bless":
				bot.ExecuteBlessCommand(update, cmd)
			case "/curse":
				bot.ExecuteCurseCommand(update, cmd)
		}
	}
}
