package memdb

import (
	"time"
)

type QuestionFlow struct {
	ID             string    `json:"id"`
	UserID         string    `json:"used_id"`
	TypeQuestionID string    `json:"question_id"`
	History        []string  `json:"history"`
	CreatedAt      time.Time `json:"created_at"`
	ClosedAt       time.Time `json:"closed_at"`
	AccuracyRate   float32   `json:"accuracy_rate"`
}

// Implement the Identifiable interface
func (q *QuestionFlow) GetID() string {
	return combineIDs(q.UserID, q.TypeQuestionID)
}
