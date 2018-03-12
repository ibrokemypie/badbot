package plugins

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
	toml "github.com/pelletier/go-toml"
)

func Commands(d *discordgo.Session, m *discordgo.MessageCreate, conf *toml.Tree) {
	// If the message is ">spin" send a spinner
	// if strings.HasPrefix(m.Content, ">spin ") || strings.HasPrefix(m.Content, ">spinner ") {
	// s := strings.SplitAfterN(m.Content, " ", 2)
	// fmt.Println(s)
	//
	// n, err := strconv.Atoi(s[1])
	// if err != nil {
	// fmt.Println(err)
	// d.ChannelMessageSend(m.ChannelID, err.Error())
	// return
	// }
	//
	// Spinner(d, m, n)
	// }

	// If the message is ">git" link to github
	if m.Content == ">git" || m.Content == ">github" {
		d.ChannelMessageSend(m.ChannelID, "https://github.com/ibrokemypie/badbot")
	}

	// If the message is ">stats" send ram and cpu usage
	if m.Content == ">s" || m.Content == ">stats" {
		d.ChannelMessageSend(m.ChannelID, Stats())
	}

	// If the message is ">help" return help
	if m.Content == ">help" || m.Content == ">h" {
		d.ChannelMessageSend(m.ChannelID, "Only god can save you.\n Try >git")
	}

	// If the message is ">nangs" define nangs"
	if m.Content == ">nang" || m.Content == ">nangs" {
		d.ChannelMessageSend(m.ChannelID, "```\nAn Australian slang term for a Nitrous oxide bulb, derived from the sound distortion that occurs when one is under the influence of the drug.\n```")
	}

	// If the message is ">qoohme" link to qooh.me "
	if m.Content == ">qooh" || m.Content == ">qoohme" {
		d.ChannelMessageSend(m.ChannelID, "http://qooh.me/ibrokemypie")
	}

	// If the message is ">sarahah" link to sarahah"
	if m.Content == ">sarahah" {
		d.ChannelMessageSend(m.ChannelID, "https://ibrokemypie.sarahah.com")
	}

	// If the message is ">woof" send randoim dog"
	if m.Content == ">woof" || m.Content == ">dog" {
		d.ChannelMessageSend(m.ChannelID, Woof())
	}

	// If the message is ">meow" send random cat"
	if m.Content == ">meow" || m.Content == ">cat" {
		d.ChannelMessageSend(m.ChannelID, Meow())
	}

	// If the message is ">ping" reply with "Pong!"
	if m.Content == ">ping" {
		d.ChannelMessageSend(m.ChannelID, "pong")
	}

	// If the message is ">pong" reply with "Ping!"
	if m.Content == ">pong" {
		d.ChannelMessageSend(m.ChannelID, "ping")
	}

	if m.Author.ID == conf.Get("ownerid").(string) {
		//If message starts with >say, say the following text
		if strings.HasPrefix(m.Content, ">say ") {
			s := strings.SplitAfterN(m.Content, " ", 2)
			fmt.Println(s)
			d.ChannelMessageSend(m.ChannelID, (s[1]))
		}

		//If message starts with >search, google the following text
		if strings.HasPrefix(m.Content, ">search ") || strings.HasPrefix(m.Content, ">google ") {
			s := strings.SplitAfterN(m.Content, " ", 2)
			fmt.Println(s)
			d.ChannelMessageSend(m.ChannelID, Search(s[1], 1, conf.Get("apiKey").(string), conf.Get("engineid").(string)))
		}

		//If message starts with >game, say the following text
		if strings.HasPrefix(m.Content, ">pfp") || strings.HasPrefix(m.Content, ">avatar") {
			fmt.Println(m.Content)
			img := m.Message.Attachments[0].URL

			baseimg := EncodeImage(img)

			conf.Set("image", baseimg)
			_, err := d.UserUpdate("", "", "", baseimg, "")
			if err != nil {
				fmt.Println(err)
				d.ChannelMessageSend(m.ChannelID, err.Error())
			}
		}

		//If message starts with >game, say the following text
		if strings.HasPrefix(m.Content, ">game ") || strings.HasPrefix(m.Content, ">status ") {
			s := strings.SplitAfterN(m.Content, " ", 2)
			fmt.Println(s)
			conf.Set("status", s[1])
			d.UpdateStatus(0, s[1])
		}

		//If message starts with >game, say the following text
		if strings.HasPrefix(m.Content, ">name ") || strings.HasPrefix(m.Content, ">nick ") {
			s := strings.SplitAfterN(m.Content, " ", 2)
			fmt.Println(s)
			guilds, err := d.UserGuilds(100, "", "")
			if err != nil {
				fmt.Println(err)
				return
			}
			for _, guild := range guilds {
				conf.Set("nickname", s[1])
				d.GuildMemberNickname(guild.ID, "@me", s[1])
			}
		}
	}
}
