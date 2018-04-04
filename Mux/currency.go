package mux

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// ListUserWallet will list the current value of the users wallet.
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

	// Get info from database.
	stmtOut, err := ctx.DatabaseConnection.Prepare("SELECT Value FROM Currency WHERE ID = ?")

	if err != nil {
		log.Printf("Error preparing query, %v", err.Error())
		return
	}

	var value int
	err = stmtOut.QueryRow(userID).Scan(&value)

	// If the user is not in the database add them
	if err == sql.ErrNoRows {
		log.Print("No rows found")
		//Insert the value and print 0
		value = 0
		_, err = ctx.DatabaseConnection.Exec(fmt.Sprintf("INSERT INTO Currency (ID, Value) VALUES(%v, %v)", userID, value))
		if err != nil {
			log.Printf("Error inserting new user into database, %v", err.Error())
			return
		}
	} else if err != nil {
		log.Printf("Error querying database, %v", err.Error())
		return
	}

	// Print information to discord
	log.Printf("**%v** has %v:cherry_blossom:", username, value)

	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xf47d42, // Peach
		Description: fmt.Sprintf("**%v** has %v:cherry_blossom:", username, value),
	}

	_, err = ds.ChannelMessageSendEmbed(mc.ChannelID, embed)

	if err != nil {
		log.Printf("Unable to send embeded message, %v", err.Error())
	}
}
