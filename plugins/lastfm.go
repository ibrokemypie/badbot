package plugins

import (
	"io/ioutil"
	"net/http"

	"github.com/buger/jsonparser"
	"github.com/bwmarrin/discordgo"
)

func lastfmPlaying(lastfmapi string, lastfmuser string) *discordgo.MessageEmbed {
	urlBase := "https://ws.audioscrobbler.com/2.0/?method=user.getrecenttracks&format=json"
	user := "&user=" + lastfmuser
	apiKey := "&api_key=" + lastfmapi
	options := "&limit=1&extended=0&nowplaying=true"

	url := urlBase + apiKey + options + user

	j, err := http.Get(url)

	if err != nil {
		return &discordgo.MessageEmbed{Title: err.Error()}
	}

	body, err := ioutil.ReadAll(j.Body)
	if err != nil {
		return &discordgo.MessageEmbed{Title: err.Error()}
	}

	artist, _ := jsonparser.GetString(body, "recenttracks", "track", "[0]", "artist", "#text")
	album, _ := jsonparser.GetString(body, "recenttracks", "track", "[0]", "album", "#text")
	track, _ := jsonparser.GetString(body, "recenttracks", "track", "[0]", "name")
	songurl, _ := jsonparser.GetString(body, "recenttracks", "track", "[0]", "url")

	embed := discordgo.MessageEmbed{
		Title: track + " - " + album + " [" + artist + "]",
		//			Color:       colour,
		URL: songurl,
	}

	return &embed
}
