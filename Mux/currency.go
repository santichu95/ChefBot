package mux

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// ListUserWallet will list the current value of the users wallet.
func ListUserWallet(ds *discordgo.Session, mc *discordgo.Message, ctx *Context) {
	log.Printf("Called ListUserWallet")
}
