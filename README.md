# Sitemap Builder

A command-line tool for generating XML sitemaps for websites. The sitemap builder crawls the specified URL up to a given depth and produces an XML sitemap that can be used by search engines for indexing.

## Project Structure

```shell
.
├── bin
    │
    └── sitemap
├── cmd
│   ├── builder.go # Contains the GenerateSitemap function to generate the XML sitemap
│   └── main.go
├── Dockerfile
├── go.mod
├── go.sum
├── Makefile
├── pkg
│   ├── collections
│   │   ├── stack.go
│   │   └── thread-safe
│   │       └── stack.go
│   ├── crawler
│   │   ├── crawler.go # Crawler logic for traversing links on a website
│   │   └── crawler_test.go
│   └── link
│       ├── parser.go # Parser for extracting links from HTML
│       └── parser_test.go
└── README.md
```


## Installation

1. Clone the repository:
   ```shell
   git clone https://github.com/yourusername/sitemap-builder.git
   cd sitemap-builder

2. Build the project using the `Makefile`:
   ```shell
   make build
   ```
This will create a `sitemap` executable in the `bin/` directory.

## Usage

Run the application with the `-url` flag to specify the URL and the `-depth` flag to set the crawling depth:
```shell
./bin/sitemap -url=https://example.com -depth=3
```

### Optional Flags
* `-url`: URL to crawl (required)
* `-depth`: Maximum depth to crawl (default: 3)
* `-out`: Output file for the sitemap XML (default: `./map.xml`)

## Running with Docker
You can also run the sitemap builder in a Docker container.

### Building the Docker Image
```shell
make docker-build
```
### Running the Container and Saving Output on the Host
To run the container and save the output file to a specified directory on your host system, use a bind mount:

1. Create an output directory on your host (if it doesn’t already exist):
```bash
mkdir -p output
```
2. Run the container, mounting the output directory and specifying the output path:
```shell
docker run --rm -v $(pwd)/output:/app/output sitemap-builder -url=https://example.com -depth=3 -out=/app/output/map.xml
```