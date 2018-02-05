package frontend

const DefaultChatID = "5"

/** get chat id for livestream */
func (r *Routes) ChatIDForStream(id string) (chatID string) {
	// TODO: fetch
	chatID = DefaultChatID
	return
}

/** get chat id for canned content */
func (r *Routes) ChatIDForVideo(id string) (chatID string) {
	// TODO: fetch
	chatID = DefaultChatID
	return
}
