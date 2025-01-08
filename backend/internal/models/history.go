package models

import "time"

type History struct {
	ID         string    `json:"id"` // new field for unique ID
	UserID     string    `json:"used_id"`
	QuestionID string    `json:"question_id"`
	Answer     string    `json:"answer"`
	IsRight    bool      `json:"is_right"`
	CreatedAt  time.Time `json:"created_at"`
}

// Implement the Identifiable interface
func (h *History) GetID() string {
	return h.ID
}
