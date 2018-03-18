package plugins

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"time"

	"github.com/buger/jsonparser"
	"github.com/bwmarrin/discordgo"
)

const urlBase = "https://www.googleapis.com/customsearch/v1?q="

var queries = make(map[string]Query)
var colour = 0xe75c5c

type result struct {
	number       int
	title        string
	snippet      string
	formattedUrl string
}

type Query struct {
	author    string
	number    int
	returned  int
	messageid string
	results   []result
}

func Googleinit(d *discordgo.Session) {
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
			queries[m.MessageID] = Query{author: queries[m.MessageID].author, messageid: queries[m.MessageID].messageid, number: i, results: queries[m.MessageID].results, returned: queries[m.MessageID].returned}

			embed := discordgo.MessageEmbed{
				Title:       queries[m.MessageID].results[queries[m.MessageID].number].title,
				Description: queries[m.MessageID].results[queries[m.MessageID].number].snippet,
				Color:       colour,
				Footer:      &discordgo.MessageEmbedFooter{Text: strconv.Itoa(queries[m.MessageID].results[queries[m.MessageID].number].number) + "/" + strconv.Itoa(queries[m.MessageID].returned)},
				URL:         queries[m.MessageID].results[queries[m.MessageID].number].formattedUrl,
			}
			_, err := s.ChannelMessageEditEmbed(m.ChannelID, m.MessageID, &embed)
			if err != nil {
				s.ChannelMessageSend(m.ChannelID, err.Error())
				return
			}
		}
	}
}

func Search(q string, n int, apiKey string, engineID string, d *discordgo.Session, m *discordgo.MessageCreate) {
	sender := m.Author.ID

	url := urlBase + url.QueryEscape(q) + "&cx=" + engineID + "&num=" + strconv.Itoa(n) + "&key=" + apiKey

	fmt.Println(url)
	j, err := http.Get(url)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	body, err := ioutil.ReadAll(j.Body)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	var results []result
	var returned int

	_, err = jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		title, err := jsonparser.GetString(value, "title")
		snippet, err := jsonparser.GetString(value, "snippet")
		url, err := jsonparser.GetString(value, "link")
		if err != nil {
			d.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		returned++

		result := result{number: returned, title: title, snippet: snippet, formattedUrl: url}
		if err != nil {
			panic(err.Error())
		}
		results = append(results, result)
	}, "items")

	if !(len(results) > 0) {
		d.ChannelMessageSend(m.ChannelID, "No results.")
		return
	}
	embed := discordgo.MessageEmbed{
		Title:       results[0].title,
		Description: results[0].snippet,
		Color:       colour,
		Footer:      &discordgo.MessageEmbedFooter{Text: strconv.Itoa(results[0].number) + "/" + strconv.Itoa(returned)},
		URL:         results[0].formattedUrl,
	}
	message, err := d.ChannelMessageSendEmbed(m.ChannelID, &embed)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	d.MessageReactionAdd(message.ChannelID, message.ID, "⬅")
	d.MessageReactionAdd(message.ChannelID, message.ID, "➡")

	var query = Query{author: sender, messageid: message.ID, number: 0, returned: returned, results: results}
	queries[message.ID] = query

	go removeQuery(d, message)
}

func removeQuery(d *discordgo.Session, message *discordgo.Message) {
	time.Sleep(2 * time.Minute)
	d.MessageReactionRemove(message.ChannelID, message.ID, "⬅", "@me")
	d.MessageReactionRemove(message.ChannelID, message.ID, "➡", "@me")

	delete(queries, message.ID)
}
