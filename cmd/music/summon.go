package music

import (
	"ChefBot/framework"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Summon will attempt to connect to the voice channel of the user who invoked the message.
// If that user is not in a channel the bot will print an error message to the chat.
func Summon(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
	// Look for user voice channel
	// Attempt to join user voice channel
	_, err := framework.JoinUserVoiceChannel(ds, mc.Author.ID)

	if err != nil {
		message := "You are not in a voice channel!"
		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Title:       "Error",
			Color:       0xFF0000, // Warning Red
			Description: message,
		}
		_, err = ds.ChannelMessageSendEmbed(mc.ChannelID, embed)
		if err != nil {
			log.Println(err)
		}
	}
	log.Println(ds.VoiceConnections)
}
