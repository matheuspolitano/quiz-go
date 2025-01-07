package memdb

import "time"

type UserProgress struct {
	ID               string    `json:"id"` // new field for unique ID
	Username         string    `json:"username"`
	CreatedAt        time.Time `json:"created_at"`
	questionsFlowsID []string  `json:"questions_flows"`
}

// Implement the Identifiable interface
func (u *UserProgress) GetID() string {
	return u.ID
}
