package currency

import (
	"ChefBot/framework"
	"log"

	"github.com/bwmarrin/discordgo"
)

// TakeCurrency will create a given amount of currency and give it to the user mentioned
// Should only be used by Bot admins
func TakeCurrency(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context) {
	if mc.Author.ID != "179776524822642688" {
		log.Printf("Award called by %v", mc.Author)
		return
	}

	AlterUsersCurrency(ds, mc, ctx, -1)

}
