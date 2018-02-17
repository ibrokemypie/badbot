package plugins

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

func Work(d *discordgo.Session, channelID string) {

	for {
		d.ChannelMessageSend(channelID, ".timely")
		d.ChannelMessageSend(channelID, "]h")
		d.ChannelMessageSend(channelID, "]h 1")
		time.Sleep(time.Hour * 1)
	}
}
