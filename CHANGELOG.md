# Formatting
Please use the following format when you update the changelog:
```
### <version> - <release date (yyyy-mm-dd)>
- Item 1
- Item 2
- ...
```

# Changelog
### v2.1.0
- Adds `today` command

### v2.0.0-beta
- VERSION 2.0.0!!!! (a.k.a. let's do this again)
- Refactors Library Codebase
	- ParseAndExecuteUpdates() now uses a more idiomatic `switch` vs `if / else`
	- Removes some of the unused commands
		- `whachu_did_there`
		- `dapun`
		- `splash` (since PokeGO died, the joke died also)
		- `why`
		- `hype` no longer can be used in the middle of a sentence
		- `labels` was a little too redundant
- Bug smashes
	- `me` and `smawk` are now hard limited to users of SMÄWK
- Test file covers more of the code, and tests deeper
- Changelog now shows most recent update first
- README has been updated with an explanation of all the commands
- Switches to a versioning approach
- Actually uses tabs like we are supposed to
- Code is now commented
- Transfers repos to central location
- Removes DB password from codebase, and requires it to be entered on the command line
- Commands are now split into several files to help find each other
- `/bless` and `/curse` have been removed. We are all equal in SMÄWK's sight
- The bot can now be deployed in any group chat, and should work fine
- The bot won't run commands if the user isn't registered
- New Commands (Say What!?)
	- `all` - (yes technically not new, but it's back)
	- `deregister` - Deregisters a user from using the chat calls
	- `here` - (literally just `all`)
	- `mute` - Will mute your username from `here`
	- `register` - Register your username on a different chat land
	- `unmute` - Will unmute your username from `here`
	- `version` - (so we can see what we are running)

### v1.1.2 - 2016-10-12
- File refactor

### v1.1.1 - 2016-10-07
- Fixes scoring between chats
- Adds a lovely 3rd person view

### v1.1.0 - 2016-08-08
- Adds Scoring
- Adds Upvoting
- Adds Downvoting

### v1.0.2 - 2016-08-05
- Lets /hype be called anywhere in a string
	- Adds appropriate tests

### v1.0.1 - 2016-08-01
- Adds new 'Whatchu Did There' gif

### v1.0.0 - 2016-07-27
- Converts from standalone program to library
- Adds self-signed certificate handler
- Updates README to reflect changes
- Adds changelog
- Removes /hello command
	- Removes /hello related unit tests
- Adds /id command
	- Adds /id unit tests
- Adds /hype command
	- Adds autofetch of gif if needed
