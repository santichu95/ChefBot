package framework

import (
	"log"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
)

// Insert will add song to the queue for a specific server.
func Insert(vi *VideoInfo, mc *discordgo.Message, ctx *Context) {
	if _, ok := ctx.Info.SongQueue[mc.GuildID]; !ok {
		ctx.Info.SongQueue[mc.GuildID] = []VideoInfo{}
		ctx.Info.AudioIsPlaying[mc.GuildID] = false
	}
	queue := append(ctx.Info.SongQueue[mc.GuildID], *vi)
	ctx.Info.SongQueue[mc.GuildID] = queue
}

// MaybePlaySong will check if there is anything playing on the related Guild, if there is not and there is something in the queue it will start playing the next video in the queue.
func MaybePlaySong(ds *discordgo.Session, mc *discordgo.Message, ctx *Context) {
	if ctx.Info.AudioIsPlaying[mc.GuildID] {
		return
	}

	vc, err := JoinUserVoiceChannel(ds, mc.Author.ID)
	if err != nil {
		log.Println("Error connecting to voice channel")
		log.Println(err.Error())
	}

	queue := ctx.Info.SongQueue[mc.GuildID]

	if len(queue) == 0 {
		return
	}

	frontOfQueue := queue[0]
	stopPlaying := make(chan bool)
	go func() {
		ctx.Info.AudioIsPlaying[mc.GuildID] = true
		dgvoice.PlayAudioFile(vc, "audio_cache/"+frontOfQueue.ID, stopPlaying)
		ctx.Info.AudioIsPlaying[mc.GuildID] = false
		// Remove recently finished song from queue.
		RemoveSong(mc.GuildID, ctx.Info)
		// Attempt to play the next thing in the queue.
		MaybePlaySong(ds, mc, ctx)
	}()
}

// RemoveSong will remove the song at the front of the queue for the related Guild
func RemoveSong(guildID string, store *Store) {

	queue := store.SongQueue[guildID]

	if len(queue) == 0 {
		return
	}

	store.SongQueue[guildID] = queue[1:]
}
