package adn

import (
	"fmt"
)

type Post struct {
	Text        string       `json:"text"`
	Annotations []Annotation `json:"annotations,omitempty"`
}

func (c *Client) SubmitPost(token string, channel ChanID, post Post) (err error) {
	resp, e := c.doHTTP("POST", fmt.Sprintf("/channels/%s/messages?include_annotations=1", channel.String()), token, post)
	err = e
	if err == nil {
		if resp.StatusCode != 200 {
			err = ErrInvalidBackendResponse
		}
	}
	return
}
