package models

import (
	"time"

	"github.com/matheuspolitano/quiz-go/backend/internal/utils"
)

type QuestionFlow struct {
	UserID       string    `json:"user_id"`
	TypeQuizName string    `json:"type_quiz"`
	History      []string  `json:"history"`
	CreatedAt    time.Time `json:"created_at"`
	ClosedAt     time.Time `json:"closed_at,omitempty"`
	AccuracyRate float32   `json:"accuracy_rate"`
}

// the Identifiable interface
func (q *QuestionFlow) GetID() string {
	return utils.CombineIDs(q.UserID, q.TypeQuizName)
}
