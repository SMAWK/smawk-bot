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

/* ================================================ */
/*                Testing functions                 */
/* ================================================ */

// TestLoadBot tests to see if the bot is loading and authenticated properly
func TestLoadBot(t *testing.T) {
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

    log.Printf("DB Connection Successful")
}

// TestUsersFetch makes a test call to the database to get the list of users in the database
func TestUsersFetch(t *testing.T) {
	db, err := connect()
	if err != nil {
      	log.Fatal(err)
		t.FailNow()
    }
    defer db.Close()

    users, err := db.Query("SELECT id, username, first_name, last_name FROM users")
    if err != nil {
            log.Fatal(err)
    }
    defer users.Close()
    for users.Next() {
            var id string
            var username string
            var first_name string
            var last_name string
            if err := users.Scan(&id, &username, &first_name, &last_name); err != nil {
                    log.Fatal(err)
            }
            fmt.Printf("ID:%s Username:%s Name: %s %s\n", id, username, first_name, last_name)
    }
    if err := users.Err(); err != nil {
            log.Fatal(err)
    }

    fmt.Printf("\n=============\n")
}

// TestScoreCommand makes sure that we can connect to the database and properly obtain the score for everybody
func TestScoreCommand(t *testing.T) {
	db, err := connect()
	if err != nil {
      	log.Fatal(err)
		t.FailNow()
    }
    defer db.Close()

    users, err := db.Query("SELECT u.username, SUM(s.point) as 'points' FROM scores s JOIN users u on u.id = s.user_id WHERE s.chat_id = -9125034 GROUP BY s.user_id")
    if err != nil {
            log.Fatal(err)
    }
    defer users.Close()
    for users.Next() {
            var username string
            var points string
            if err := users.Scan(&username, &points); err != nil {
                    log.Fatal(err)
            }
            fmt.Printf("%s: %s\n", username, points)
    }
    if err := users.Err(); err != nil {
            log.Fatal(err)
    }
    fmt.Printf("\n=============\n")
}

func TestUserScoreCommand(t *testing.T) {
	db, err := connect()
	if err != nil {
      	log.Fatal(err)
		t.FailNow()
    }
    defer db.Close()

    update, _ := generateUpdate(t,"/score @bnmtthews")
    cmd := strings.Split(update.Message.Text, " ")

    var total_points sql.NullString
    err = db.QueryRow("SELECT SUM(s.point) FROM scores s JOIN users u ON s.user_id = u.id WHERE u.username=?", cmd[1]).Scan(&total_points)
    if err != nil {
            log.Fatal(err)
    } else if err == sql.ErrNoRows {
    	fmt.Printf("User %s does not exist.\n",cmd[1])
    	return
    } else if !total_points.Valid {
    	fmt.Printf("User %s does not exist.\n",cmd[1])
    	return
    }

    fmt.Printf("%s has %s points, of which:\n",cmd[1],total_points.String)

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
            fmt.Printf("%s is for %s\n", points, reason)
    }
    if err := users.Err(); err != nil {
            log.Fatal(err)
    }
}
