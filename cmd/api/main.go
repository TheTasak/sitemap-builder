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

type Link struct {
	Href   string
	Source string
	Depth  int
}

func main() {
	var urlFlag = flag.String("url", "https://www.interia.pl", "url from where to parse links")
	var maxDepthFlag = flag.Int("depth", 3, "max depth to which to parse links")
	var fileFlag = flag.String("file", "result.txt", "path to file where to store results of program execution")
	var cmdFlag = flag.Bool("showCmd", false, "show results in command line")
	var domainFlag = flag.Bool("sameDomain", true, "include in searching only domain specified by the url flag")
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
			linksToParse = linksToParse[1:]
			continue
		}

		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			linksToParse = linksToParse[1:]
			continue
		}

		ctype := resp.Header.Get("Content-Type")
		if strings.HasPrefix(ctype, "text/html") {
			pageLinks := getPageLinks(resp, rootLink, *domainFlag)
			for _, link := range pageLinks {
				fullLink := makeLink(link, rootLink, *domainFlag, *urlFlag)

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

	file, err := os.Create(*fileFlag)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()

	for _, link := range linksArray {
		fmt.Fprintf(file, " %s - %s %d\n", link.Href, link.Source, link.Depth)
		if *cmdFlag {
			fmt.Printf(" %s - %s %d\n", link.Href, link.Source, link.Depth)
		}
	}
}

func makeLink(link string, rootLink string, onlyRootLink bool, urlFlag string) string {
	if !strings.Contains(link, rootLink) {
		if !strings.Contains(link, "://") || onlyRootLink {
			if len(link) > 0 && link[0] == '/' {
				return urlFlag + link
			} else {
				return urlFlag + "/" + link
			}
		}
	}
	return link
}

func checkLink(link string, rootLink string, onlyRootLink bool) bool {
	if onlyRootLink {
		return strings.Contains(link, rootLink) || (!strings.Contains(link, rootLink) && !strings.Contains(link, "://"))
	}
	return true
}

func getPageLinks(resp *http.Response, rootLink string, onlyRootLink bool) []string {
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
					if string(attr) == "href" && checkLink(string(val), rootLink, onlyRootLink) {
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
