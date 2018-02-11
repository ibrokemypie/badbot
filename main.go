package main

import (
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/pelletier/go-toml"
)

var (
	token   string
	ownerid string
)

func config() {
	config, _ := toml.LoadFile("config.toml")

	token = config.Get("token").(string)
	ownerid = config.Get("ownerid").(string)
}

func main() {
	//Load config
	config()

	//Create discord object
	discord, err := discordgo.New(token)
	if err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Connected")
	}

	//Connect
	err = discord.Open()
	if err != nil {
		fmt.Println("error opening connection,", err)
		return
	}

	//Verify connection
	_, err = discord.User("@me")
	if err != nil {
		fmt.Println("Something went wrong with logging in, check twice if your token is correct.\nYou can do so by editing/deleting config.toml")
		fmt.Println("Press CTRL-C to exit.")
		<-make(chan struct{})
		return
	}
	fmt.Println("Logged in")

	//Init
	botInit(discord)
	fmt.Println("Init completed")

	//Add handlers
	discord.AddHandler(messageCreate)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Press Control + C to quit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
}

func botInit(s *discordgo.Session) {
	s.UpdateStatus(0, "with potato")

	guilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, guild := range guilds {
		s.GuildMemberNickname(guild.ID, "@me", "pretty bad ey")
	}
}

// message is created on any channel that the autenticated bot has access to.
func messageCreate(d *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID != ownerid {
		return
	}

	//If message starts with >Say, say the following text
	if strings.HasPrefix(m.Content, ">say ") {
		s := strings.SplitAfterN(m.Content, " ", 2)
		fmt.Println(s)
		d.ChannelMessageSend(m.ChannelID, (s[1]))
	}

	// If the message is "ping" reply with "Pong!"
	if m.Content == ">ping" {
		d.ChannelMessageSend(m.ChannelID, "pong")
	}

	// If the message is "pong" reply with "Ping!"
	if m.Content == ">pong" {
		d.ChannelMessageSend(m.ChannelID, "ping")
	}
}
