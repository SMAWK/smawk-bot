package smawk

import (
	"bytes"
	"fmt"
	"gopkg.in/telegram-bot-api.v4"
	"os"
	"os/exec"
)

// ExecuteHypeCommand send that amazing most hypeful gif to the smawk chat
func (bot *SmawkBot) ExecuteHypeCommand(update tgbotapi.Update) (tgbotapi.Message, error) {
	// Make sure that we have the hype command in our working directory
	if _, err := os.Stat("hype.gif"); os.IsNotExist(err) {
		// NOOOO!!!! WE DON'T HAVE THE GIF!!!!!
		// Fetch it from the SMAWK source
		cmdname := "curl"
		cmdargs := []string{"-O","http://www.benjaminrmatthews.com/img/smawk-bot/hype.gif"}

		cmd := exec.Command(cmdname,cmdargs...)
		var stderr bytes.Buffer
		cmd.Stderr = &stderr
		err := cmd.Run()
		if err != nil {
			fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
			return tgbotapi.Message{}, nil
		}
	}

	doc := tgbotapi.NewDocumentUpload(update.Message.Chat.ID, "hype.gif")
	return bot.API.Send(doc)
}
