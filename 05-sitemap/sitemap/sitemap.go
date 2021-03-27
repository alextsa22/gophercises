package sitemap

import (
	"github.com/alextsa22/gophercises/04-link/link"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type Sitemap struct {
	URL      *url.URL
	MaxDepth int
}

func NewSitemap(URL *url.URL, maxDepth int) *Sitemap {
	return &Sitemap{URL: URL, MaxDepth: maxDepth}
}

func (s *Sitemap) Build() []string {
	seen := make(map[string]struct{})
	q := make(map[string]struct{})
	nq := map[string]struct{}{
		s.URL.String(): {},
	}

	for i := 0; i <= s.MaxDepth; i++ {
		q, nq = nq, make(map[string]struct{})
		if len(q) == 0 {
			break
		}

		for url, _ := range q {
			if _, ok := seen[url]; ok {
				continue
			}

			seen[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for url, _ := range seen {
		ret = append(ret, url)
	}

	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		return []string{}
	}
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	base := baseUrl.String()

	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)

	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}

	return ret
}

func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
