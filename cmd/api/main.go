package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"golang.org/x/net/html"
)

func checkLink(link string, rootLink string) bool {
	return strings.Contains(link, rootLink) || (!strings.Contains(link, rootLink) && !strings.Contains(link, "://"))
}

func getPageLinks(resp *http.Response, rootLink string) []string {
	tokenizer := html.NewTokenizer(resp.Body)

	var pageLinks []string

	for {
		tt := tokenizer.Next()
		switch tt {
		case html.ErrorToken:
			return pageLinks
		case html.StartTagToken:
			tagName, _ := tokenizer.TagName()
			if tagName[0] == 'a' && len(tagName) == 1 {
				for {
					attr, val, moreArgs := tokenizer.TagAttr()
					if string(attr) == "href" && checkLink(string(val), rootLink) {
						pageLinks = append(pageLinks, string(val))
						break
					}
					if moreArgs == false {
						break
					}
				}
			}
		}
	}
}

type Link struct {
	Href   string
	Source string
	Depth  int
}

func main() {
	var urlFlag = flag.String("url", "https://atos.net", "url from where to parse links")
	var maxDepthFlag = flag.Int("depth", 3, "max depth to which to parse links")
	flag.Parse()

	pageLinksSet := make(map[string]bool) // only includes one copy of every path
	var linksArray []Link
	var linksToParse []Link // links that are waiting to be parsed by for paths
	linksToParse = append(linksToParse, Link{
		Href:   *urlFlag,
		Source: "",
		Depth:  1,
	})

	var rootLink string = *urlFlag
	rootLinkIndex := strings.Index(*urlFlag, "www.")
	if rootLinkIndex != -1 {
		rootLink = rootLink[rootLinkIndex+len("www."):]
	}

	for len(linksToParse) > 0 {
		if linksToParse[0].Depth >= *maxDepthFlag {
			linksToParse = linksToParse[1:]
			continue
		}
		resp, err := http.Get(linksToParse[0].Href)
		if err != nil {
			log.Print("Something went wrong with connection to:", linksToParse[0].Href)
			break
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			linksToParse = linksToParse[1:]
			continue
		}

		ctype := resp.Header.Get("Content-Type")
		if strings.HasPrefix(ctype, "text/html") {
			pageLinks := getPageLinks(resp, rootLink)
			for _, link := range pageLinks {
				fullLink := link
				if !strings.Contains(fullLink, rootLink) {
					if len(link) > 0 && link[0] == '/' {
						fullLink = *urlFlag + link
					} else {
						fullLink = *urlFlag + "/" + link
					}
				}

				// if link is not already present inside the main map, then append it to parsing
				if !pageLinksSet[fullLink] {
					newLink := Link{
						Href:   fullLink,
						Source: linksToParse[0].Href,
						Depth:  linksToParse[0].Depth + 1,
					}
					pageLinksSet[fullLink] = true
					if newLink.Depth < *maxDepthFlag {
						linksToParse = append(linksToParse, newLink)
					}
					linksArray = append(linksArray, newLink)
				}
			}
		}

		// dequeue first element of links
		fmt.Println(len(linksToParse))
		linksToParse = linksToParse[1:]
	}

	file, err := os.Create("links.txt")
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, link := range linksArray {
		fmt.Fprintf(file, " %s - %s %d\n", link.Href, link.Source, link.Depth)
		fmt.Printf(" %s - %s %d\n", link.Href, link.Source, link.Depth)
	}
}
