package entity

import "time"

type GetActiveChatsResponse []Chat
type Chat struct {
	ID           int
	Title        string
	Author       string
	TimeToLive   int
	Tags         []ChatTag
	LastMessages []Message
}
type ChatTag struct {
	Title string
}
type Message struct {
	Author string
	Text   string
	Time   time.Time
}
