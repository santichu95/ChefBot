package framework

// FindQueue ...
func FindQueue(ctx *Context) []AudioItem {
	if _, ok := ctx.Info.SongQueue[ctx.GuildID]; !ok {
		ctx.Info.SongQueue[ctx.GuildID] = make([]AudioItem, 0)
	}
	return ctx.Info.SongQueue[ctx.GuildID]
}
