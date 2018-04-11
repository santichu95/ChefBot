[![Go Report Card](https://goreportcard.com/badge/github.com/santichu95/ChefBot)](https://goreportcard.com/report/github.com/santichu95/ChefBot)

# ChefBot
Discord bot centered around the currency system.
### TODOs
| Filename | line # | TODO
|:------|:------:|:------
| cmd/currency.go | 15 | create function to print error message to discord
| cmd/currency.go | 16 | GetValuesFromMessage
| cmd/currency.go | 17 | GetMentionsFromMessage
| cmd/currency.go | 22 | pagination
| cmd/currency.go | 94 | Read the transaction value from message
| cmd/currency.go | 146 | Ensure that these will both happen or neither will happen
| cmd/currency.go | 187 | Print error message to discord
| cmd/currency.go | 204 | Abstract this to CheckIfExistsDatabase( <DB>, <table name>, <Primary Key>)
| cmd/currency.go | 224 | add a default value for users to start with
| cmd/gambling.go | 16 | create all gambling functions
| cmd/gambling.go | 17 | Wheel
| cmd/gambling.go | 18 | Slots
| framework/mux.go | 57 | return error and allow calling function to handle it
| framework/mux.go | 151 | Add server specific prefixes
| framework/mux.go | 192 | Add a help message mentioning the unknown command
| framework/mux.go | 195 | Create a way to groups the routes, i.e. not list every single route
