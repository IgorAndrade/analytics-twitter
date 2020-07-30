package model

import (
	"encoding/json"
	"time"
)

type Post struct {
	CreatedAt time.Time `json:"created_at,omitempty"`
	Author    string    `json:"author,omitempty"`
	Hastags   []string  `json:"hastags,omitempty"`
	Text      string    `json:"text,omitempty"`
	Location  string    `json:"location,omitempty"`
	Lang      string    `json:"lang,omitempty"`
}

func (p Post) String() string {
	b, _ := json.Marshal(p)
	return string(b)
}
