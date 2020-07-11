package music

import (
	"ChefBot/framework"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Play ...
func Play(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
	log.Println("Finding youtubeURL")
	youtubeURL, err := framework.ParseYoutubeInput(strings.Join(strings.Fields(mc.Content)[1:], " "))

	if err != nil {
		log.Println("Error searching for the YT Url")
		log.Println(err.Error())
	}

	log.Println("Downloading video")
	videoInfo, err := framework.DownloadVideo(youtubeURL)
	if err != nil {
		log.Println("Error downloading video")
		log.Println(err.Error())
	}

	go func() {
		framework.Insert(videoInfo, mc, ctx)
		framework.MaybePlaySong(ds, mc, ctx)
	}()
}
