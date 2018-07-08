package music

import (
	"ChefBot/framework"
	"log"

	"github.com/bwmarrin/discordgo"
)

// Disconnect will have the bot disconnect from the voice channel that it is currently in.
// This should be limited to admin/moderators
func Disconnect(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
	log.Println("Disconnect")
	log.Println(ctx.GuildID)
	if vc, ok := ds.VoiceConnections[ctx.GuildID]; ok {
		err := vc.Disconnect()
		if err != nil {
			log.Println(err)
		}

	} else {
		message := "ChefBot is not connected to a voice channel!"
		embed := &discordgo.MessageEmbed{
			Author:      &discordgo.MessageEmbedAuthor{},
			Title:       "Error",
			Color:       0xFF0000, // Warning Red
			Description: message,
		}
		_, err := ds.ChannelMessageSendEmbed(mc.ChannelID, embed)
		if err != nil {
			log.Println(err)
		}
	}
}
