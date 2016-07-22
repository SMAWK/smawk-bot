package smawk

import (
    "gopkg.in/telegram-bot-api.v4"
    "log"
    "net/http"
)

// BotAPI allows you to interact with the Telegram Bot API.
type SmawkBot struct {
    API *tgbotapi.BotAPI
}

// Connect takes a provided access token, and returns a pointer
// to the Telegram Bot API. This function must be called in order
// to have access to any of the commands
func Connect(tkn string, debug bool) (*SmawkBot) {
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
    return sbot
}

// OpenWebhookWithCert is the wrapper function that calls Telegram's Bot API
// and listens to a self-signed https webhook for commands
func (bot *SmawkBot) OpenWebhookWithCert(url string, cert string) {
    _, err := bot.API.SetWebhook(tgbotapi.NewWebhookWithCert(url, cert))
    if err != nil {
        log.Fatal(err)
    }
}

func (bot *SmawkBot) Listen(token string) {
    // Start listening on our webhook for the commands
    // Spin off a goroutine to handle listening elsewhere
    //updates := bot.API.ListenForWebhook("/309LKj2390gklj1LJF2")
    go http.ListenAndServeTLS("0.0.0.0:8443", "smawk_cert.pem", "smawk_key.pem", nil)
}
