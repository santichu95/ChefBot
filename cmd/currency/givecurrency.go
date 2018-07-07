package currency

import (
	"ChefBot/framework"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// GiveCurrency will give a given amount of currency from the author of the message
// to the user who is mentioned.
// TODO Read the transaction value from message
func GiveCurrency(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
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
		log.Printf("Unable to send embedded message, %v", err.Error())
	}
}
