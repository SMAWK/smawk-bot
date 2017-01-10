#Formatting
Please use the following format when you update the changelog:
```
###<version> - <release date (yyyy-mm-dd)>
- Item 1
- Item 2
- ...
```

#Changelog
###v1.2.0
- Refactors Library Codebase
	- ParseAndExecuteUpdates() now uses a more idiomatic `switch` vs `if / else`
	- Removes some of the unused commands
		- `whachu_did_there`
		- `dapun`
		- `splash` (since PokeGO died, the joke died also)
		- `why`
		- `hype` no longers can be used in the middle of a sentence
- Test file covers more of the code, and tests deeper
- Changelog now shows most recent update first
- README has been updated with an explanation of all the commands
- Switches to a versioning approach

###v1.1.2 - 2016-10-12
- File refactor

###v1.1.1 - 2016-10-07
- Fixes scoring between chats
- Adds a lovely 3rd person view

###v1.1.0 - 2016-08-08
- Adds Scoring
- Adds Upvoting
- Adds Downvoting

###v1.0.2 - 2016-08-05
- Lets /hype be called anywhere in a string
	- Adds appropriate tests

###v1.0.1 - 2016-08-01
- Adds new 'Whatchu Did There' gif

###v1.0.0 - 2016-07-27
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
