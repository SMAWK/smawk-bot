# SMÄWKBot (a.k.a smawk_bot)
###A Telegram Bot for the Almighty ~~She Mans Woman Haters Club~~ SMÄWK

# Getting Started
SMÄWKBot is written in Golang / Go. Therefore, the first step to working on this bot is [setting up Go](https://golang.org/doc/install).

If you are new to Go, I would hightly recommend [taking the tour](tour.golang.org), or reading through [GoByExample](gobyexample.com). These resources will help provide some light about what's going on (and maybe attract you to the language).

####Setting up local environment
Once you are ready to start developing with the bot and have [set up your workspace](https://golang.org/doc/code.html#Organization), clone this repo using the following command (assuming you have SSH set up):
```bash
git clone git@github.com:SMAWK/smawk-bot.git $GOPATH/src/github.com/SMAWK/smawk-bot
```

Once the repo is cloned in, run `make` to download all the dependencies required to make the bot run.

# Adding Features
If you wish to add a new feature, make a new branch and name it with the following sheme: `feature/<your_feature_name>`. You are welcome to add any feature you wish, as long as you follow the rules.

####Feature Rules
1) If you add a third-party library, make sure you update the buildscript (i.e., Makefile) to make it easier for everyone else to get all the needed dependencies. The script on the server calls the `make` command before it builds the project; by updating the files we make sure that everything stays consistent.

2) **NEVER EVER EVER EVER EVER PUSH DIRECTLY TO MASTER.** When you are done with the feature (and have tested it sufficiently), submit a pull request; I will merge everything in, and make sure that it builds in production.

3) Please do not attempt to use this code to create your own bot (i.e., don't try to run this program directly). You are more than welcome to use the logic in your own repositories, but please do not try to execute this program directly; it will cause conflicts on Telegrams end with two bots sharing websockets. See [Future Happenings](#future-happenings) for more info about this.

4) If you add a new feature, please write sufficient unit tests for it (and place them inside `smawk_test.go`). See the [Testing](#testing) section for a little more info about writing unit tests.

# Reporting / Fixing Bugs
Seeing as we are all human (and not Dr. K), bugs are bound to happen. ~~If~~ When that happens, simply open an issue and try to explain in detail what's going on. If you are feeling really adventurous, you can create a new branch title `bugs/<issue#>` and try to patch it yourself. Otherwise, I'll work to get it resolved.

# Testing
Testing the bot is as easy as calling `go test` from the command line of the `smawk-bot` working directory. This command will execute all the tests located inside of `smawk-test.go`. The results of the test will be shown on the command line.

Note: The results of the test command will be sent to the chat id specified by the `ChatID` const inside `smawk-test.go` (Approx. line 14). Currently, this ID is set to be the personal chat between [@bmatt468](https://github.com/orgs/SMAWK/people/bmatt468) and SMÄWKBot. When you test, please change this ID to match the ID of your personal chat with the bot. Your chat id can be obtained by calling starting a private chat with the bot (use the search bar or the group chat) and then typing `/id` into the chat. The bot will respond with your personal chat ID. Please note that this command will not work in a group chat.

For more info about testing, see the official Golang docs [here](https://golang.org/pkg/testing/).

# Future Happenings
Currently, SMÄWKBot acts as a Go program (i.e., it's in the `package main`). In the near future, this repo will be converted into a library that can be plugged into other packages; this change will allow everyone to launch their own bots (if they wish)
