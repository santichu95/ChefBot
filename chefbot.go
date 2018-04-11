package main

import (
	"ChefBot/cmd"
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
	r.Route("$", "Display value of users wallet", cmd.ListUserWallet)
	r.Route("give", "Give currency to another user", cmd.GiveCurrency)
	r.Route("award", "Award currency to a user", cmd.AwardCurrency)
	r.Route("take", "Take currency from a user", cmd.TakeCurrency)
	r.Route("bf", "Make a bet on a flip of a coin", cmd.BetFlip)
	r.Route("betflip", "Make a bet on a flip of a coin", cmd.BetFlip)
	r.Route("lb", "Show a leaderboard of currency for the server", cmd.ShowLeaderBoard)
	r.Route("leaderboard", "Show a leaderboard of currency for the server", cmd.ShowLeaderBoard)
	r.Route("br", "Make a bet on the roll of a d100", cmd.BetRoll)
	r.Route("betroll", "Make a bet on the roll of a d100", cmd.BetRoll)
	return nil
}
