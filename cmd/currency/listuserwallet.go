package currency

import (
	"ChefBot/framework"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// ListUserWallet will list the current value of the users wallet.h
func ListUserWallet(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
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
