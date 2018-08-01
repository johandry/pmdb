package mdb

import (
	"fmt"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

const (
	dubbedTxt = " (Doblada)"
	hdTxt1     = " (1080p HD)"
	hdTxt2     = " (HD)"
	colon      = "_"
)

// Movie store the movie information
type Movie struct {
	ID      string      `json:"id" yaml:"id"`
	Title   string      `json:"title" yaml:"title"`
	Dubbed bool        `json:"dubbed" yaml:"dubbed"`
	HD      bool        `json:"hd" yaml:"hd"`
	Count   int         `json:"count" yaml:"count"`
	Files   []MovieFile `json:"files" yaml:"files"`
}

// MovieFile store the movie file information
type MovieFile struct {
	Size     int64     `json:"size" yaml:"size"`
	Location string    `json:"location" yaml:"location"`
	ModTime  time.Time `json:"mtime" yaml:"mtime"`
	Stat     string    `json:"status,omitempty" yaml:"status,omitempty"`
}

// NewMovie creates a new Movie
func NewMovie(title string, location string, size int64, mtime time.Time) *Movie {
	uuid, _ := uuid.NewV4()
	id := uuid.String()

	var dubbed bool
	if strings.Contains(title, dubbedTxt) {
		dubbed = true
		title = strings.Replace(title, dubbedTxt, "", 1)
	}

	var hd bool
	if strings.Contains(title, hdTxt1) {
		hd = true
		title = strings.Replace(title, hdTxt1, "", 1)
	}

	if strings.Contains(title, hdTxt2) {
		hd = true
		title = strings.Replace(title, hdTxt2, "", 1)
	}

	title = strings.Replace(title, colon, ":", -1)

	file := MovieFile{
		Size:     size,
		Location: location,
		ModTime:  mtime,
	}

	mov := Movie{
		ID:      id,
		Title:   title,
		Dubbed: dubbed,
		HD:      hd,
		Count:   1,
		Files:   []MovieFile{file},
	}

	return &mov
}

// AddFile append a new file to the movie
func (m *Movie) AddFile(location string, size int64, mtime time.Time) {
	file := MovieFile{
		Size:     size,
		Location: location,
		ModTime:  mtime,
	}

	m.Files = append(m.Files, file)
	m.Count = len(m.Files)
}

func (f MovieFile) filename() (string, error) {
	// This may not be a good or generic idea, not all the deparators could be "/"
	items := strings.Split(f.Location, "/")
	if len(items) == 0 || len(items) == 1 {
		return "", fmt.Errorf("filename not found for %s: %v", f.Location, items)
	}
	return items[len(items)-1], nil

}

func (f MovieFile) equalTo(file MovieFile) bool {
	if f.Size != file.Size {
		return false
	}

	fname, err := f.filename()
	if err != nil {
		return false
	}
	filename, err := file.filename()
	if err != nil {
		return false
	}

	return fname == filename
}

func (m *Movie) fileInStat(file MovieFile, stat string) bool {
	for _, f := range m.Files {
		if f.Location == file.Location {
			// Exactly same file, go next
			continue
		}
		if f.Stat == stat && f.equalTo(file) {
			return true
		}
	}

	return false
}
