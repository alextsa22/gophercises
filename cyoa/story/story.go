package story

import (
	"bytes"
	"encoding/json"
	"os"
)

type Stories map[string]Story

type Story struct {
	Title   string
	Story   []string
	Options []Option `json:"options"`
}

type Option struct {
	Text string `json:"text"`
	Arc  string `json:"arc"'`
}

func NewStories(filename string) (Stories, error) {
	bytes, err := getFileBytes(filename)
	if err != nil {
		return nil, err
	}

	var stories Stories
	if err := json.Unmarshal(bytes, &stories); err != nil {
		return nil, err
	}

	return stories, nil
}

func getFileBytes(fileName string) ([]byte, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(f)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
