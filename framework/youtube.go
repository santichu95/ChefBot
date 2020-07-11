package framework

import (
	"bytes"
	"fmt"
	"net/http"
	"net/url"
	"os/exec"
	"strings"

	"google.golang.org/api/googleapi/transport"
	"google.golang.org/api/youtube/v3"
)

// VideoInfo stores information needed to display in discord
// TODO(sandaluz) add the information about who added this song to the queue
type VideoInfo struct {
	URL   string
	Title string
	ID    string
}

// ParseYoutubeInput will parse the input and determine if it was a URL linking to a video or a search query,
// If it was a search query it will query youtube and return a link to the first result
func ParseYoutubeInput(input string) (string, error) {
	url, err := url.ParseRequestURI(input)

	if err == nil {
		return url.String(), nil
	}

	id := Search(input)

	return "youtu.be/" + id, nil
}

// Search ...
func Search(query string) string {
	developerKey := "AIzaSyDyWX6x3Ak9i0P7o1QPN0BKG0IB9PjZuk8"

	client := &http.Client{
		Transport: &transport.APIKey{Key: developerKey},
	}

	service, err := youtube.New(client)
	if err != nil {
		fmt.Println("Error creating new YouTube client: %v", err)
	}

	// Make the API call to YouTube.
	call := service.Search.List("id,snippet").
		Q(query).
		MaxResults(1)
	response, err := call.Do()
	if err != nil {
		fmt.Println("Error creating new YouTube client: %v", err)
	}

	//TODO get video title
	return response.Items[0].Id.VideoId
}

// DownloadVideo ...
func DownloadVideo(input string) (*VideoInfo, error) {
	cmd := exec.Command("youtube-dl", "-f", "140", "-o", "%(id)s", input)
	fmt.Println(cmd)
	var out bytes.Buffer
	cmd.Stdout = &out
	err := cmd.Run()

	videoInfo := new(VideoInfo)
	if err != nil {
		return videoInfo, err
	}

	videoInfo.URL = input
	videoInfo.Title = "Filler"
	videoInfo.ID = strings.Split(input, "=")[1]

	return videoInfo, nil
}
