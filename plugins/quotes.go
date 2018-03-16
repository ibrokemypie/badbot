package plugins

import (
	"bufio"
	"bytes"
	"io"
	"math/rand"
	"os"
)

func WriteQuote(key string, value string) {
	f, err := os.OpenFile("quotes/"+key, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		panic(err)
	}

	defer f.Close()

	if _, err = f.WriteString(value + "\n"); err != nil {
		panic(err)
	}
}
func lineCounter(r io.Reader) (int, error) {
	buf := make([]byte, 32*1024)
	count := 0
	lineSep := []byte{'\n'}

	for {
		c, err := r.Read(buf)
		count += bytes.Count(buf[:c], lineSep)

		switch {
		case err == io.EOF:
			return count, nil

		case err != nil:
			return count, err
		}
	}
}

func ReadQuote(key string) string {
	fileHandle, _ := os.Open("quotes/" + key)
	lines, err := lineCounter(fileHandle)
	fileHandle.Close()

	if err != nil {
		return "No quotes exist under that name."
	}

	fileHandle, _ = os.Open("quotes/" + key)
	scanner := bufio.NewScanner(fileHandle)

	line := 0
	text := ""

	chosenline := rand.Intn(lines) + 1

	for scanner.Scan() {
		line++
		if line == chosenline {
			text = scanner.Text()
		}
	}
	fileHandle.Close()
	return text
}
