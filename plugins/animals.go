package plugins

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

func Woof() string {
	url := "https://random.dog/"

	resp, err := http.Get(url + "woof")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)
	return url + string(body)
}

func Meow() string {
	type files struct {
		File string `json:"file"`
	}

	url := "https://random.cat/"

	resp, err := http.Get(url + "meow")
	if err != nil {
		fmt.Println(err)
		return ""
	}
	defer resp.Body.Close()

	body, _ := ioutil.ReadAll(resp.Body)

	filej := files{}
	json.Unmarshal(body, &filej)

	return filej.File
}
