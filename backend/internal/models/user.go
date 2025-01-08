package models

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID               string    `json:"id"` // new field for unique ID
	Username         string    `json:"username"`
	CreatedAt        time.Time `json:"created_at"`
	QuestionsFlowsID []string  `json:"questions_flows"`
}

// Implement the Identifiable interface
func (u *User) GetID() string {
	return u.ID
}

// NewUser Its create a new user with username
func NewUser(username string) User {
	return User{
		ID:        uuid.NewString(),
		Username:  username,
		CreatedAt: time.Now(),
	}
}
