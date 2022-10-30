# Sitemap Builder
Command line tool for recursive scraping links of a website source
## How to use
Git clone this repository to your local machine, then use 
```sh
go build ./cmd/api
```
and execute the built source code
or just
```sh
go run ./cmd/api
```
## Options
The app doesn't require any flags, but there's couple of optional ones
* url=String - url of the website to scrape, default: https://atos.net
* depth=Int - max depth of recursion while searching for links, default 3
* file=String - path to file where to store results of execution, default ./result.txt
* showCmd=Bool - show results in command line?, default false