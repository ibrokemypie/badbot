package plugins

import (
	"encoding/gob"
	"math/rand"
	"os"
	"strconv"
)

var numquotes int
var quoteids = make(map[int]string)
var nameid = make(map[string][]int)
var idquote = make(map[int]string)

func WriteQuote(key string, value string) {
	var i int
	numquotes++
	i = numquotes
	quoteids[i] = key

	nameid[key] = append(nameid[key], i)
	idquote[i] = value
}

func ReadQuote(key string) (string, int) {
	qlen := len(nameid[key])

	if qlen == 0 {
		return "", 0
	}
	chosenq := rand.Intn(qlen)
	id := nameid[key][chosenq]

	text := idquote[id]
	return text, id
}

func ReadQuoteID(id string) (string, int, string) {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return "Not a Number.", 0, ""
	}
	if idquote[rid] != "" {
		text := idquote[rid]
		qname := quoteids[rid]
		return text, rid, qname
	}
	return "No quotes with that ID found.", 0, ""
}

func RemoveQuote(id string) string {
	rid, err := strconv.Atoi(id)
	if err != nil {
		return "Not a Number."
	}

	if idquote[rid] != "" {
		qname := quoteids[rid]

		delete(quoteids, rid)
		delete(nameid, qname)
		delete(idquote, rid)
		return "Quote removed."
	}
	return "No quotes with that ID found."
}

func SaveData() {
	encodeFile, err := os.Create("quotes/numquotes.gob")
	if err != nil {
		panic(err)
	}
	encoder := gob.NewEncoder(encodeFile)
	if err := encoder.Encode(numquotes); err != nil {
		panic(err)
	}
	encodeFile.Close()

	encodeFile, err = os.Create("quotes/quoteids.gob")
	if err != nil {
		panic(err)
	}
	encoder = gob.NewEncoder(encodeFile)
	if err = encoder.Encode(quoteids); err != nil {
		panic(err)
	}
	encodeFile.Close()

	encodeFile, err = os.Create("quotes/nameid.gob")
	if err != nil {
		panic(err)
	}
	encoder = gob.NewEncoder(encodeFile)
	if err = encoder.Encode(nameid); err != nil {
		panic(err)
	}
	encodeFile.Close()

	encodeFile, err = os.Create("quotes/idquote.gob")
	if err != nil {
		panic(err)
	}
	encoder = gob.NewEncoder(encodeFile)
	if err = encoder.Encode(idquote); err != nil {
		panic(err)
	}
	encodeFile.Close()
}

func LoadData() {
	decodeFile, err := os.Open("quotes/numquotes.gob")
	if err != nil {
		panic(err)
	}

	decoder := gob.NewDecoder(decodeFile)
	decoder.Decode(&numquotes)
	decodeFile.Close()

	decodeFile, err = os.Open("quotes/quoteids.gob")
	if err != nil {
		panic(err)
	}

	decoder = gob.NewDecoder(decodeFile)
	decoder.Decode(&quoteids)
	decodeFile.Close()

	decodeFile, err = os.Open("quotes/nameid.gob")
	if err != nil {
		panic(err)
	}

	decoder = gob.NewDecoder(decodeFile)
	decoder.Decode(&nameid)
	decodeFile.Close()

	decodeFile, err = os.Open("quotes/idquote.gob")
	if err != nil {
		panic(err)
	}
	decoder = gob.NewDecoder(decodeFile)
	decoder.Decode(&idquote)
	decodeFile.Close()

}
