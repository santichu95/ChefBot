package mux

import (
	"database/sql"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// TODO create function to print error message to discord
// TODO abstract duplicated code into stand alone query functions
// TODO create award function
// TODO create all gambling functions
// TODO GetValuesFromMessage
// TODO GetMentionsFromMessage

// TODO BetRoll
// TODO Wheel
// TODO Slots

// ShowLeaderBoard will display the leaderboard in chat
// *lb <page>
func ShowLeaderBoard(ds *discordgo.Session, mc *discordgo.Message, ctx *Context) {
	// TODO pagination
	// Get info from database.
	rows, err := ctx.DatabaseConnection.Query("SELECT Value, Username, Discriminator FROM Currency JOIN Users USING (ID) ORDER BY Value DESC limit 9")

	if err != nil {
		log.Printf("Error preparing query, %v", err.Error())
		return
	}

	defer rows.Close()

	var fields []*discordgo.MessageEmbedField
	counter := 1
	for rows.Next() {
		var wealth int
		var username string
		var discriminator int
		err = rows.Scan(&wealth, &username, &discriminator)
		if err != nil {
			log.Printf("Error processing row for leaderboard, %v", err.Error())
			continue
		}
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("#%v **%v#%v**", counter, username, discriminator),
			Value:  fmt.Sprintf("%v:cherry_blossom:", wealth),
			Inline: true,
		})
		counter++
	}

	message := ":cherry_blossom: Leaderboard"
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xf47d42, // Peach
		Fields:      fields,
		Description: message,
	}

	_, err = ds.ChannelMessageSendEmbed(mc.ChannelID, embed)

	if err != nil {
		log.Printf("Unable to send embedded message, %v", err.Error())
	}
}

// BetFlip will retrive a value and guess from the message, then flip a coin and will then reward/take accordingly.
func BetFlip(ds *discordgo.Session, mc *discordgo.Message, ctx *Context) {
	payoutMulti := .9

	log.Printf("Called BetFlip")

	const (
		heads = iota
		tails
	)

	username := mc.Author
	userID := username.ID

	bet := 0
	temp := ""

	// Get the value from the message
	for _, word := range strings.Fields(mc.Content) {
		log.Printf("%v", word)
		if val, err := strconv.Atoi(word); err == nil {
			bet = val
		} else {
			temp = word
		}
	}

	if bet == 0 {
		// Print message
		return
	}

	guess := -1
	if temp != "" && strings.HasPrefix("heads", strings.ToLower(temp)) {
		guess = heads
	}

	if temp != "" && strings.HasPrefix("tails", strings.ToLower(temp)) {
		guess = tails
	}

	if guess == -1 {
		// Print message
		return
	}

	rand.Seed(time.Now().UnixNano())

	flip := rand.Intn(100)

	message := ""
	log.Printf("Rolled: %v\nGuessed: %v", flip, guess)
	// 0 <= tails < 50
	// 50 <= heads < 100
	if (flip < 50 && guess == tails) || (flip >= 50 && guess == heads) {
		log.Print((int)(math.Round(payoutMulti * (float64)(bet))))
		ChangeValue(ctx.DatabaseConnection, (int)(math.Round(payoutMulti*(float64)(bet))), userID)
		message = fmt.Sprintf("<@%v> Correct! You won %v:cherry_blossom:", mc.Author.ID, (int)(math.Round(payoutMulti*(float64)(bet))))
	} else {
		ChangeValue(ctx.DatabaseConnection, -1*bet, userID)
		message = fmt.Sprintf("<@%v> Thanks for the %v:cherry_blossom:! Better luck next time.", mc.Author.ID, bet)
	}

	// Print information to discord
	log.Print(message)
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xf47d42, // Peach
		Description: message,
	}

	_, err := ds.ChannelMessageSendEmbed(mc.ChannelID, embed)

	if err != nil {
		log.Printf("Unable to send embedded message, %v", err.Error())
	}
}

// AwardCurrency will create a given amount of currency and give it to the user mentioned
// Should only be used by Bot admins
func AwardCurrency(ds *discordgo.Session, mc *discordgo.Message, ctx *Context) {
	if mc.Author.ID != "179776524822642688" {
		log.Printf("Award called by %v", mc.Author)
		return
	}

	AlterUsersCurrency(ds, mc, ctx, 1)

	// Chat Message
}

// TakeCurrency will create a given amount of currency and give it to the user mentioned
// Should only be used by Bot admins
func TakeCurrency(ds *discordgo.Session, mc *discordgo.Message, ctx *Context) {
	if mc.Author.ID != "179776524822642688" {
		log.Printf("Award called by %v", mc.Author)
		return
	}

	AlterUsersCurrency(ds, mc, ctx, -1)

}

