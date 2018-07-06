package cmd

import (
	"ChefBot/framework"
	"io"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
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

// Play ...
func Play(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {

	// Change these accordingly
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	log.Println("Grabbing youtubeURL")
	youtubeURL, err := framework.ParseYoutubeInput(strings.Join(strings.Fields(mc.Content)[1:], " "))
	log.Println("Have youtubeURL")
	if err != nil {
		log.Println("Error searching for the YT Url")
		log.Println(err.Error())
	}

	log.Println("Grabbing downloadURL")
	downloadURL, _, err := framework.GetVideoInfo(youtubeURL)
	if err != nil {
		log.Println("Error 2")
		log.Println(err.Error())
	}

	encodingSession, err := dca.EncodeFile(downloadURL, options)

	if err != nil {
		log.Println("Error 3")
	}
	defer encodingSession.Cleanup()

	done := make(chan error)

	vc, _ := framework.JoinUserVoiceChannel(ds, mc.Author.ID)

	dca.NewStream(encodingSession, vc, done)
	err = <-done
	if err != nil && err != io.EOF {
		log.Println("Error 4")
	}
}
