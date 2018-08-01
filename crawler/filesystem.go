package crawler

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/sirupsen/logrus"

	"github.com/johandry/gomdb/mdb"
	"github.com/johandry/log"
)

// FilesystemCrawler implements the Crawler interface to search for movies in
// the FileSystem
type FilesystemCrawler struct {
	scheme string
	dir    []string
	total  int
	logger *log.Logger
}

// NewFSCrawler returns a new FilesystemCrawler
func NewFSCrawler(dirPath ...string) *FilesystemCrawler {
	logger := log.NewDefault()
	logger.SetPrefix("Crawler")

	fsCrawler := FilesystemCrawler{
		scheme: "file",
		dir:    dirPath,
		logger: logger,
	}
	return &fsCrawler
}

// SetLogLevel change the log level of the crawler logger
func (fsc *FilesystemCrawler) SetLogLevel(level logrus.Level) {
	fsc.logger.SetLevel(level)
}

// Crawl search for movies on the given directories
func (fsc *FilesystemCrawler) Crawl(rules ...string) (*mdb.DB, error) {
	db := mdb.New()

	for _, dir := range fsc.dir {
		if err := filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
			if info == nil || info.IsDir() {
				fsc.logger.Debugf("getting into directory %s", path)
				return nil
			}
			for _, rule := range rules {
				if filepath.Ext(path) == "."+rule {
					filename := filepath.Base(path)
					title := strings.TrimSuffix(filename, "."+rule)
					fsc.logger.Debugf("adding movie %q", title)
					db.Add(title, fmt.Sprintf("%s://%s", fsc.scheme, path), info.Size(), info.ModTime())
					return nil
				}
			}
			fsc.logger.Warnf("skipping file %s", path)
			return nil
		}); err != nil {
			return db, err
		}
	}

	return db, nil
}
