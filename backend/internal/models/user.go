package models

import (
	"time"
)

type User struct {
	Username         string    `json:"username"`
	CreatedAt        time.Time `json:"created_at"`
	QuestionsFlowsID []string  `json:"questions_flows_id"`
}

// Implement the Identifiable interface
func (u *User) GetID() string {
	return u.Username
}
