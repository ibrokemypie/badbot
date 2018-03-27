package plugins

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bwmarrin/discordgo"
	toml "github.com/pelletier/go-toml"
)

func Commands(d *discordgo.Session, m *discordgo.MessageCreate, conf *toml.Tree, trustedusers map[string]bool) {
	// If the message is ">spin" send a spinner
	// if strings.HasPrefix(m.Content, ">spin ") || strings.HasPrefix(m.Content, ">spinner ") {
	// s := strings.SplitN(m.Content, " ", 2)
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

	if strings.HasPrefix(m.Content, ">>> ") {
		s := strings.SplitN(m.Content, " ", 3)
		if len(s) < 3 {
			d.ChannelMessageSend(m.ChannelID, "Usage is ``>>> name quote``")
			return
		}
		fmt.Println(s)
		s[2] = strings.Replace(s[2], "\n", "\\n", -1)
		WriteQuote(s[1], s[2])
		d.ChannelMessageSend(m.ChannelID, "Quote "+s[1]+" added.")
	}

	if strings.HasPrefix(m.Content, ">> ") {
		s := strings.SplitN(m.Content, " ", 2)
		if len(s) != 2 {
			d.ChannelMessageSend(m.ChannelID, "Usage is ``>> quotename``")
			return
		}
		fmt.Println(s)
		q, i := ReadQuote(s[1])
		if i != 0 {
			q = strings.Replace(q, "\\n", "\n", -1)
			d.ChannelMessageSend(m.ChannelID, "``"+s[1]+" "+strconv.Itoa(i)+"``"+": \n"+q)
		} else {
			d.ChannelMessageSend(m.ChannelID, "No quotes with that name found.")
		}
	}

	if strings.HasPrefix(m.Content, ">qid ") || strings.HasPrefix(m.Content, ">quoteid ") {
		s := strings.SplitN(m.Content, " ", 2)
		if len(s) != 2 {
			d.ChannelMessageSend(m.ChannelID, "Usage is ``> quoteid``")
			return
		}
		fmt.Println(s)
		q, i, n := ReadQuoteID(s[1])
		if i != 0 {
			q = strings.Replace(q, "\\n", "\n", -1)
			d.ChannelMessageSend(m.ChannelID, "``"+n+" "+strconv.Itoa(i)+"``"+": \n"+q)
		} else {
			d.ChannelMessageSend(m.ChannelID, q)
		}
	}

	if strings.HasPrefix(m.Content, ">qdel ") || strings.HasPrefix(m.Content, ">qrm ") {
		s := strings.SplitN(m.Content, " ", 2)
		if len(s) != 2 {
			d.ChannelMessageSend(m.ChannelID, "Usage is ``> qdel``")
			return
		}
		fmt.Println(s)
		r := RemoveQuote(s[1])
		d.ChannelMessageSend(m.ChannelID, r)
	}

	if trustedusers[m.Author.ID] {
		//If message starts with >say, say the following text
		if strings.HasPrefix(m.Content, ">say ") {
			s := strings.SplitN(m.Content, " ", 2)
			fmt.Println(s)
			d.ChannelMessageSend(m.ChannelID, (s[1]))
		}

		//If message starts with >search, google the following text
		if strings.HasPrefix(m.Content, ">search ") || strings.HasPrefix(m.Content, ">google ") {
			s := strings.SplitN(m.Content, " ", 2)
			fmt.Println(s)
			go SearchGoogle(s[1], 10, conf.Get("googleapi").(string), conf.Get("engineid").(string), d, m)
		}

		//If message starts with >yt, youtube search the following text
		if strings.HasPrefix(m.Content, ">yt ") || strings.HasPrefix(m.Content, ">youtube ") {
			s := strings.SplitN(m.Content, " ", 2)
			fmt.Println(s)
			go SearchYoutube(s[1], 10, conf.Get("youtubeapi").(string), d, m)
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
			s := strings.SplitN(m.Content, " ", 2)
			fmt.Println(s)
			conf.Set("status", s[1])
			d.UpdateStatus(0, s[1])
		}

		//If message starts with >trust, add id to trusted listt
		if strings.HasPrefix(m.Content, ">trust ") {
			//s := strings.SplitN(m.Content, " ", 2)
			fmt.Println(m.Content)
			user := m.Mentions[0].ID
			trustedusers[user] = true

			conf.Set("trustedusers", trustedusers)
			f, _ := os.Create("config.toml")
			f.WriteString(conf.String())
		}

		//If message starts with >untrust, remove id from trusted list
		if strings.HasPrefix(m.Content, ">untrust ") {
			fmt.Println(m.Content)
			user := m.Mentions[0].ID
			if trustedusers[user] {
				delete(trustedusers, user)

				conf.Set("trustedusers", trustedusers)
				f, _ := os.Create("config.toml")
				f.WriteString(conf.String())
			}
		}

		//If message starts with >trusted, send list of trusted ids
		if m.Content == ">trusted" {
			var s string
			for user := range trustedusers {
				s = s + "<@" + user + ">, "
			}
			d.ChannelMessageSend(m.ChannelID, s)
		}

		//If message starts with >game, say the following text
		if strings.HasPrefix(m.Content, ">name ") || strings.HasPrefix(m.Content, ">nick ") {
			s := strings.SplitN(m.Content, " ", 2)
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
