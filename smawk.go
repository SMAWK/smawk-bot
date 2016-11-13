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

func (bot *SmawkBot) ParseAndExecuteUpdate(update tgbotapi.Update) {
    if update.Message.Text != "" {
        cmd := strings.Split(update.Message.Text, " ")
        if (cmd[0] == "/start" || cmd[0] == "/start@smawk_bot") {
            bot.ExecuteStartCommand(update)
        } else if (cmd[0] == "/id" || cmd[0] == "/id@smawk_bot") {
            bot.ExecuteIDCommand(update)
        } else if (cmd[0] == "/hype" || cmd[0] == "/hype@smawk_bot" || strings.Contains(update.Message.Text, "/hype") || strings.Contains(update.Message.Text, "/hype@smawk_bot")) {
            bot.ExecuteHypeCommand(update)
        } else if (cmd[0] == "/whatchu_did_there" || cmd[0] == "/whatchu_did_there@smawk_bot") {
            bot.ExecuteWhatchuDidThereCommand(update)
        } else if (cmd[0] == "/score" || cmd[0] == "/score@smawk_bot") {
            bot.ExecuteScoreCommand(update, cmd)
        } else if (cmd[0] == "/upvote" || cmd[0] == "/upvote@smawk_bot") {
            bot.ExecuteUpvoteCommand(update, cmd)
        } else if (cmd[0] == "/downvote" || cmd[0] == "/downvote@smawk_bot") {
            bot.ExecuteDownvoteCommand(update, cmd)
        } else if (cmd[0] == "/bless" || cmd[0] == "/bless@smawk_bot") {
            bot.ExecuteBlessCommand(update, cmd)
        } else if (cmd[0] == "/curse" || cmd[0] == "/curse@smawk_bot") {
            bot.ExecuteCurseCommand(update, cmd)
        } else if (cmd[0] == "/splash" || cmd[0] == "/splash@smawk_bot") {
            bot.ExecuteSplashCommand(update)
        } else if (cmd[0] == "/why" || cmd[0] == "/why@smawk_bot") {
            bot.ExecuteWhyCommand(update)
        } else if (cmd[0] == "/smawk" || cmd[0] == "/smawk@smawk_bot") {
            bot.ExecuteSMAWKCommand(update, cmd)
        } else if (cmd[0] == "/me" || cmd[0] == "/me@smawk_bot") {
            bot.ExecuteSMAWKCommand(update, cmd)
        } else if (cmd[0] == "/dapun" || cmd[0] == "/dapun@smawk_bot") {
            bot.ExecuteDaPunCommand(update)
        }
    }
}
