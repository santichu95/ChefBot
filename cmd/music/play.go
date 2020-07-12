package music

import (
	"ChefBot/framework"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Play ...
func Play(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
	input := strings.Join(strings.Fields(mc.Content)[1:], " ")

	log.Println("Downloading video")
	videoInfo, err := framework.DownloadYoutubeVideo(input)

	if err != nil {
		log.Println("Error downloading video")
		log.Println(err.Error())
		return
	}

	go func() {
		framework.Insert(videoInfo, mc, ctx)
		framework.MaybePlaySong(ds, mc, ctx)
	}()
}
