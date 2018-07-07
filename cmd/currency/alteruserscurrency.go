package currency

import (
	"ChefBot/framework"
	"fmt"
	"log"
	"math"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// AlterUsersCurrency will parse through ds, and mc to get a value from the message and then
// make a change to the mentioned user. The value will be multiplied with the multipier to allow
// both positive and negative changes.
func AlterUsersCurrency(ds *discordgo.Session, mc *discordgo.Message, ctx *framework.Context, multiplier float64) {
	var value int
	var targetID string

	if len(mc.Mentions) > 0 {
		targetID = mc.Mentions[0].ID
	} else {
		// Error message about misused command
		return
	}

	tokens := strings.Fields(mc.Content)

	// *award <value> <@recipient>
	log.Print(len(tokens))
	if len(tokens) != 3 {
		log.Print("Missing information")
		return
	}

	for _, word := range tokens {
		log.Printf("%v", word)
		if val, err := strconv.Atoi(word); err == nil {
			value = val
			break
		}
	}

	deltaCurrency := (int)(math.Round(multiplier * (float64)(value)))

	ChangeValue(ctx.DatabaseConnection, deltaCurrency, targetID)

	var message string
	if multiplier > 0 {
		message = fmt.Sprintf("<@%v> has received %v :cherry_blossom:. Congratulations!", targetID, value)

	} else {
		message = fmt.Sprintf("<@%v> has lost %v :cherry_blossom:. What did you do?", targetID, value)
	}

	// Chat Message
	embed := &discordgo.MessageEmbed{
		Author:      &discordgo.MessageEmbedAuthor{},
		Color:       0xf47d42, // Peach
		Description: message,
	}

	_, err := ds.ChannelMessageSendEmbed(mc.ChannelID, embed)

	if err != nil {
		log.Printf("Unable to send embedded message, %v", err.Error())
	}
}
