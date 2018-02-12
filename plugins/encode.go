package plugins

import (
	"bytes"
	"encoding/base64"
	"image/jpeg"
	"image/png"
	"io"
	"net/http"
	"os"
	"strings"
)

func EncodeImage(url string) string {
	s := strings.Split(url, "/")
	filename := s[len(s)-1]

	existingImageFile, _ := os.Create(filename)
	defer existingImageFile.Close()

	resp, _ := http.Get(url)
	defer resp.Body.Close()

	io.Copy(existingImageFile, resp.Body)

	existingImageFile.Seek(0, 0)

	extension := filetype(filename)
	var buff bytes.Buffer

	if strings.EqualFold(extension, "jpg ") || strings.EqualFold(extension, "jpeg") {
		loadedImage, _ := jpeg.Decode(existingImageFile)
		jpeg.Encode(&buff, loadedImage, nil)
	} else if strings.EqualFold(extension, "png") {
		loadedImage, _ := png.Decode(existingImageFile)
		png.Encode(&buff, loadedImage)
	}

	encodedString := base64.StdEncoding.EncodeToString(buff.Bytes())
	os.Remove(filename)
	// You can embed it in an html doc with this string
	baseimg := "data:image/" + extension + ";base64," + encodedString
	return baseimg
}

func filetype(name string) string {
	s := strings.Split(name, ".")
	extension := s[len(s)-1]

	return extension
}
