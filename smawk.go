package smawk

import (
    "bytes"
    "encoding/json"
    "fmt"
    "gopkg.in/telegram-bot-api.v4"
    "log"
    "os"
    "os/exec"
    "strconv"
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
    cmd := strings.Split(update.Message.Text, " ")
    if (cmd[0] == "/start" || cmd[0] == "/start@smawk_bot") {
        bot.ExecuteStartCommand(update)
    } else if (cmd[0] == "/id" || cmd[0] == "/id@smawk_bot") {
        bot.ExecuteIDCommand(update)
    } else if (cmd[0] == "/hype" || cmd[0] == "/hype@smawk_bot" || strings.Contains(update.Message.Text, "/hype") || strings.Contains(update.Message.Text, "/hype@smawk_bot")) {
        bot.ExecuteHypeCommand(update)
    } else if (cmd[0] == "/whatchu_did_there" || cmd[0] == "/whatchu_did_there@smawk_bot") {
        bot.ExecuteWhatchuDidThereCommand(update)
    }
}

/* ================================================ */
/*                Command functions                 */
/* ================================================ */

func (bot *SmawkBot) ExecuteStartCommand(update tgbotapi.Update) {
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Lo, the official SMÄWKBot rises!")
    bot.API.Send(msg)
}

func (bot *SmawkBot) ExecuteIDCommand(update tgbotapi.Update) {
    if update.Message.Chat.Type == "private" {
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your chat ID is: "+strconv.FormatInt(update.Message.Chat.ID,10))
        bot.API.Send(msg)
    }
}

func (bot *SmawkBot) ExecuteHypeCommand(update tgbotapi.Update) {
    // Make sure that we have the hype command in our working directory
    if _, err := os.Stat("hype.gif"); os.IsNotExist(err) {
        // NOOOO!!!! WE DON'T HAVE THE GIF!!!!!
        // Fetch it from the SMAWK source
        cmdname := "curl"
        cmdargs := []string{"-O","http://mysimplethings.xyz/img/smawk-bot/hype.gif"}

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

func (bot *SmawkBot) ExecuteWhatchuDidThereCommand(update tgbotapi.Update) {
    // Make sure that we have the hype command in our working directory
    if _, err := os.Stat("whoa.gif"); os.IsNotExist(err) {
        // NOOOO!!!! WE DON'T HAVE THE GIF!!!!!
        // Fetch it from the SMAWK source
        cmdname := "curl"
        cmdargs := []string{"-O","http://mysimplethings.xyz/img/smawk-bot/whoa.gif"}

        cmd := exec.Command(cmdname,cmdargs...)
        var stderr bytes.Buffer
        cmd.Stderr = &stderr
        err := cmd.Run()
        if err != nil {
            fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
        }
    }

    doc := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, "whoa.gif")
    bot.API.Send(doc)
}

/* ================================================ */
/*                 Helper functions                 */
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

// GenerateCertificate is used to create a self signed certificate for use with any
// instance of the bot
func GenerateCertificate(c string, st string, ct string, org string, dom string, key string, cert string) {
    // Generate our string for the certificate
    certstring := "/C="+c+"/ST="+st+"/L="+ct+"/O="+org+"/CN="+dom

    cmdname := "openssl"
    cmdargs := []string{"req","-newkey","rsa:2048","-sha256","-nodes","-keyout",key,"-x509","-days","365","-out",cert,"-subj",certstring}

    cmd := exec.Command(cmdname,cmdargs...)
    var stderr bytes.Buffer
    cmd.Stderr = &stderr
    err := cmd.Run()
    if err != nil {
        fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
    }
}
