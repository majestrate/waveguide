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

type AtomFeed struct {
	Title    string          `xml:"title"`
	SubTitle string          `xml:"subtitle"`
	ID       string          `xml:"id"`
	Link     Link            `xml:"link"`
	Updated  time.Time       `xml:"updated"`
	Entries  []AtomFeedEntry `xml:"entry"`
}

func (feed *AtomFeed) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {
	start.Name.Local = "feed"
	start.Name.Space = "http://www.w3.org/2005/Atom"
	err = e.EncodeElement(feed, start)
	return
}
