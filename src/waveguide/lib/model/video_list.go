package model

import (
	"encoding/xml"
	"net/url"
	"time"
	"waveguide/lib/util"
)

type VideoList struct {
	LastUpdated time.Time
	Videos      []VideoInfo
}

type VideoFeed struct {
	List  VideoList
	Owner *UserInfo
	URL   *url.URL
}

func (feed *VideoFeed) Title() string {
	return "Video Feed for " + feed.Owner.Name
}

func (feed *VideoFeed) Entries() (entries []AtomFeedEntry) {
	for _, v := range feed.List.Videos {
		v.BaseURL = feed.URL
		entries = append(entries, v)
	}
	return
}

func (feed *VideoFeed) ID() string {
	return util.SHA256(feed.URL.String())
}

func (feed *VideoFeed) MarshalXML(e *xml.Encoder, start xml.StartElement) (err error) {

	latest := time.Unix(0, 0)
	for _, v := range feed.List.Videos {
		u := v.CreatedAt()
		if u.After(latest) {
			latest = u
		}
	}

	u := feed.URL
	start.Name.Local = "feed"
	start.Name.Space = "http://www.w3.org/2005/Atom"
	err = e.EncodeElement(&atomFeedImpl{
		Title:    feed.Title(),
		SubTitle: feed.Title(),
		Link:     NewLink(u.Host, u.RequestURI()),
		ID:       feed.ID(),
		Entries:  feed.Entries(),
		Updated:  feed.List.LastUpdated,
	}, start)
	return
}
