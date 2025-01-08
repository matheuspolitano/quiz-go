package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/spf13/cobra"
)

// Question represents a single quiz question.
type Question struct {
	Prompt  string
	Options []string // e.g. []string{"A: Madrid", "B: Paris", "C: Berlin", "D: Rome"}
	Answer  string   // e.g. "B"
}

// Quiz holds the list of questions and the current index.
type Quiz struct {
	questions []Question
	current   int
	score     int
}

// NewQuiz creates a new Quiz instance with the provided questions.
func NewQuiz(qs []Question) *Quiz {
	return &Quiz{
		questions: qs,
		current:   0,
		score:     0,
	}
}

// GetNextQuestion returns the next question if it exists; otherwise, returns nil.
func (q *Quiz) GetNextQuestion() *Question {
	if q.current < len(q.questions) {
		return &q.questions[q.current]
	}
	return nil
}

// CheckAnswer compares the user’s answer to the correct answer and updates the score.
func (q *Quiz) CheckAnswer(answer string) bool {
	question := q.questions[q.current]
	if strings.EqualFold(answer, question.Answer) {
		q.score++
		return true
	}
	return false
}

// MoveToNext moves the quiz to the next question.
func (q *Quiz) MoveToNext() {
	q.current++
}

// IsFinished checks whether the quiz has no more questions.
func (q *Quiz) IsFinished() bool {
	return q.current >= len(q.questions)
}

// ScoreReport returns a formatted string with the final score.
func (q *Quiz) ScoreReport() string {
	return fmt.Sprintf("Quiz finished! Your score is %d out of %d.", q.score, len(q.questions))
}

// readAnswer repeatedly prompts the user for an answer until a valid one is provided.
func readAnswer(reader *bufio.Reader, validOptions []string) (string, error) {
	for {
		fmt.Print("Your answer (A, B, C, D): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(strings.ToUpper(input))

		// Check if the input is one of the valid options: "A", "B", "C", or "D".
		for _, opt := range validOptions {
			if input == opt {
				return input, nil
			}
		}

		fmt.Printf("Invalid answer: %s. Please enter one of %v.\n", input, validOptions)
	}
}

// runQuiz runs the quiz logic in the terminal.
func runQuiz() {
	// Sample questions
	questions := []Question{
		{
			Prompt:  "What is the capital of France?",
			Options: []string{"A: Madrid", "B: Paris", "C: Berlin", "D: Rome"},
			Answer:  "B",
		},
		{
			Prompt:  "Which country is known as the Land of the Rising Sun?",
			Options: []string{"A: China", "B: South Korea", "C: Japan", "D: Thailand"},
			Answer:  "C",
		},
		{
			Prompt:  "What is the largest country by area?",
			Options: []string{"A: Russia", "B: Canada", "C: China", "D: United States"},
			Answer:  "A",
		},
		{
			Prompt:  "Which country has the highest population?",
			Options: []string{"A: India", "B: China", "C: United States", "D: Indonesia"},
			Answer:  "B",
		},
	}

	quiz := NewQuiz(questions)
	reader := bufio.NewReader(os.Stdin)
	validOptions := []string{"A", "B", "C", "D"}

	fmt.Println("Welcome to the Countries Quiz!")
	fmt.Println("=============================")

	for {
		// Retrieve next question
		q := quiz.GetNextQuestion()
		if q == nil {
			break // No more questions
		}

		// Print question
		fmt.Printf("\nQuestion %d: %s\n", quiz.current+1, q.Prompt)
		for _, opt := range q.Options {
			fmt.Println(opt)
		}

		// Read user answer (with validation)
		answer, err := readAnswer(reader, validOptions)
		if err != nil {
			fmt.Printf("Error reading answer: %v\n", err)
			return
		}

		// Check answer correctness
		if quiz.CheckAnswer(answer) {
			fmt.Println("Correct!")
		} else {
			fmt.Printf("Incorrect! The correct answer is %s.\n", q.Answer)
		}

		// Move to next question
		quiz.MoveToNext()
	}

	// Display final score
	fmt.Printf("\n%s\n", quiz.ScoreReport())
}

// --------------------------------------------------------
// Cobra CLI setup below
// --------------------------------------------------------

var rootCmd = &cobra.Command{
	Use:   "quiz",
	Short: "A CLI application for country quizzes",
	Long:  "quiz is a terminal-based application that tests your knowledge about world countries.",
}

var startCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the quiz",
	Run: func(cmd *cobra.Command, args []string) {
		runQuiz()
	},
}

func init() {
	rootCmd.AddCommand(startCmd)
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
