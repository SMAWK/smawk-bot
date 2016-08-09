package smawk

import (
    "bytes"
    "database/sql"
    "encoding/json"
    "fmt"
    "github.com/go-sql-driver/mysql"
    "gopkg.in/telegram-bot-api.v4"
    "log"
    "os"
    "os/exec"
    "strconv"
    "strings"
    "time"
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
        }
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

func (bot *SmawkBot) ExecuteScoreCommand(update tgbotapi.Update, cmd []string) {
    // Connect to our database
    db, err := ConnectDB()
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if len(cmd) == 1 {
        // Create our query
        users, err := db.Query("SELECT u.username, SUM(s.point) as 'points' FROM scores s JOIN users u on u.id = s.user_id WHERE s.chat_id = "+strconv.FormatInt(update.Message.Chat.ID,10)+" GROUP BY s.user_id")
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
            msg_string += "\n"+username+": "+points
        }
        if err := users.Err(); err != nil {
                log.Fatal(err)
        }

        msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
        bot.API.Send(msg)
    } else if len(cmd) == 2 {
        var total_points sql.NullString
        err = db.QueryRow("SELECT SUM(s.point) FROM scores s JOIN users u ON s.user_id = u.id WHERE u.username=?", cmd[1]).Scan(&total_points)
        if err != nil {
                log.Fatal(err)
        } else if err == sql.ErrNoRows || !total_points.Valid {
            msg_string := "User "+cmd[1]+"  does not exist.\n"
            msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
            bot.API.Send(msg)
            return
        }

        msg_string := cmd[1]+" has "+total_points.String+" points, of which:\n"

        users, err := db.Query("SELECT SUM(s.point) as points, s.reason FROM scores s JOIN users u ON s.user_id = u.id WHERE u.username = '"+cmd[1]+"' GROUP BY s.reason")
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

func (bot *SmawkBot) ExecuteSplashCommand(update tgbotapi.Update) {
    msg_string := "@"+update.Message.From.UserName+" used splash....."
    msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
    bot.API.Send(msg)

    duration := time.Duration(5)*time.Second
    time.Sleep(duration)
    msg_string2 := "... but nothing happened"
    msg2 := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string2)
    bot.API.Send(msg2)
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

// ConnectDB takes care of opening a proper connection to the database to retrieve the scores that we need
func ConnectDB() (*sql.DB, error) {
    cfg := &mysql.Config {
        User: "smawk-bot",
        Passwd: "SM@WKisGR8",
        Net: "tcp",
        Addr: "107.170.45.12:3306",
        DBName: "smawk-bot",
    }
    return sql.Open("mysql", cfg.FormatDSN())
}
