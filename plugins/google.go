package plugins

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/buger/jsonparser"
	"github.com/bwmarrin/discordgo"
)

func SearchGoogle(q string, n int, apiKey string, engineID string, d *discordgo.Session, m *discordgo.MessageCreate) {
	var colour = 0xe75c5c
	var urlBase = "https://www.googleapis.com/customsearch/v1?q="

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

	var query = Query{author: sender, messageid: message.ID, number: 0, returned: returned, results: results, colour: colour, qtype: "google"}
	queries[message.ID] = query

	go removeQuery(d, message)
}
