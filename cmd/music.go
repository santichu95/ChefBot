package cmd

import (
	"ChefBot/framework"
	"io"
	"log"

	"github.com/bwmarrin/discordgo"
	"github.com/jonas747/dca"
)

// Play ...
func Play(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {

	// Change these accordingly
	options := dca.StdEncodeOptions
	options.RawOutput = true
	options.Bitrate = 96
	options.Application = "lowdelay"

	downloadURL, _, err := framework.GetVideoInfo("Elephant Gym Finger")
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