// GiveCurrency will give a given amount of currency from the author of the message
// to the user who is mentioned.
// TODO Read the transaction value from message
func GiveCurrency(ds *discordgo.Session, mc *discordgo.Message, ctx *Context) {
	log.Printf("Called GiveCurrency")

	fromUsername := mc.Author
	fromUserID := fromUsername.ID

	var toUserID string
	var toUsername *discordgo.User

	if len(mc.Mentions) > 0 {
		toUsername = mc.Mentions[0]
		toUserID = toUsername.ID
	} else {
		// Error message about misused command
		return
	}

	// Sending currency to yourself
	if toUsername == fromUsername {
		return
	}

	// Get transfer value
	transferValue := 0

	// Get the value from the message
	for _, word := range strings.Fields(mc.Content) {
		log.Printf("%v", word)
		if val, err := strconv.Atoi(word); err == nil {
			transferValue = val
			break
		}
	}

	// Either no value specified or 0 given, no point to continue either way.
	if transferValue == 0 {
		return
	}

	// Print information to discord
	fromValue, err := CheckForCurrency(ctx.DatabaseConnection, fromUserID)
	log.Printf("**%v** has %v:cherry_blossom:", fromUsername, fromValue)

	toValue, err := CheckForCurrency(ctx.DatabaseConnection, toUserID)
	log.Printf("**%v** has %v:cherry_blossom:", toUsername, toValue)

	message := fmt.Sprintf("<@%v> has gifted %v:cherry_blossom: to **%v**", fromUserID, transferValue, toUsername)

	if transferValue > fromValue {
		message = fmt.Sprintf("**%v**, you do not have %v:cherry_blossom: to give", fromUsername, transferValue)
	} else {
		// TODO Ensure that these will both happen or neither will happen
		// IDEA Create a change value struct, pass them to a ChangeValue function
		// the function will ensure that they are all in the same transaction
		ChangeValue(ctx.DatabaseConnection, -1*transferValue, fromUserID)
		ChangeValue(ctx.DatabaseConnection, transferValue, toUserID)
	}

	// Print information to discord
	log.Print(message)
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xf47d42, // Peach
		Description: message,
	}

	_, err = ds.ChannelMessageSendEmbed(mc.ChannelID, embed)

	if err != nil {
		log.Printf("Unable to send embeded message, %v", err.Error())
	}
}

// ListUserWallet will list the current value of the users wallet.h
func ListUserWallet(ds *discordgo.Session, mc *discordgo.Message, ctx *Context) {
	log.Printf("Called ListUserWallet")
	// Figure out who to get information about

	var userID string
	var username *discordgo.User

	if len(mc.Mentions) > 0 {
		username = mc.Mentions[0]
		userID = username.ID
	} else {
		username = mc.Author
		userID = username.ID
	}

	value, err := CheckForCurrency(ctx.DatabaseConnection, userID)

	if err != nil {
		// TODO Print error message to discord
		return
	}

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xf47d42, // Peach
		Description: fmt.Sprintf("**%v** has %v:cherry_blossom:", username, value),
	}

	_, err = ds.ChannelMessageSendEmbed(mc.ChannelID, embed)

	if err != nil {
		log.Printf("Unable to send embedded message, %v", err.Error())
	}
}

// TODO Abstract this to CheckIfExistsDatabase( <DB>, <table name>, <Primary Key>)

// CheckForCurrency will query the database and return the value of the targetUserID's wallet
// If that user is not in the database they will be added and given the starting amount of currency
func CheckForCurrency(db *sql.DB, targetUserID string) (int, error) {
	// Get info from database.
	stmtOut, err := db.Prepare("SELECT Value FROM Currency WHERE ID = ?")

	if err != nil {
		log.Printf("Error preparing query, %v", err.Error())
		return 0, err
	}

	var wealth int
	err = stmtOut.QueryRow(targetUserID).Scan(&wealth)

	// If the user is not in the database add them
	if err == sql.ErrNoRows {
		log.Print("No rows found")
		// Insert the value and print 0
		// TODO add a default value for users to start with
		wealth = 0
		_, err = db.Exec(fmt.Sprintf("INSERT INTO Currency (ID, Value) VALUES(%v, %v)", targetUserID, wealth))
		if err != nil {
			log.Printf("Error inserting new user into database, %v", err.Error())
			return 0, err
		}
	} else if err != nil {
		log.Printf("Error querying database, %v", err.Error())
		return 0, err
	}

	return wealth, nil
}

// AlterUsersCurrency will parse through ds, and mc to get a value from the message and then
// make a change to the mentioned user. The value will be multiplied with the multipier to allow
// both positive and negative changes.
func AlterUsersCurrency(ds *discordgo.Session, mc *discordgo.Message, ctx *Context, multiplier float64) {
	var value int
	var targetID string

	if len(mc.Mentions) > 0 {
		targetID = mc.Mentions[0].ID
	} else {
		// Error message about misused command
		return
	}

	tokens := strings.Fields(mc.Content)

	// *award <value> <@recipient>
	log.Print(len(tokens))
	if len(tokens) != 3 {
		log.Print("Missing information")
		return
	}

	for _, word := range tokens {
		log.Printf("%v", word)
		if val, err := strconv.Atoi(word); err == nil {
			value = val
			break
		}
	}

	deltaCurrency := (int)(math.Round(multiplier * (float64)(value)))

	ChangeValue(ctx.DatabaseConnection, deltaCurrency, targetID)

	var message string
	if multiplier > 0 {
		message = fmt.Sprintf("<@%v> has received %v :cherry_blossom:. Congratulations!", targetID, value)

	} else {
		message = fmt.Sprintf("<@%v> has lost %v :cherry_blossom:. What did you do?", targetID, value)
	}

	// Chat Message
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xf47d42, // Peach
		Description: message,
	}

	_, err := ds.ChannelMessageSendEmbed(mc.ChannelID, embed)

	if err != nil {
		log.Printf("Unable to send embedded message, %v", err.Error())
	}
}

// ChangeValue will alter the targetUserID's wallet by currencyDelta on the database pointed to by db
func ChangeValue(db *sql.DB, currencyDelta int, targetUserID string) {
	_, err := db.Exec(fmt.Sprintf("INSERT INTO Currency (ID, Value) VALUES (%v, %v) ON DUPLICATE KEY UPDATE Value=Value + %v;",
		targetUserID, currencyDelta, currencyDelta))
	if err != nil {
		log.Printf("Error changing value(%v) for %v, %v", currencyDelta, targetUserID, err.Error())
		return
	}
}
