package oauth

type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Token    string `json:"-"`
}

type Token struct {
	User User `json:"user"`
}

type TokenRequest struct {
	AccessToken string `json:"access_token"`
	Token       Token  `json:"token"`
}
