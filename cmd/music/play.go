package music

import (
	"ChefBot/framework"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Play ...
func Play(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
	userInput := strings.Fields(mc.Content)
	if len(userInput) <= 1 {
		// Do nothing if we see no additional input other than the command.
		return
	}
	input := strings.Join(userInput[1:], " ")

	log.Println("Downloading video")
	videoInfo, err := framework.DownloadYoutubeVideo(input)

	if err != nil {
		log.Println("Error downloading video")
		log.Println(err.Error())
		return
	}

	for i, vi := range *videoInfo {
		if i == 0 {
			go func() {
				framework.MaybePlaySong(ds, mc, ctx)
			}()
		}
		framework.Insert(&vi, mc, ctx)
	}
}
