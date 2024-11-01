package main

import (
	"encoding/xml"
	"io"
)

const XML_NS = "http://www.sitemaps.org/schemas/sitemap/0.9"

type loc struct {
	Value string `xml:"loc"`
}

type urlSet struct {
	Urls  []loc  `xml:"url"`
	XMLns string `xml:"xmlns,attr"`
}

func GenerateSitemap(pages []string, w io.Writer) {
	siteMap := urlSet{
		XMLns: XML_NS,
	}
	for _, page := range pages {
		siteMap.Urls = append(siteMap.Urls, loc{page})
	}

	encoder := xml.NewEncoder(w)
	encoder.Indent("", "  ")

	w.Write([]byte(xml.Header))
	if err := encoder.Encode(siteMap); err != nil {
		panic(err)
	}
}
