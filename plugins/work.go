package plugins

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

var (
	channelID = "231252749218611200"
)

func Work(d *discordgo.Session) {

	for {
		d.ChannelMessageSend(channelID, ".timely")
		d.ChannelMessageSend(channelID, "]h")
		d.ChannelMessageSend(channelID, "]h 1")
		time.Sleep(time.Hour * 1)
	}
}
