package oauth

type TokenInfo struct {
	Scopes []string `json:"scopes"`
	User   User     `json:"user"`
}

type TokenInfoRequest struct {
	Data TokenInfo    `json:"data"`
	Meta MetaResponse `json:"meta"`
}
