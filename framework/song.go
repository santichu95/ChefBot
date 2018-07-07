package framework

import (
	"os/exec"
)

// Song ...
type Song struct {
	Media    string
	Title    string
	Duration *string
	ID       string
}

// Ffmpeg ...
func (song Song) Ffmpeg() *exec.Cmd {
	return exec.Command("ffmpeg", "-i", song.Media, "-f", "s16le", "-ar", "48000", "-ac", "2", "pipe:1")
}

// NewSong ...
func NewSong(media, title, id string) *Song {
	song := new(Song)
	song.Media = media
	song.Title = title
	song.ID = id
	return song
}
