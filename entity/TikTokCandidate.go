package entity

import "time"

type TikTokCandidate struct {
	ID          uint64    `json:"-"`
	TikTokID    int       `json:"id"`
	Description string    `json:"desc"`
	Username    string    `json:"username"`
	ShareLink   string    `json:"view_link"`
	Thumbnail   string    `json:"thumbnail"`
	Duration    int       `json:"duration"`
	DatePosted  time.Time `json:"date_posted"`
}

var Candidates = []TikTokCandidate{
	{ID: 1, TikTokID: 1234567, Description: "Cool video", Username: "dolbaeb"},
	{ID: 2, TikTokID: 7654321, Description: "Not a Cool video", Username: "dolbaeb"},
	{ID: 3, TikTokID: 1111111, Description: "Tool tip missing", Username: "konch"},
}

func NewTikTokCandidate() *TikTokCandidate {
	return &TikTokCandidate{}
}
