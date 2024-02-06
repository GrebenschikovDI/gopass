package records

import "time"

type Record struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Name      string    `json:"name,omitempty"`
	Site      string    `json:"site,omitempty"`
	Login     string    `json:"login,omitempty"`
	Password  string    `json:"password,omitempty"`
	Info      string    `json:"info,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}
