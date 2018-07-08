package gambling

import (
	"ChefBot/cmd/currency"
	"ChefBot/framework"
	"fmt"
	"log"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// BetRoll will retrieve a value and guess from the message, then flip a coin and will then reward/take accordingly.
func BetRoll(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
	log.Printf("Called BetRoll")

	username := mc.Author
	userID := username.ID

	bet := 0

	// Get the value from the message
	for _, word := range strings.Fields(mc.Content) {
		log.Printf("%v", word)
		if val, err := strconv.Atoi(word); err == nil {
			bet = val
		}
	}

	if bet == 0 {
		// Print message
		return
	}

	authorWallet, err := currency.CheckForCurrency(ctx.DatabaseConnection, userID)

	var message string
	if bet > authorWallet {
		message = fmt.Sprintf("<@%v> you do not have enough :cherry_blossom:.", mc.Author.ID)
	} else {
		rand.Seed(time.Now().UnixNano())

		roll := rand.Intn(100) + 1

		log.Printf("Rolled: %v", roll)
		payoutMulti := -1
		message = fmt.Sprintf("<@%v> you rolled a %v.", mc.Author.ID, roll)
		if roll < 66 {
			payoutMulti = -1
		} else if roll < 90 {
			payoutMulti = 2
		} else if roll < 100 {
			payoutMulti = 4
		} else {
			payoutMulti = 10
		}

		currency.ChangeValue(ctx.DatabaseConnection, payoutMulti*bet, userID)
		if payoutMulti < 0 {
			message = fmt.Sprintf("%vThanks for the %v:cherry_blossom:!", message, bet)
		} else {
			message = fmt.Sprintf("%vCongratulations! You won %v:cherry_blossom:", message, (1+payoutMulti)*bet)
		}
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
