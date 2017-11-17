package model

import (
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

func (feed *VideoFeed) ToAtom() *AtomFeed {

	latest := time.Unix(0, 0)
	for _, v := range feed.List.Videos {
		u := v.CreatedAt()
		if u.After(latest) {
			latest = u
		}
	}

	u := feed.URL

	return &AtomFeed{
		Title:    feed.Title(),
		SubTitle: feed.Title(),
		Link:     NewLink(u.Host, u.RequestURI()),
		ID:       feed.ID(),
		Entries:  feed.Entries(),
		Updated:  feed.List.LastUpdated,
	}
}
