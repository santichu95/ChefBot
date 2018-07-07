package cmd

import (
	"ChefBot/cmd/currency"
	"ChefBot/framework"
	"fmt"
	"log"
	"math"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// TODO create all gambling functions
// TODO Wheel
// TODO Slots
// TODO

// BetFlip will retrive a value and guess from the message, then flip a coin and will then reward/take accordingly.
func BetFlip(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
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
		currency.ChangeValue(ctx.DatabaseConnection, (int)(math.Round(payoutMulti*(float64)(bet))), userID)
		message = fmt.Sprintf("<@%v> Correct! You won %v:cherry_blossom:", mc.Author.ID, (int)(math.Round((1+payoutMulti)*(float64)(bet))))
	} else {
		currency.ChangeValue(ctx.DatabaseConnection, -1*bet, userID)
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
