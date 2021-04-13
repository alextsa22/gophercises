package main

import (
	"errors"
	"flag"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/url"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/alextsa22/gophercises/13-quiet-hn/hn"
)

func main() {
	// parse flags
	var port, numStories int
	flag.IntVar(&port, "port", 3000, "the port to start the web server on")
	flag.IntVar(&numStories, "num_stories", 30, "the number of top stories to display")
	flag.Parse()

	tpl := template.Must(template.ParseFiles("index.html"))

	http.HandleFunc("/", handler(numStories, tpl))

	// Start the server
	fmt.Printf("server start on port %d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), nil))
}

func handler(numStories int, tpl *template.Template) http.HandlerFunc {
	sc := storyCache{
		numStories: numStories,
		duration:   6 * time.Second,
	}

	go func() {
		ticker := time.NewTicker(3 * time.Second)
		for {
			tmp := storyCache{
				numStories: numStories,
				duration:   6 * time.Second,
			}
			tmp.stories()

			sc.mutex.Lock()
			sc.cache = tmp.cache
			sc.expiration = tmp.expiration
			sc.mutex.Unlock()

			<-ticker.C
		}
	}()

	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		stories, err := sc.stories()
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		data := templateData{
			Stories: stories,
			Time:    time.Now().Sub(start),
		}
		err = tpl.Execute(w, data)
		if err != nil {
			http.Error(w, "Failed to process the template", http.StatusInternalServerError)
			return
		}
	}
}

type storyCache struct {
	numStories int
	cache      []item
	expiration time.Time
	duration   time.Duration
	mutex      sync.Mutex
}

func (c *storyCache) stories() ([]item, error) {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	if time.Now().Sub(c.expiration) < 0 {
		return c.cache, nil
	}

	stories, err := getTopStories(c.numStories)
	if err != nil {
		return nil, err
	}

	c.expiration = time.Now().Add(c.duration)
	c.cache = stories

	return c.cache, nil
}

func getTopStories(numStories int) ([]item, error) {
	var client hn.Client
	ids, err := client.TopItems()
	if err != nil {
		return nil, errors.New("failed to load top stories")
	}

	var stories []item
	at := 0
	for len(stories) < numStories {
		need := (numStories - len(stories)) * 5 / 4
		stories = append(stories, getStories(ids[at:at+need])...)
		at += need
	}

	return stories[:numStories], nil
}

func getStories(ids []int) []item {
	type result struct {
		idx  int
		item item
		err  error
	}
	resultCh := make(chan result, 100)

	for i := 0; i < len(ids); i++ {
		go func(idx, id int) {
			var client hn.Client
			hnItem, err := client.GetItem(id)
			if err != nil {
				resultCh <- result{idx: idx, err: err}
			}
			resultCh <- result{idx: idx, item: parseHNItem(hnItem)}
		}(i, ids[i])
	}

	var results []result
	for i := 0; i < len(ids); i++ {
		results = append(results, <-resultCh)
	}

	sort.Slice(results, func(i int, j int) bool {
		return results[i].idx < results[j].idx
	})

	var stories []item
	for _, res := range results {
		if res.err != nil {
			continue
		}

		if isStoryLink(res.item) {
			stories = append(stories, res.item)
			if len(stories) >= len(ids) {
				break
			}
		}
	}

	return stories
}

func isStoryLink(item item) bool {
	return item.Type == "story" && item.URL != ""
}

func parseHNItem(hnItem hn.Item) item {
	ret := item{Item: hnItem}
	url, err := url.Parse(ret.URL)
	if err == nil {
		ret.Host = strings.TrimPrefix(url.Hostname(), "www.")
	}

	return ret
}

// item is the same as the hn.Item, but adds the Host field
type item struct {
	hn.Item
	Host string
}

type templateData struct {
	Stories []item
	Time    time.Duration
}
