package plugins

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"

	"github.com/buger/jsonparser"
)

const urlBase = "https://www.googleapis.com/customsearch/v1?q="

type result struct {
	title        string `json:title`
	snippet      string `json:snippet`
	formattedUrl string `json:formattedUrl`
}

func Search(q string, n int, apiKey string, engineID string) string {

	url := urlBase + url.QueryEscape(q) + "&cx=" + engineID + "&num=" + strconv.Itoa(n) + "&key=" + apiKey
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

	jsonparser.ArrayEach(body, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {
		title, err := jsonparser.GetString(value, "title")
		snippet, err := jsonparser.GetString(value, "snippet")
		formattedUrl, err := jsonparser.GetString(value, "formattedUrl")

		result := result{title: title, snippet: snippet, formattedUrl: formattedUrl}
		if err != nil {
			panic(err.Error())
		}
		results = append(results, result)
	}, "items")

	for i, result := range results {
		fmt.Println(i)
		fmt.Println(result.title)
		fmt.Println(result.snippet)
		fmt.Println(result.formattedUrl)
	}

	var returnString string

	returnString = results[0].title + "\n" + results[0].snippet + "\n<" + results[0].formattedUrl + ">"
	return returnString
}
