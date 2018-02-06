package api

type Resource struct {
	Name    string `json:"name"`
	URL     string `json:"url"`
	SHA1Hex string `json:"sha1"`
	Size    int64  `json:"size"`
}
