package framework

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
)

// VideoInfo stores information needed to display in discord
// TODO(sandaluz) add the information about who added this song to the queue
type VideoInfo struct {
	WebpageURL string
	Title      string
	ID         string
}

// DownloadYoutubeVideo will parse the input and grab the appropriate video. The input can either be a search term or a direct link to a single video or a playlist.
func DownloadYoutubeVideo(input string) (*[]VideoInfo, error) {
	if strings.Contains(strings.ToLower(input), "playlist") {
		return DownloadMultipleYoutubeVideo(input)
	}
	return DownloadSingleYoutubeVideo(input)
}

// DownloadMultipleYoutubeVideo will download an entire youtube playlist video a url to the playlist is required.
func DownloadMultipleYoutubeVideo(input string) (*[]VideoInfo, error) {
	// TODO(sandaluz) when you queue a large playlist it waitings for every video to be downloaded before it start playing music
	cmd := exec.Command("youtube-dl", "-i", "--print-json", "-f", "140", "-o", "audio_cache/%(id)s", input)
	fmt.Println(cmd)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	jsonList := strings.Split(strings.TrimSpace(out.String()), "\n")
	videoInfoList := make([]VideoInfo, len(jsonList))
	for i, videoJSON := range jsonList {
		fmt.Println(i, videoJSON)
		videoInfo := new(VideoInfo)
		err = json.Unmarshal([]byte(videoJSON), &videoInfo)
		if err != nil {
			return nil, err
		}
		videoInfoList[i] = *videoInfo
	}

	return &videoInfoList, nil
}

// DownloadSingleYoutubeVideo will download a single youtube video when a search term or a direct link is provided.
func DownloadSingleYoutubeVideo(input string) (*[]VideoInfo, error) {
	cmd := exec.Command("youtube-dl", "-i", "--print-json", "-f", "140", "-o", "audio_cache/%(id)s", "ytsearch:\""+input+"\"")
	fmt.Println(cmd)
	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		return nil, err
	}

	videoInfo := new(VideoInfo)
	err = json.Unmarshal(out.Bytes(), &videoInfo)
	if err != nil {
		return nil, err
	}

	return &[]VideoInfo{*videoInfo}, nil
}
