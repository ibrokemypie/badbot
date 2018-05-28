package plugins

import (
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var queries = make(map[string]Query)

type result struct {
	number       int
	title        string
	snippet      string
	formattedUrl string
	thumbnail    string
	image        string
}

type Query struct {
	author    string
	number    int
	returned  int
	messageid string
	qtype     string
	colour    int
	results   []result
}

func ButtonInit(d *discordgo.Session) {
	d.AddHandler(buttons)
}

func buttons(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
	if queries[m.MessageID].messageid != "" {
		if m.UserID == queries[m.MessageID].author {
			i := queries[m.MessageID].number

			if m.Emoji.Name == "➡" {
				i++
				if i >= queries[m.MessageID].returned {
					i = 0
				}
			} else if m.Emoji.Name == "⬅" {
				if i == 0 {
					i = queries[m.MessageID].returned - 1
				} else {
					i--
				}
			}

			if queries[m.MessageID].qtype == "google" {
				queries[m.MessageID] = Query{author: queries[m.MessageID].author, messageid: queries[m.MessageID].messageid, number: i, results: queries[m.MessageID].results, returned: queries[m.MessageID].returned, colour: queries[m.MessageID].colour, qtype: queries[m.MessageID].qtype}

				embed := discordgo.MessageEmbed{
					Title:       queries[m.MessageID].results[queries[m.MessageID].number].title,
					Description: queries[m.MessageID].results[queries[m.MessageID].number].snippet,
					Color:       queries[m.MessageID].colour,
					Footer:      &discordgo.MessageEmbedFooter{Text: strconv.Itoa(queries[m.MessageID].results[queries[m.MessageID].number].number) + "/" + strconv.Itoa(queries[m.MessageID].returned)},
					URL:         queries[m.MessageID].results[queries[m.MessageID].number].formattedUrl,
				}
				_, err := s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, &embed)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
					return
				}
				return
			} else if queries[m.MessageID].qtype == "youtube" {
				queries[m.MessageID] = Query{author: queries[m.MessageID].author, messageid: queries[m.MessageID].messageid, number: i, results: queries[m.MessageID].results, returned: queries[m.MessageID].returned, colour: queries[m.MessageID].colour, qtype: queries[m.MessageID].qtype}

				mtext := queries[m.MessageID].results[queries[m.MessageID].number].title + "\n" + queries[m.MessageID].results[queries[m.MessageID].number].formattedUrl + "\n ``" + strconv.Itoa(queries[m.MessageID].results[queries[m.MessageID].number].number) + "/" + strconv.Itoa(queries[m.MessageID].returned) + "``"
				_, err := s.ChannelMessageEdit(m.ChannelID, m.MessageID, mtext)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
					return
				}
				return
			} else if queries[m.MessageID].qtype == "image" {
				queries[m.MessageID] = Query{author: queries[m.MessageID].author, messageid: queries[m.MessageID].messageid, number: i, results: queries[m.MessageID].results, returned: queries[m.MessageID].returned, colour: queries[m.MessageID].colour, qtype: queries[m.MessageID].qtype}

				embed := discordgo.MessageEmbed{
					Title:  queries[m.MessageID].results[queries[m.MessageID].number].title,
					Image:  &discordgo.MessageEmbedImage{URL: queries[m.MessageID].results[queries[m.MessageID].number].image},
					Color:  queries[m.MessageID].colour,
					Footer: &discordgo.MessageEmbedFooter{Text: strconv.Itoa(queries[m.MessageID].results[queries[m.MessageID].number].number) + "/" + strconv.Itoa(queries[m.MessageID].returned)},
					URL:    queries[m.MessageID].results[queries[m.MessageID].number].formattedUrl,
				}
				_, err := s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, &embed)
				if err != nil {
					s.ChannelMessageSend(m.ChannelID, err.Error())
					return
				}
				return
			}
		}
	}
}

func removeQuery(d *discordgo.Session, message *discordgo.Message) {
	time.Sleep(2 * time.Minute)
	d.MessageReactionRemove(message.ChannelID, message.ID, "⬅", "@me")
	d.MessageReactionRemove(message.ChannelID, message.ID, "➡", "@me")

	delete(queries, message.ID)
}
