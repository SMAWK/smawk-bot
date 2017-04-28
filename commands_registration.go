package smawk

import (
	"gopkg.in/telegram-bot-api.v4"
	"log"
)

func (bot *SmawkBot) ExecuteRegisterCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {
	db, err := ConnectDB(bot.dbPass)
    if err != nil {
        log.Fatal(err)
        return tgbotapi.Message{}, nil
    }
    defer db.Close()

    if bot.isUser(update.Message.From.UserName, update.Message.Chat.ID) {
    	msg_string := "@"+update.Message.From.UserName+" already registered."
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
        return bot.API.Send(msg)
    }

    if len(cmd) > 1 {
        msg_string := "Correct Usage: /register"
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
        return bot.API.Send(msg)
    } else {
        _, err := db.Query("INSERT INTO users(username,chat_id,flag_muted) VALUES(?,?,0)","@"+update.Message.From.UserName,update.Message.Chat.ID)
        if err != nil {
            log.Fatal(err)
            return tgbotapi.Message{}, nil
        }

        msg_string := "@"+update.Message.From.UserName+" has been registered."
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
        return bot.API.Send(msg)
    }

    return tgbotapi.Message{}, nil
}

func (bot *SmawkBot) ExecuteDeregisterCommand(update tgbotapi.Update, cmd []string) (tgbotapi.Message, error) {
	db, err := ConnectDB(bot.dbPass)
    if err != nil {
        log.Fatal(err)
        return tgbotapi.Message{}, nil
    }
    defer db.Close()

    if len(cmd) > 1 {
        msg_string := "Correct Usage: /deregister"
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
        return bot.API.Send(msg)
    } else {
        _, err := db.Query("DELETE FROM users WHERE username=? AND chat_id=?","@"+update.Message.From.UserName,update.Message.Chat.ID)
        if err != nil {
            log.Fatal(err)
            return tgbotapi.Message{}, nil
        }

        msg_string := "@"+update.Message.From.UserName+" has been registered."
        msg := tgbotapi.NewMessage(update.Message.Chat.ID, msg_string)
        return bot.API.Send(msg)
    }

    return tgbotapi.Message{}, nil
}

