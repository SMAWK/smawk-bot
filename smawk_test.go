package smawk_test

import (
	//"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	//"os"
	//"os/exec"
	"strconv"
	//"strings"
	"testing"
	"time"
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

func connect() (*sql.DB, error) {
	cfg := &mysql.Config {
		User: "smawk-bot",
		Passwd: "SM@WKisGR8",
		Net: "tcp",
		Addr: "107.170.45.12:3306",
		DBName: "smawk-bot",
	}
	return sql.Open("mysql", cfg.FormatDSN())
}

func timestamp() string {
	t := time.Now()
	return t.Format("2006-01-02 15:04:05 ");
}

/* ================================================ */
/*                Testing functions                 */
/* ================================================ */
func TestHelpers(t *testing.T) {
	fmt.Println("======= Starting Helper Tests =======")
	/** === Loading Bot === **/
	fmt.Print(timestamp()+"Loading SMÄWK_bot.... ")

	// Fetch our bot using the helper function
	_, err := getBot(t)

	// Check to see if something bad happened and break if need be
	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	// Otherwise, log to the console that we authenticated properly
	fmt.Println("done")

	/** === Database connection === **/
	fmt.Print(timestamp()+"Connecting to database.... ")
	db, err := connect()
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

    fmt.Println("done")
    fmt.Println("======= Helper Tests Succeeded =======\n")
}

func TestCommands(t *testing.T) {
	fmt.Println("======= Starting Command Tests =======")
	/*
	start
id
hype
score
upvote
downvote
bless
curse
smawk/me
*/
    fmt.Println("======= Command Tests Succeeded =======")
}
