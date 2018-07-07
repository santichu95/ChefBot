package framework

import (
	"bytes"
	"encoding/json"
	"net/url"
	"os/exec"
)

type videoResponse struct {
	Formats []struct {
		URL string `json:"url"`
	} `json:"formats"`
	Title string `json:"title"`
}

// AudioItem  stores information about songs that are queued to play.
type AudioItem struct {
	vr videoResponse
}

// ParseYoutubeInput will parse the input and determine if it was a URL linking to a video or a search query,
// If it was a search query it will query youtube and return a link to the first result
func ParseYoutubeInput(input string) (string, error) {
	url, err := url.ParseRequestURI(input)

	if err == nil {
		return url.String(), nil
	}

	//TODO(sandaluz) search for the yt url based on keyword searches.

	//Temp video url
	return "youtu.be/dQw4w9WgXcQ", nil
}

//GetVideoInfo will get a download url, the title of the video, and an error if onee arises.
func GetVideoInfo(input string) (string, string, error) {
	cmd := exec.Command("youtube-dl", "--skip-download", "--print-json", "--flat-playlist", input)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	if err != nil {
		return "", "", err
	}

	str := out.String()

	var resp videoResponse
	err = json.Unmarshal([]byte(str), &resp)
	if err != nil {
		return "", "", err
	}

	return resp.Formats[0].URL, resp.Title, nil

}
