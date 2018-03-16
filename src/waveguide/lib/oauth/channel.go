package oauth

type ChannelMembers struct {
	Users []string `json:"user_ids"`
}

type ChannelStats struct {
	Messages    int64 `json:"messages"`
	Subscribers int64 `json:"subscribers"`
}

type Channel struct {
	Counts  ChannelStats   `json:"counts"`
	Type    string         `json:"type"`
	Readers ChannelMembers `json:"readers"`
	Writers ChannelMembers `json:"writers"`
	ID      string         `json:"id"`
	Owner   User           `json:"owner"`
}

type ChannelResponse struct {
	Data Channel      `json:"data"`
	Meta MetaResponse `json:"meta"`
}
