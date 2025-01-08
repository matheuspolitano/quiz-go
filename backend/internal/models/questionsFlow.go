package models

import (
	"time"

	"github.com/matheuspolitano/quiz-go/backend/internal/utils"
)

type QuestionFlow struct {
	ID             string    `json:"id"`
	UserID         string    `json:"user_id"`
	TypeQuestionID string    `json:"question_id"`
	History        []string  `json:"history"`
	CreatedAt      time.Time `json:"created_at"`
	ClosedAt       time.Time `json:"closed_at,omitempty"`
	AccuracyRate   float32   `json:"accuracy_rate"`
}

// the Identifiable interface
func (q *QuestionFlow) GetID() string {
	return utils.CombineIDs(q.UserID, q.TypeQuestionID)
}
