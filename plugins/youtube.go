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

func SearchYoutube(q string, n int, apiKey string, d *discordgo.Session, m *discordgo.MessageCreate) {
	var colour = 0xe75c5c
	var urlBase = "https://www.googleapis.com/youtube/v3/search?q="

	sender := m.Author.ID

	url := urlBase + url.QueryEscape(q) + "&maxResults=" + strconv.Itoa(n) + "&key=" + apiKey + "&part=snippet&type=video"

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
		title, err := jsonparser.GetString(value, "snippet", "title")
		id, err := jsonparser.GetString(value, "id", "videoId")
		url := "https://www.youtube.com/watch?v=" + id
		thumbnail, err := jsonparser.GetString(value, "snippet", "thumbnails", "high", "url")
		if err != nil {
			d.ChannelMessageSend(m.ChannelID, err.Error())
			return
		}

		returned++

		result := result{number: returned, title: title, formattedUrl: url, thumbnail: thumbnail}
		if err != nil {
			panic(err.Error())
		}
		results = append(results, result)
	}, "items")

	if !(len(results) > 0) {
		d.ChannelMessageSend(m.ChannelID, "No results.")
		return
	}
	mtext := results[0].title + "\n" + results[0].formattedUrl + "\n ``" + strconv.Itoa(results[0].number) + "/" + strconv.Itoa(returned) + "``"
	message, err := d.ChannelMessageSend(m.ChannelID, mtext)
	if err != nil {
		d.ChannelMessageSend(m.ChannelID, err.Error())
		return
	}

	d.MessageReactionAdd(message.ChannelID, message.ID, "⬅")
	d.MessageReactionAdd(message.ChannelID, message.ID, "➡")

	var query = Query{author: sender, messageid: message.ID, number: 0, returned: returned, results: results, colour: colour, qtype: "youtube"}
	queries[message.ID] = query

	go removeQuery(d, message)
}
