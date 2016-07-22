package smawk

import (
    "gopkg.in/telegram-bot-api.v4"
    "log"
    "encoding/json"
    "strconv"
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
    cmd := update.Message.Text
    if (cmd == "/start" || cmd == "/start@smawk_bot") {
        bot.ExecuteStartCommand(update)
    } else if (cmd == "/hello" || cmd == "/hello@smawk_bot") {
        bot.ExecuteHelloCommand(update)
    }
}

/* ================================================ */
/*                Command functions                 */
/* ================================================ */

func (bot *SmawkBot) ExecuteStartCommand(update tgbotapi.Update) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Lo, the official SMÄWKBot rises!")
    bot.API.Send(msg)
}

func (bot *SmawkBot) ExecuteHelloCommand(update tgbotapi.Update) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Hello, @" + update.Message.From.UserName + "!")
    bot.API.Send(msg)
}

/* ================================================ */
/*                Testing functions                 */
/* ================================================ */

// GenerateUpdate is a helper function that generates a test update
// (see Update sruct in tgbotapi/types). This function can be called
// from the test files of programs that implement this library
func (bot *SmawkBot) GenerateUpdate(cmd string) (tgbotapi.Update, error) {
    // Create our Update Var
    var upd tgbotapi.Update

    // Create our JSON blob
    var updjson = []byte(`{
        "update_id":322176086,
        "message":{
            "message_id":178,
            "from":{
                "id":55997207,
                "first_name":"Benjamin",
                "last_name":"Matthews",
                "username":"bnmtthews"
            },
            "chat":{
                "id":55997207,
                "first_name":"Benjamin",
                "last_name":"Matthews",
                "username":"bnmtthews",
                "type":"private"
            },
            "date":1468013062,
            "text":"`+cmd+`",
            "entities":[{
                "type":"bot_command",
                "offset":0,
                "length":`+strconv.Itoa(len(cmd))+`
            }]
        }
    }`)

    // Create our update
    json.Unmarshal(updjson, &upd)

    // Return our update
    return upd, nil
}
