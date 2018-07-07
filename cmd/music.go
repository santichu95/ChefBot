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

func findQueue(ctx *framework.Context) []framework.AudioItem {
	if _, ok := ctx.Info.SongQueue[ctx.GuildID]; !ok {
		ctx.Info.SongQueue[ctx.GuildID] = make([]framework.AudioItem, 0)
	}
	return ctx.Info.SongQueue[ctx.GuildID]
}

// Play ...
func Play(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
	/*
		// Find the queue for the server(GuildID)
		queue := findQueue(ctx)

		// Add song to queue
		if len(queue) > 0 {
		} else {
			// Add song to queue and start audio transmition
		}
		// Create one if it does not exist
		// Add song to queue, if empty beforehand start audio play back.
		// TODO(sandaluz) figure out how to call another function at the end of streaming
	*/

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
