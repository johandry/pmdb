package main

import (
	"fmt"
	"io/ioutil"
	"sort"
	"strings"

	"github.com/johandry/gomdb/crawler"
	"github.com/sirupsen/logrus"
)

const (
	moviesFilename           = "mdb.json"
	movieFilesToMoveFilename = "file_to_move.txt"
	backupMoviesFilename     = "backup_movies.txt"
	sourceLocation           = "/Volumes/SMC/Shared iTunes/iTunes Media/Movies"
)

func main() {
	var fsCrawler = crawler.NewFSCrawler(
		sourceLocation,
		"/Volumes/SMC/_To Organize/Movies",
	)
	// fsCrawler.SetLogLevel(logrus.DebugLevel)
	fsCrawler.SetLogLevel(logrus.WarnLevel)

	movies, _ := fsCrawler.Crawl(
		"m4v",
		"mp4",
	)

	moveMov, bkpMov := movies.MarkBakup(sourceLocation)
	if len(moveMov) > 0 {
		total := len(moveMov)
		sort.Strings(moveMov)
		content := strings.Join(moveMov, "\n")
		if err := ioutil.WriteFile(movieFilesToMoveFilename, []byte(content), 0644); err != nil {
			panic(err)
		}
		fmt.Printf("A total of %d movie files to move to %q have been saved in file %q\n", total, sourceLocation, movieFilesToMoveFilename)
	}

	if len(bkpMov) > 0 {
		total := len(bkpMov)
		sort.Strings(bkpMov)
		content := strings.Join(bkpMov, "\n")
		if err := ioutil.WriteFile(backupMoviesFilename, []byte(content), 0644); err != nil {
			panic(err)
		}
		fmt.Printf("A total of %d backup movie files have been saved in file %q\n", total, backupMoviesFilename)
	}

	movies.WriteFile(moviesFilename)
	fmt.Printf("Movies have been saved in file %q\n", moviesFilename)
	fmt.Printf("\tMovies: %d\n\tMovie Files: %d\n\tHD Movies: %d\n\tDubbed Movies: %d\n", movies.Stats.Movies, movies.Stats.Files, movies.Stats.HD, movies.Stats.Dubbed)
}
