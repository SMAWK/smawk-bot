package smawk_test

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
	"testing"
)

// Create our constants for use throughout the testing functions
const (
	SMAWKToken              = "249930361:AAHz1Gksb-eT0SQG47lDb7WbJxujr7kGCkU"
	ChatID                 	= 55997207
)

/* ================================================ */
/*                 Helper functions                 */
/* ================================================ */

// getBot is a helper function that returns the bot object
// to each of the test functions
func getBot(t *testing.T) (*tgbotapi.BotAPI, error) {
	// Get the bot using the SMÄWK token
	bot, err := tgbotapi.NewBotAPI(SMAWKToken)

	// Check to see if something bad happened and break if need be
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	// Otherwise return our bot
	return bot, err
}

// generateUpdate is a helper function that generates a test update
// (see Update sruct in tgbotapi/types)
func generateUpdate(t *testing.T, cmd string) (tgbotapi.Update, error) {
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

/* ================================================ */
/*                Testing functions                 */
/* ================================================ */

// TestLoadBot tests to see if the bot is loading and authenticated properly
func iTestLoadBot(t *testing.T) {
	// Fetch our bot using the helper function
	_, err := getBot(t)

	// Check to see if something bad happened and break if need be
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	// Otherwise, log to the console that we authenticated properly
	log.Printf("SMÄWK_bot authenticated")
}

// TestSendMessage tests to see if the bot can properly send a message to the provided
// chatID (the private chat of the user)
func iTestSendMessage(t *testing.T) {
	// Fetch our bot using the helper function
	bot, _ := getBot(t)

	// Generate our message and send it to the private chat
	msg := tgbotapi.NewMessage(ChatID, "SMÄWKBot:Go - Test message.")
	_, err := bot.Send(msg)

	// Check to see if something bad happened and break if need be
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	// Otherwise, log to the console that the message was sent
	log.Printf("Test message sent")
}

// TestIDCommand tests to make sure that the bot will properly send an
// ID to a private chat while refusing to send to a public chat
func iTestIDCommand(t *testing.T) {
	// Fetch our bot using the helper function
	bot, _ := getBot(t)

	update, _ := generateUpdate(t,"/id")

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, "Your chat ID is: "+strconv.FormatInt(update.Message.Chat.ID,10))
    bot.Send(msg)
}

// TestContainedHypeCommand tests to make sure that the bot will properly handle a /hype
// command that is contained inside of a string
func iTestContainedHypeCommand(t *testing.T) {
	// Fetch our bot using the helper function
	bot, _ := getBot(t)

	update, _ := generateUpdate(t,"test /hype")

	if (strings.Contains(update.Message.Text, "/hype") || strings.Contains(update.Message.Text, "/hype@smawk_bot")) {
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
	    bot.Send(doc)
    }
}

// TestDatabaseConnection Makes sure that we are connected to our database where the scores are located
func TestDatabaseConnection(t *testing.T) {
	cfg := &mysql.Config {
		User: "smawk-bot",
		Passwd: "SM@WKisGR8",
		Net: "tcp",
		Addr: "107.170.45.12:3306",
	}
	db, err := sql.Open("mysql", cfg.FormatDSN())
    if err != nil {
      	log.Fatal(err)
		t.FailNow()
    }
    defer db.Close()

    err = db.Ping()
    if err != nil {
        log.Fatal(err)
		t.FailNow()
    }
}
