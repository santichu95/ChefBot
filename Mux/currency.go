package mux

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// TODO create function to print error message to discord
// TODO abstract duplicated code into stand alone query functions
// TODO create award function
// TODO create all gambling functions

// GiveCurrency will give a given amount of currency from the author of the message
// to the user who is mentioned.
// TODO Read the transaction value from message
func GiveCurrency(ds *discordgo.Session, mc *discordgo.Message, ctx *Context) {
	log.Printf("Called GiveCurrency")

	// Get transfer value
	transerValue := 0

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

	// Get info from database.
	stmtOut, err := ctx.DatabaseConnection.Prepare("SELECT Value FROM Currency WHERE ID = ?")

	if err != nil {
		log.Printf("Error preparing query, %v", err.Error())
		return
	}

	var fromValue int
	err = stmtOut.QueryRow(fromUserID).Scan(&fromValue)

	// If the user is not in the database add them
	if err == sql.ErrNoRows {
		log.Print("No rows found")
		//Insert the value and print 0
		fromValue = 0
		_, err = ctx.DatabaseConnection.Exec(fmt.Sprintf("INSERT INTO Currency (ID, Value) VALUES(%v, %v)", fromUserID, fromValue))
		if err != nil {
			log.Printf("Error inserting new user into database, %v", err.Error())
			return
		}
	} else if err != nil {
		log.Printf("Error querying database, %v", err.Error())
		return
	}

	// Print information to discord
	log.Printf("**%v** has %v:cherry_blossom:", fromUsername, fromValue)

	var toValue int
	err = stmtOut.QueryRow(toUserID).Scan(&toValue)

	// If the user is not in the database add them
	if err == sql.ErrNoRows {
		log.Print("No rows found")
		//Insert the value and print 0
		toValue = 0
		_, err = ctx.DatabaseConnection.Exec(fmt.Sprintf("INSERT INTO Currency (ID, Value) VALUES(%v, %v)", toUserID, toValue))
		if err != nil {
			log.Printf("Error inserting new user into database, %v", err.Error())
			return
		}
	} else if err != nil {
		log.Printf("Error querying database, %v", err.Error())
		return
	}

	// Print information to discord
	log.Printf("**%v** has %v:cherry_blossom:", toUsername, toValue)
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xf47d42, // Peach
		Description: fmt.Sprintf("%v has gifted %v:cherry_blossom: to **%v**", fromUsername, transerValue, toUsername),
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
