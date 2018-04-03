package main

import (
	mux "ChefBot/Mux"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
)

// Version is a constant that store the version of ChefBot
const Version = "v0.0.1-alpha"

var (
	// Token is the Discord Bot Token
	Token string

	// Router will be the multiplexer for all of the commands
	Router = mux.New()
)

func init() {
	Token = os.Getenv("DG_TOKEN")

	if Token == "" {
		flag.StringVar(&Token, "t", "", "Discord Authentication Token")
		flag.Parse()
	}

	mux.AddAllCommands(Router)
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
