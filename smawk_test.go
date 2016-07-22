package smawk_test

import (
	"encoding/json"
	"gopkg.in/telegram-bot-api.v4"
	"log"
	"strconv"
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
func TestSendMessage(t *testing.T) {
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

// TestParseHello emulates the /hello{@smawk_bot} command and responds as it does in production
func TestParseHello(t *testing.T) {
	// Fetch our bot using the helper function
	bot, _ := getBot(t)

	// Generate our update using the helper function
	upd, _  := generateUpdate(t,"/hello")

	// Check our Message to ensure text is fine
	if (upd.Message.Text == "/hello") {
		// Generate our message to send back
		msg := tgbotapi.NewMessage(upd.Message.Chat.ID, "Hello, @" + upd.Message.From.UserName + "!")
		_, err := bot.Send(msg)

		// Check to see if something bad happened and break if need be
		if err != nil {
			log.Fatal(err)
			t.FailNow()
		}

		// Otherwise, log to the console that the message was sent
		log.Printf("/hello command test success")
	} else {
		log.Fatalf("generateUpdate fail: Expected /hello. Received %v", upd.Message.Text)
		t.FailNow()
	}
}
