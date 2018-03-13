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

const urlBase = "https://www.googleapis.com/customsearch/v1?q="

//var queries = make(map[string]Query)

type result struct {
	number       int
	title        string `json:title`
	snippet      string `json:snippet`
	formattedUrl string `json:formattedUrl`
}

type Query struct {
	author string
	//id        int
	messageid string
}

func Search(q string, n int, apiKey string, engineID string, d *discordgo.Session, m *discordgo.MessageCreate) {
	sender := m.Author.ID

	url := urlBase + url.QueryEscape(q) + "&cx=" + engineID + "&num=" + strconv.Itoa(n) + "&key=" + apiKey
	fmt.Println(sender)
	fmt.Println(url)
	j, err := http.Get(url)
	if err != nil {
		panic(err.Error())
	}

	body, err := ioutil.ReadAll(j.Body)

	if err != nil {
		panic(err.Error())
	}
	var results []result
	var returned int

	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		title, err := jsonparser.GetString(value, "title")
		snippet, err := jsonparser.GetString(value, "snippet")
		formattedUrl, err := jsonparser.GetString(value, "formattedUrl")
		returned++

		result := result{number: returned, title: title, snippet: snippet, formattedUrl: formattedUrl}
		if err != nil {
			panic(err.Error())
		}
		results = append(results, result)
	}, "items")

	var returnString string

	returnString = strconv.Itoa(results[0].number) + "/" + strconv.Itoa(returned) + "\n" + results[0].title + "\n" + results[0].snippet + "\n<" + results[0].formattedUrl + ">"
	message, _ := d.ChannelMessageSend(m.ChannelID, returnString)

	d.MessageReactionAdd(message.ChannelID, message.ID, "⬅")
	d.MessageReactionAdd(message.ChannelID, message.ID, "➡")

	// TODO: REMOVE LISTENER
	//var query = Query{author: sender, messageid: m.Message.ID}
	//queries[m.Message.ID] = query

	var i = 0

	d.AddHandler(func(s *discordgo.Session, m *discordgo.MessageReactionAdd) {
		if m.UserID == sender {
			if m.Emoji.Name == "➡" {
				i++
				if i >= returned {
					i = 0
				}
			} else if m.Emoji.Name == "⬅" {
				if i == 0 {
					i = returned - 1
					fmt.Println("returned")
				} else {
					i--
				}
			}

			returnString = strconv.Itoa(results[i].number) + "/" + strconv.Itoa(returned) + "\n" + results[i].title + "\n" + results[i].snippet + "\n<" + results[i].formattedUrl + ">"
			d.ChannelMessageEdit(message.ChannelID, message.ID, returnString)
		}
	})
}
