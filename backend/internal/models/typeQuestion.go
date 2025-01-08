package models

type TypeQuiz struct {
	Name        string   `json:"name"`
	QuestionsID []string `json:"questions_id"`
}

// Implement the Identifiable interface
func (u *TypeQuiz) GetID() string {
	return u.Name
}
