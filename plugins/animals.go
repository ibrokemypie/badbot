package plugins

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

func Woof() string {
	url := "https://random.dog/"

	resp, _ := http.Get(url + "woof")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return url + string(body)
}

func Meow() string {
	type files struct {
		File string `json:"file"`
	}

	url := "https://random.cat/"

	resp, _ := http.Get(url + "meow")
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	filej := files{}
	json.Unmarshal(body, &filej)

	return filej.File
}
