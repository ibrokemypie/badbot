package plugins

import (
	"strings"

	"github.com/bwmarrin/discordgo"
	toml "github.com/pelletier/go-toml"
)

func Replies(d *discordgo.Session, m *discordgo.MessageCreate, conf *toml.Tree) {

	if m.Author.ID == "361540375443275777" {
		//If message starts with >game, say the following text
		if strings.HasSuffix(m.Content, "has left the server.") {
			d.ChannelMessageSend(m.ChannelID, "f")
		}
	}
}
