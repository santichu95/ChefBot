package framework

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/bwmarrin/discordgo"
	// import a mysql driver for database/sql
	_ "github.com/go-sql-driver/mysql"
)

// Route holds information about a specific message route handler
type Route struct {
	Pattern     []string
	Description string
	Help        string
	Run         HandlerFunc
}

// Context holds extra data that is passed along to route handlers
// This way processing some of this only need to happen once
type Context struct {
	Fields             []string
	Content            string
	DatabaseConnection *sql.DB
	GuildID            string
	IsDirected         bool
	IsPrivate          bool
	HasPrefix          bool
	HasMention         bool
	HasMentionFirst    bool
	Info               *Store
}

// HandlerFunc is the function signature required for a message route handler
type HandlerFunc func(*discordgo.Session, *discordgo.Message, *Context)

// Mux is the main struct for all mux methods.
type Mux struct {
	Routes             []*Route
	Default            *Route
	Prefix             string
	DatabaseConnection *sql.DB
	Info               *Store
}

// Store is a generic place to maintain long term information.
type Store struct {
	SongQueue      map[string][]VideoInfo
	AudioIsPlaying map[string]bool
}

// NewMux returns a new Discord message route mux
func NewMux() *Mux {
	m := &Mux{}
	m.Prefix = "*"
	emptyInfo := Store{
		SongQueue:      make(map[string][]VideoInfo),
		AudioIsPlaying: make(map[string]bool)}
	m.Info = &emptyInfo
	return m
}

// TODO return error and allow calling function to handle it

// ConnectDB will connect the multiplexer to a database that holds information
// about the users on the discord server
func (m *Mux) ConnectDB(filename string) {
	file, err := os.Open(filename)

	if err != nil {
		log.Printf("Unable to open %v", filename)
		return
	}

	defer file.Close()

	data := struct {
		Username, Password string
	}{}

	decoder := json.NewDecoder(file)

	err = decoder.Decode(&data)

	if err != nil {
		log.Printf("Unable to parse %v", filename)
		log.Println(err.Error())
		os.Exit(1)
	}

	dsn := fmt.Sprintf("%v:%v@/chefbot", data.Username, data.Password)

	m.DatabaseConnection, err = sql.Open("mysql", dsn)

	if err != nil {
		log.Print("Error opening database")
		log.Println(err.Error())
		os.Exit(1)
	}

	err = m.DatabaseConnection.Ping()

	if err != nil {
		log.Print("Error connecting to database")
		log.Println(err.Error())
		os.Exit(1)
	}

}

// Route allows you to register a route
func (m *Mux) Route(pattern []string, desc string, cb HandlerFunc) (*Route, error) {

	r := Route{}
	r.Pattern = pattern
	r.Description = desc
	r.Run = cb
	m.Routes = append(m.Routes, &r)

	return &r, nil
}

// Match attempts to find the route for the given message
func (m *Mux) Match(msg string) (*Route, error) {

	// Tokenize the msg string into a slice of words
	command := strings.Fields(msg)[0]

	for _, routeValue := range m.Routes {
		for _, pattern := range routeValue.Pattern {
			if pattern == command[1:] {
				return routeValue, nil
			}
		}
	}

	return nil, errors.New("No route found")
}

// OnMessageCreate is a DiscordGo Event Handler function. This function will
// receive all Discord message and parse themm for matches to registered routes.
func (m *Mux) OnMessageCreate(ds *discordgo.Session, mc *discordgo.MessageCreate) {

	// Ignore all messages created by the Bot
	if mc.Author.ID == ds.State.User.ID {
		return
	}

	// Check if the user has been seen
	// if not seen, insert into DB
	UpsertUserInfo(m.DatabaseConnection, mc.Author)

	// Creating a context struct
	ctx := &Context{
		Content:            strings.TrimSpace(mc.Content),
		DatabaseConnection: m.DatabaseConnection,
		Info:               m.Info,
	}

	// TODO Add server specific prefixes
	// If the message does not start with the bot prefix do nothing
	if !strings.HasPrefix(ctx.Content, m.Prefix) {
		log.Printf("Message missing bot prefix, %v", ctx.Content)
		log.Printf("Author, %v;Time %v", mc.Author, mc.Timestamp)
		return
	}

	// Fetch the channel for this Message
	var c *discordgo.Channel
	var err error

	c, err = ds.State.Channel(mc.ChannelID)
	if err == nil {
		// Try fetching via REST API
		c, err = ds.Channel(mc.ChannelID)
		if err != nil {
			log.Printf("unable to fetch Channel for Message, %v", err)
		} else {
			log.Println("OnMessageCreate")
			// Attempt to add this channel into our State
			err = ds.State.ChannelAdd(c)
			if err != nil {
				log.Printf("error updatin State with Channel, %v", err)
			}

			// Add Channel info into Context
			ctx.GuildID = c.GuildID
			if c.Type == discordgo.ChannelTypeDM {
				ctx.IsPrivate = true
				ctx.IsDirected = true
			}
		}
	} else {
		log.Println(err)
	}

	// Run the route that was found
	r, err := m.Match(ctx.Content)
	if err == nil {
		r.Run(ds, mc.Message, ctx)
		return
	}

	log.Printf(err.Error())
	// TODO Add a help message mentioning the unknown command
}

// TODO Create a way to groups the routes, i.e. not list every single route

// ListRoutes will list all of the routes into the chat
func (m *Mux) ListRoutes(ds *discordgo.Session) {
	log.Printf("Printing all of the routes to chat")
}

// UpsertUserInfo will query the database and return the value of the targetUserID's wallet
// If that user is not in the database they will be added and given the starting amount of currency
func UpsertUserInfo(db *sql.DB, user *discordgo.User) error {
	// Upsert User info into database
	_, err := db.Exec("INSERT INTO Users (ID, Username, Discriminator) VALUES(?, ?, ?) ON DUPLICATE KEY UPDATE Username = ?, Discriminator = ?;",
		user.ID, user.Username, user.Discriminator, user.Username, user.Discriminator)

	if err != nil {
		log.Printf("Unable to upsert user info, %v", err.Error())
		return err
	}

	return nil
}
