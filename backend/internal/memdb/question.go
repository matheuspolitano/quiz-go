package memdb

type Question struct {
	ID      string   `json:"id"` // new field for unique ID
	Prompt  string   `json:"prompt"`
	Options []string `json:"options"`
	Answer  string   `json:"answer"`
}

// Implement the Identifiable interface
func (q *Question) GetID() string {
	return q.ID
}
