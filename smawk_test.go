package smawk_test

import (
    "gopkg.in/telegram-bot-api.v4"
    "log"
	"testing"
)

const (
	SMAWKToken              = "249930361:AAHz1Gksb-eT0SQG47lDb7WbJxujr7kGCkU"
	ChatID                 	= 55997207
)

func getBot(t *testing.T) (*tgbotapi.BotAPI, error) {
	bot, err := tgbotapi.NewBotAPI(SMAWKToken)

	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	return bot, err
}

func TestLoadBot(t *testing.T) {
	_, err := getBot(t)

	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	log.Printf("SMÄWK_bot authenticated")
}

func TestSendMessage(t *testing.T) {
	bot, _ := getBot(t)

	msg := tgbotapi.NewMessage(ChatID, "SMÄWKBot:Go - Test message.")
	_, err := bot.Send(msg)

	if err != nil {
		log.Fatal(err)
		t.FailNow()
	}

	log.Printf("Test message sent")
}
