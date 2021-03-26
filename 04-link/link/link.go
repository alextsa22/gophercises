package link

import (
	"golang.org/x/net/html"
	"io"
	"strings"
)

type Link struct {
	Href string
	Text string
}

func Parse(r io.Reader) ([]Link, error) {
	root, err := html.Parse(r)
	if err != nil {
		return nil, err
	}

	return getLinks(root), nil
}

func getLinks(n *html.Node) []Link {
	var links []Link
	if n.Type == html.ElementNode && n.Data == "a" {
		for _, a := range n.Attr {
			if a.Key == "href" {
				text := getText(n)
				links = append(links, Link{a.Val, text})
			}
		}
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		childLinks := getLinks(c)
		links = append(links, childLinks...)
	}

	return links
}

func getText(n *html.Node) string {
	var text string
	if n.Type != html.ElementNode && n.Data != "a" && n.Type != html.CommentNode {
		text = n.Data
	}

	for c := n.FirstChild; c != nil; c = c.NextSibling {
		text += getText(c)
	}

	return strings.Trim(text, "\n")
}
