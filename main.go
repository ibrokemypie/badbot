package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/ibrokemypie/badbot/plugins"
	"github.com/pelletier/go-toml"
)

var (
	confFile = "config.toml"
	conf     *toml.Tree
	token    string
	ownerid  string
	nickname string
	status   string
	image    string
	workchan string
	replies  bool
	err      error
)

func config() {
	conf, err = toml.LoadFile(confFile)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	token = conf.Get("token").(string)
	ownerid = conf.Get("ownerid").(string)
	nickname = conf.Get("nickname").(string)
	status = conf.Get("status").(string)
	image = conf.Get("image").(string)
	replies = conf.Get("replies").(bool)
	workchan = conf.GetDefault("workchan", "").(string)
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
	//discord.AddHandler(messageReactionAdd)

	// Wait here until CTRL-C or other term signal is received.
	fmt.Println("Press Control + C to quit")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	// Cleanly close down the Discord session.
	discord.Close()
	f, _ := os.Create(confFile)
	f.WriteString(conf.String())
	plugins.SaveData()
}

func botInit(s *discordgo.Session) {
	s.UpdateStatus(0, status)

	// Dont need to update this on launch I think, or any tbf but this one gets limited
	//s.UserUpdate("", "", "", image, "")

	guilds, err := s.UserGuilds(100, "", "")
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, guild := range guilds {
		s.GuildMemberNickname(guild.ID, "@me", nickname)
	}
	if workchan != "" {
		go plugins.Work(s, conf.Get("workchan").(string))
	}
	plugins.Googleinit(s)
	plugins.LoadData()
}

// message is created on any channel that the autenticated bot has access to.
func messageCreate(d *discordgo.Session, m *discordgo.MessageCreate) {
	if replies == true {
		go plugins.Replies(d, m, conf)
	}
	go plugins.Commands(d, m, conf)
}

func messageReactionAdd(d *discordgo.Session, m *discordgo.MessageReactionAdd) {
	fmt.Println(m.MessageID)
	fmt.Println(m.Emoji.ID)
	fmt.Println(m.Emoji.Name)
	if m.Emoji.ID != "" {
		d.MessageReactionAdd(m.ChannelID, m.MessageID, m.Emoji.Name+":"+m.Emoji.ID)
	} else {
		d.MessageReactionAdd(m.ChannelID, m.MessageID, m.Emoji.Name)
	}
}
