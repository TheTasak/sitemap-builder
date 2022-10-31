# Sitemap Builder
Command line tool for recursive links scraping of a website source
## How to use
Git clone this repository to your local machine, then use 
```sh
go build ./cmd/api
```
and execute the built source code
or just (with default flag settings)
```sh
go run ./cmd/api
```
## Options
The app doesn't require any flags, but there's couple of optional ones
* url=String - url of the website to scrape
* depth=Int - max depth of recursion while searching for links
* file=String - path to file where to store results of execution
* showCmd=Bool - show results in command line?
* sameDomain=Bool - include in searching only domain specified by the url flag