package models

// AccessTokenResponse is the structure we expect when we login successfully.
type AccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}

// QuizType represents the structure of a quiz type from GET /quiz/types.
type QuizType struct {
	Name        string   `json:"name"`
	QuestionsID []string `json:"questions_id"`
}

// Question represents the structure of each question from the server.
type Question struct {
	ID      string   `json:"id"`
	Prompt  string   `json:"prompt"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"` // The correct answer (often not sent in real scenarios)
}

// ScoreResponse is the structure of the final score response from the server.
type ScoreResponse struct {
	UserQuiz struct {
		UserID       string   `json:"user_id"`
		TypeQuiz     string   `json:"type_quiz"`
		History      []string `json:"history"`
		CreatedAt    string   `json:"created_at"`
		ClosedAt     string   `json:"closed_at"`
		AccuracyRate float64  `json:"accuracy_rate"`
	} `json:"user_quiz"`
	GeneralAccuracyRates float64 `json:"general_accuracy_rates"`
}
