# SMÄWKBot (a.k.a smawk_bot) v1.5.0
This bot was written for a telegram group chat where several nerdy guys hang out. It was built (and is primary maintained) by [bmatt468](https://github.com/orgs/SMAWK/people/bmatt468) as a learning project / sandbox for Google Go.

# Getting Started
SMÄWKBot is written in Golang / Go. Therefore, the first step to working on this bot is [setting up Go](https://golang.org/doc/install).

If you are new to Go, I would hightly recommend [taking the tour](tour.golang.org), or reading through [GoByExample](gobyexample.com). These resources will help provide some light about what's going on (and maybe attract you to the language).

Once you have taken the steps to understand the nature of Go, you are ready to jump in and set up a bot of your own!

#Setting up the bot
###Controllers
This bot is a [library](https://golang.org/doc/code.html#Library), meaning that it can not run on its own. Rather, it offers methods to be used by an external controller. Most controllers will look nearly 100% the same, there isn't much difference between them. This bot is currently being run by [this controller](https://github.com/bmatt468/smawk-bot). I would recommend visiting that repo and looking through the README to understand the more detailed inner workings of a bot controller. Essentially the controller is responsible for:

- Controlling the access token
- Opening the webhook
- Listening on the webhook / fetching updates
- Calling this bots execute method for the updates

###Interfacing with the Bot
To interface with this API, add the following to the top of your controller:
```Go
import (
    "github.com/SMAWK/smawk-bot"
)
```

After importing, you can access each of the SMAWK/smawk-bot methods from your controller.

#Launching the bot
There are two ways to lauch the bot. One is launching inisde of the controller source, the other is installing a compiled binary and launching there.

###Controller Source (quick and easy)
To launch from the controller source, navigate to the working directory of your controller and run
```Shell
go build
./<compiled_binary>
```

###Installed Binary (slightly more complex but recommended)
Go has the ability to install a binary to `$GOPATH/bin`. During your setup of Go, you should have added this directory to your `PATH`, thus allowing the controller to be launched directly from the command line. To install the binary you have two options:

####Go Get (recommended)
These commands will take care of fetching your controller (and any updates to it), any dependencies (and updates), and compilation of the program. It will then move a binary to `$GOPATH/bin`.
```Shell
cd $GOPATH
go get -u github.com/<user_name>/<controller_repo>
<binary_name> #Generally this is just <controller_repo>
```

####Go Install
**NOTE:** to use this method you must have instantiated `$GOBIN`. Check the Golang docs for more info on this variable.
You can check if it's set by running `go env | grep GOBIN`
```Shell
cd $GOPATH/src/<user_name>/<controller_repo>
go install
<binary_name> #Generally this is just <controller_repo>
```

#Commands
###/start
`/start` is used to start the bot in a chat (both group and personal). It is a command that is required by Telegram; when run, it will return a message to the chat if it started successfully.

###/id
`id` can only be used in a personal chat. It will return your unique chat ID; this ID is helpful if you want to do work on your own instance of the bot.

###/hype
`/hype` returns the most-hyped-up gif of 2016. You just have to see it to believe it.

###/score
`/score` will display the current score of everyone in the chat (ordered by points). If you provide the optional username parameter, it will output the specific reasons for that users score (again, ordered by points)

Usage: `/score [<user>]`


###/upvote
`/upvote` is the command to give a user points. A username is required for the command to work, and you are allowed to add an optional reason (should you desire)

Usage: `/upvote <user> [<reason>]`


###/downvote
`/downvote` is the command to take points away from a user. You must provide a username for a person, but the reason for your downvote is optional.

Usage: `/downvote <user> [<reason>]`


###/bless
`/bless` is a special command that is reserved for the 'clergy' (a.k.a. the people ordained by the church of dude). It grants 3 points to the user of choice.

Usage: `/bless <user>`

###/curse
`/curse` is another special command that is reserved for the 'clergy'. It takes 3 points from a user of choice.

Usage: `/curse <user>`

###/smawk or /me
`/smawk` or `/me` is a special command that can only be executed from a private chat (and only by a member of SMÄWK). When run, it will print a message in the group chat in third-person.

Usage: `/smawk <statement>`

Example: `/smawk sighs` would output `<username> sighs` in the group chat.

# Adding Features
Any command will have an assigned controller function attached to it. These functions should be placed together inside of `smawk.go`. Please look at how the other commands are called, and replicate that style.

If you wish to add a new feature, make a new branch and name it with the following sheme: `feature/<your_feature_name>`. You are welcome to add any feature you wish, as long as you follow the rules.

###Feature Rules
1) **NEVER EVER EVER EVER EVER PUSH DIRECTLY TO MASTER.** When you are done with the feature (and have tested it sufficiently), submit a pull request; I will merge everything in, and make sure that it builds in production.

2) If you add a new feature, please write sufficient unit tests for it (and place them inside `smawk_test.go`). See the [Testing](#testing) section for a little more info about writing unit tests.

3) When you add a feature, please update the CHANGELOG (and respect formatting please)

# Reporting / Fixing Bugs
Seeing as we are all human (and not Dr. K), bugs are bound to happen. ~~If~~ When that happens, simply open an issue and try to explain in detail what's going on. If you are feeling really adventurous, you can create a new branch title `bugs/<issue#>` and try to patch it yourself. Otherwise, I'll work to get it resolved.

# Testing
Testing the bot is as easy as calling `go test` from the command line of the `smawk-bot` working directory. This command will execute all the tests located inside of `smawk-test.go`. The results of the test will be shown on the command line.

Note: The results of the test command will be sent to the chat id specified by the `ChatID` const inside `smawk-test.go` (Approx. line 14). Currently, this ID is set to be the personal chat between [@bmatt468](https://github.com/orgs/SMAWK/people/bmatt468) and SMÄWKBot. When you test, please change this ID to match the ID of your personal chat with the bot. Your chat id can be obtained by calling starting a private chat with the bot (use the search bar or the group chat) and then typing `/id` into the chat. The bot will respond with your personal chat ID. Please note that this command will not work in a group chat.

For more info about testing, see the official Golang docs [here](https://golang.org/pkg/testing/).

#Changelog
Changelog can be viewed [here](CHANGELOG.md)
