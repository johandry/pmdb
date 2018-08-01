package crawler

import "github.com/johandry/gomdb/mdb"

// MovieCrawler defines what a crawler does, which is to collect movies from a source
type MovieCrawler interface {
	Crawl(rules ...string) mdb.DB
}
