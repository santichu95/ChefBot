package currency

import (
	"ChefBot/framework"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// ShowLeaderBoard will display the leaderboard in chat
// *lb <page>
func ShowLeaderBoard(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
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
