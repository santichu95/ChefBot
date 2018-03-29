package mux

import (
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// Route holds information about a specific message route handler
type Route struct {
	Pattern     string
	Description string
	Help        string
	Run         HandlerFunc
}

// Context holds extra data that is passed along to route handlers
// This way processing some of this only need to happen once
type Context struct {
	Fields          []string
	Content         string
	GuildID         string
	IsDirected      bool
	IsPrivate       bool
	HasPrefix       bool
	HasMention      bool
	HasMentionFirst bool
}

// HandlerFunc is the function signature required for a message route handler
type HandlerFunc func(*discordgo.Session, *discordgo.Message, *Context)

// Mux is the main struct for all mux methods.
type Mux struct {
	Routes  []*Route
	Default *Route
	Prefix  string
}

// New returns a new Discord message route mux
func New() *Mux {
	m := &Mux{}
	m.Prefix = "*"
	return m
}

// Route allows you to register a route
func (m *Mux) Route(pattern, desc string, cb HandlerFunc) (*Route, error) {

	r := Route{}
	r.Pattern = pattern
	r.Description = desc
	r.Run = cb
	m.Routes = append(m.Routes, &r)

	return &r, nil
}

// OnMessageCreate is a DiscordGo Event Handler function. This function will
// recieve all Discord message and parse themm for matches to registed routes.
func (m *Mux) OnMessageCreate(ds *discordgo.Session, mc *discordgo.MessageCreate) {
	var err error

	// Ignore all messages created by the Bot
	if mc.Author.ID == ds.State.User.ID {
		return
	}

	// Creating a context struct
	ctx := &Context{
		Content: strings.TrimSpace(mc.Content),
	}

	// If the message does not start with the bot prefix do nothing
	if !strings.HasPrefix(ctx.Content, m.Prefix) {
		log.Printf("Message missing bot prefix, ", ctx.Content)
		return
	}

	// Fetch the channel for this Message
	var c *discordgo.Channel
	c, err = ds.State.Channel(mc.ChannelID)
	if err != nil {
		// Try fetching via REST API
		c, err = ds.Channel(mc.ChannelID)
		if err != nil {
			log.Printf("unable to fetch Channel for Message,", err)
		} else {
			// Attempt to add this channel into our State
			err = ds.State.ChannelAdd(c)
			if err != nil {
				log.Printf("error updatin State with Channel,", err)
			}

			// Add Channel info into Context
			ctx.GuildID = c.GuildID
			if c.Type == discordgo.ChannelTypeDM {
				ctx.IsPrivate = true
				ctx.IsDirected = true
			}
		}
	}
}
