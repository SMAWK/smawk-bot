package main

import (
    "gopkg.in/telegram-bot-api.v4"
    "log"
    "net/http"
)

func main() {
    bot, err := tgbotapi.NewBotAPI("249930361:AAHz1Gksb-eT0SQG47lDb7WbJxujr7kGCkU")
    if err != nil {
        log.Fatal(err)
    }

    bot.Debug = true

    log.Printf("Authorized on account %s", bot.Self.UserName)

    _, err = bot.SetWebhook(tgbotapi.NewWebhook("https://mysimplethings.xyz:8443/" + bot.Token))
    if err != nil {
        log.Fatal(err)
    }

    updates := bot.ListenForWebhook("/" + bot.Token)
    go http.ListenAndServeTLS("0.0.0.0:8443", "cert.pem", "key.pem", nil)

    for update := range updates {
        log.Printf("%+v\n", update)
    }
}
