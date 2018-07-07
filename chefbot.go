package main

import (
	"ChefBot/cmd"
	"ChefBot/cmd/currency"
	"ChefBot/framework"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// TODO Create logging package
// TODO figure out the best way to do botwide constants i.e color for embedded messages
// TODO Add auth functionality to allow only certain users to call certain functions

// Version is a constant that store the version of ChefBot
const Version = "v0.0.1-alpha"

var (
	// Token is the Discord Bot Token
	Token string

	// Router will be the multiplexer for all of the commands
	Router = framework.NewMux()
)

func init() {
	Token = os.Getenv("DG_TOKEN")

	if Token == "" {
		flag.StringVar(&Token, "t", "", "Discord Authentication Token")
		flag.Parse()
	}

	Router.ConnectDB("config.secret")
	addAllRoutes(Router)
}

func main() {
	dg, err := discordgo.New(Token)
	if err != nil {
		fmt.Println("Error creating Discrod session", err)
		return
	}

	err = dg.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	dg.AddHandler(Router.OnMessageCreate)

	fmt.Println("Bot is now running. Press CTRL-C to exit.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	dg.Close()
}

// addAllRoutes will add all of the commands to the mux
func addAllRoutes(r *framework.Mux) error {
	r.Route([]string{"$"}, "Display value of users wallet", currency.ListUserWallet)
	r.Route([]string{"give"}, "Give currency to another user", currency.GiveCurrency)
	r.Route([]string{"award"}, "Award currency to a user", currency.AwardCurrency)
	r.Route([]string{"take"}, "Take currency from a user", currency.TakeCurrency)
	r.Route([]string{"bf", "betflip"}, "Make a bet on a flip of a coin", cmd.BetFlip)
	r.Route([]string{"leaderboard", "lb"}, "Show a leaderboard of currency for the server", currency.ShowLeaderBoard)
	r.Route([]string{"br", "betroll"}, "Make a bet on the roll of a d100", cmd.BetRoll)
	r.Route([]string{"test"}, "used to test commands", cmd.Play)
	r.Route([]string{"summon"}, "Summons the bot into the voice channel you are in", cmd.Summon)
	r.Route([]string{"disconnect"}, "Disconnects the bot from the voice channel", cmd.Disconnect)

	addPersonalCommands(r)

	return nil
}
