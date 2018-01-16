package oauth

type User struct {
	Token    string `json:"access_token"`
	Username string `json:"username"`
	ID       string `json:"user_id"`
}
