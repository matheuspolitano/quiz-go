package models

type TypeQuestion struct {
	ID          string   `json:"id"` // new field for unique ID
	Name        string   `json:"name"`
	QuestionsID []string `json:"questions_id"`
}

// Implement the Identifiable interface
func (u *TypeQuestion) GetID() string {
	return u.ID
}
