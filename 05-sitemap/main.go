package main

import (
	"encoding/xml"
	"flag"
	"fmt"
	"github.com/alextsa22/gophercises/05-sitemap/sitemap"
	"log"
	"net/url"
	"os"
)

const xmlns = "http://www.sitemaps.org/schemas/sitemap/0.9"

var (
	urlFlag = flag.String("url", "https://gophercises.com/", "the url that you want to build a sitemap for")
	depth   = flag.Int("depth", 10, "the maximum number of links deep to traverse")
)

type loc struct {
	Value string `xml:"loc"`
}

type urlset struct {
	Urls  []loc  `xml:"url"`
	Xmlns string `xml:"xmlns,attr"`
}

func main() {
	flag.Parse()

	urlObj, err := url.Parse(*urlFlag)
	if err != nil {
		log.Fatal(err)
	}

	pages := sitemap.NewSitemap(urlObj, 10).Build()

	toXml := urlset{
		Xmlns: xmlns,
	}

	for _, page := range pages {
		toXml.Urls = append(toXml.Urls, loc{page})
	}

	fmt.Print(xml.Header)

	enc := xml.NewEncoder(os.Stdout)
	enc.Indent("", "  ")
	if err := enc.Encode(toXml); err != nil {
		panic(err)
	}

	fmt.Println()
}
