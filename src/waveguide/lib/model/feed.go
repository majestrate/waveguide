package model

import (
	"encoding/xml"
	"net/url"
	"time"
)

type Link struct {
	URL string `xml:"href,attr"`
}

func NewLink(domain, path string) Link {
	u := &url.URL{
		Scheme: "https",
		Host:   domain,
		Path:   path,
	}
	return Link{
		URL: u.String(),
	}
}

type AtomFeedEntry interface {
	xml.Marshaler
	CreatedAt() time.Time
}

type atomFeedImpl struct {
	Title    string          `xml:"title"`
	SubTitle string          `xml:"subtitle"`
	ID       string          `xml:"id"`
	Link     Link            `xml:"link"`
	Updated  time.Time       `xml:"updated"`
	Entries  []AtomFeedEntry `xml:"entry"`
}
