package mdb

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v2"
)

// DB is a collection of movies
type DB struct {
	Movies map[string]*Movie `json:"movies" yaml:"movies"`
	Stats  Stats             `json:"stats" yaml:"stats"`
}

// Stats contain statistics about the movie colection
type Stats struct {
	Movies int `json:"movies" yaml:"movies"`
	Files  int `json:"files" yaml:"files"`
	Dubbed int `json:"dubbed" yaml:"dubbed"`
	HD     int `json:"hd" yaml:"hd"`
}

// New creates a new and empty movies database
func New() *DB {
	movies := make(map[string]*Movie, 0)

	return &DB{
		Movies: movies,
		Stats:  Stats{},
	}
}

// Add adds a new movie generating a unique id
func (db *DB) Add(title string, location string, size int64, mtime time.Time) {
	mov := NewMovie(title, location, size, mtime)
	id := base64.StdEncoding.EncodeToString([]byte(mov.Title))
	if _, ok := db.Movies[id]; ok {
		if db.Movies[id].Title != mov.Title {
			panic(fmt.Sprintf("base64 of title %q and %q are the same:\nIn DB: %v\nNew: %s", db.Movies[id].Title, title, db.Movies[id], location))
		}
		if !db.Movies[id].HD && mov.HD {
			db.Stats.HD++
		}
		if !db.Movies[id].Dubbed && mov.Dubbed {
			db.Stats.Dubbed++
		}
		db.Movies[id].HD = db.Movies[id].HD || mov.HD
		db.Movies[id].Dubbed = db.Movies[id].Dubbed || mov.Dubbed
		db.Movies[id].AddFile(location, size, mtime)
		db.Stats.Files++
		return
	}
	db.Movies[id] = mov
	db.Stats.Movies++
	db.Stats.Files++
	if db.Movies[id].HD {
		db.Stats.HD++
	}
	if db.Movies[id].Dubbed {
		db.Stats.Dubbed++
	}
}

// Marshal returns the movies collection in the requested format: yaml, json or
// pretty json
func (db *DB) Marshal(format string, pretty ...bool) ([]byte, error) {
	switch format {
	case "json", "js":
		if len(pretty) > 0 && pretty[0] {
			return json.MarshalIndent(db, "", "  ")
		}
		return json.Marshal(db)
	case "yaml", "yml":
		return yaml.Marshal(db)
	default:
		return nil, fmt.Errorf("unknown format %q", format)
	}
}

// WriteFile writes a file with the movies in the requested format: yaml, json or
// pretty json
func (db *DB) WriteFile(filename string) error {
	format := strings.TrimPrefix(filepath.Ext(filename), ".")
	pretty := (format == "json") || (format == "js")
	content, err := db.Marshal(format, pretty)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(filename, content, 0644); err != nil {
		return err
	}
	return nil
}

// MarkBakup mark all the movie files that are backup. If there is a backup
// movie not in the source location, will be mark as 'move_to_source' and return
// the list of file locations
func (db *DB) MarkBakup(sourceLoc ...string) ([]string, []string) {
	moveToSource := []string{}
	backupMovies := []string{}

	for _, mov := range db.Movies {
		// Just one file, this is OK, next.
		if len(mov.Files) <= 1 {
			continue
		}

		// fmt.Printf("[DEBUG] checking movie files of %s\n", mov.Title)
		for id, file := range mov.Files {
			for _, loc := range sourceLoc {
				if strings.Contains(file.Location, loc) {
					mov.Files[id].Stat = "source"
				} else {
					// fmt.Printf("[DEBUG]   >> found backup: %s\n", file.Location)
					backupMovies = append(backupMovies, file.Location)
					mov.Files[id].Stat = "backup"
				}
			}
		}

		for id, file := range mov.Files {
			if file.Stat == "backup" && !mov.fileInStat(file, "source") {
				// fmt.Printf("[DEBUG]   >>>> found backup that should be source: %s\n", file.Location)
				moveToSource = append(moveToSource, file.Location)
				mov.Files[id].Stat = "move_to_source"
			}
		}
	}

	return moveToSource, backupMovies
}
