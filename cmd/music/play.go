package music

import (
	"ChefBot/framework"
	"io"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

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
