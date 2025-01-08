package quiz

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"

	"github.com/matheuspolitano/quiz-go/client/internal/api"
	"github.com/matheuspolitano/quiz-go/client/internal/models"
)

// RunQuizFlow orchestrates the entire quiz process:
// 1. Ask user for username and login
// 2. Retrieve quiz types
// 3. User selects a quiz type
// 4. Join the quiz
// 5. Fetch next question, answer, repeat
// 6. Retrieve final score
func RunQuizFlow(baseURL string) {
	// Create a new client for the quiz API.
	client := api.NewClient(baseURL)
	reader := bufio.NewReader(os.Stdin)

	color.Cyan("Welcome to the Quiz CLI!")
	fmt.Println(strings.Repeat("=", 40))

	// 1. Prompt for username and login
	username, err := promptForUsername(reader)
	if err != nil {
		color.Red("Error reading username: %v", err)
		return
	}

	// Use the client to log in
	err = client.Login(username)
	if err != nil {
		color.Red("Login failed: %v", err)
		return
	}

	color.Green("Successfully logged in! Your access token is saved.\n")

	// Outer loop to allow multiple quiz attempts
	for {
		// 2. Retrieve quiz types
		quizTypes, err := client.GetQuizTypes()
		if err != nil {
			color.Red("Unable to fetch quiz types: %v", err)
			return
		}
		if len(quizTypes) == 0 {
			color.Yellow("No quiz types available.")
			return
		}

		// 3. Prompt user to select quiz type
		selectedQuizType, err := promptForQuizType(reader, quizTypes)
		if err != nil {
			color.Red("Invalid quiz type selection: %v", err)
			continue
		}

		// 4. Join the quiz
		if err := client.JoinQuiz(selectedQuizType); err != nil {
			color.Red("Cannot join quiz: %v", err)
			continue
		}
		color.Green("Joined quiz: %s", selectedQuizType)

		// 5. Question loop
		err = questionLoop(reader, client, selectedQuizType)
		if err != nil {
			color.Red("Error during question flow: %v", err)
		}

		// 6. Fetch final score
		err = fetchAndDisplayScore(client, selectedQuizType)
		if err != nil {
			color.Red("Error fetching final score: %v", err)
			return
		}

		// 7. Prompt to try another quiz type
		tryAnother, err := promptForAnotherQuiz(reader)
		if err != nil {
			color.Red("Error reading input: %v", err)
			return
		}
		if !tryAnother {
			color.Cyan("Thank you for playing! Goodbye.")
			break
		}
	}
}

// promptForAnotherQuiz asks the user if they want to try another quiz type.
func promptForAnotherQuiz(reader *bufio.Reader) (bool, error) {
	for {
		fmt.Print("Do you want to try another quiz type? (Y/N): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return false, err
		}
		input = strings.TrimSpace(strings.ToUpper(input))

		if input == "Y" || input == "YES" {
			return true, nil
		} else if input == "N" || input == "NO" {
			return false, nil
		} else {
			fmt.Println("Please answer Y or N.")
		}
	}
}

// questionLoop repeatedly fetches the next question, prompts for an answer,
// and sends it to the server until no more questions remain (or an error occurs).
func questionLoop(reader *bufio.Reader, client *api.Client, quizType string) error {
	for {
		// Attempt to fetch the next question
		question, err := client.GetNextQuestion(quizType)
		if err != nil {
			// If the server says “question flow is already closed”, we’re done.
			if strings.Contains(strings.ToLower(err.Error()), "already closed") {
				color.Yellow("\nNo more questions. The quiz is finished.\n")
				return nil
			}
			return err
		}

		// Display the question
		color.Cyan("\nQuestion (%s)\n", question.ID)
		fmt.Println(question.Prompt)
		for _, opt := range question.Options {
			fmt.Println(opt)
		}

		// Prompt for answer
		answer, err := promptForAnswer(reader)
		if err != nil {
			return err
		}

		// Submit the answer
		resultAnswer, err := client.SubmitAnswer(quizType, question.ID, answer)
		if err != nil {
			// If "AddAnswer: question already answer", we simply move to the next
			if strings.Contains(strings.ToLower(err.Error()), "already answer") {
				color.Yellow("It seems you already answered this question. Moving on...")
				continue
			}
			return err
		}
		color.Green("Your answer (%s) was submitted.\n", answer)
		if resultAnswer.Answer == resultAnswer.ExpectedAnswer {
			color.Green("Your answer is right :) \n")
		} else {
			color.Red("Your answer is wrong :(. The right is (%s) \n", resultAnswer.ExpectedAnswer)
		}

	}
}

// fetchAndDisplayScore calls the final score endpoint and displays the result.
func fetchAndDisplayScore(client *api.Client, quizType string) error {
	scoreResp, err := client.GetScore(quizType)
	if err != nil {
		return err
	}

	color.Magenta("========================================")
	color.Magenta("              FINAL SCORE              ")
	color.Magenta("========================================")

	fmt.Printf("Quiz Type        : %s\n", scoreResp.UserQuiz.TypeQuiz)
	fmt.Printf("Answered         : %d questions\n", len(scoreResp.UserQuiz.History))
	fmt.Printf("Accuracy Rate    : %.2f%%\n", scoreResp.UserQuiz.AccuracyRate*100)
	fmt.Printf("General Avg Rate : %.2f%%\n", scoreResp.GeneralAccuracyRates*100)
	fmt.Printf("Quiz closed at   : %s\n", scoreResp.UserQuiz.ClosedAt)

	color.Magenta("========================================")

	return nil
}

// promptForUsername reads the username from stdin.
func promptForUsername(reader *bufio.Reader) (string, error) {
	fmt.Print("Enter your username: ")
	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(input), nil
}

// promptForQuizType displays the available quiz types and lets the user select one.
func promptForQuizType(reader *bufio.Reader, quizTypes []models.QuizType) (string, error) {
	fmt.Println("\nAvailable quiz types:")
	for i, qt := range quizTypes {
		fmt.Printf("  %d) %s\n", i+1, qt.Name)
	}
	fmt.Print("Select a quiz type by number: ")

	input, err := reader.ReadString('\n')
	if err != nil {
		return "", err
	}
	input = strings.TrimSpace(input)

	var index int
	_, err = fmt.Sscanf(input, "%d", &index)
	if err != nil {
		return "", fmt.Errorf("please enter a number")
	}

	index -= 1 // Convert to zero-based index
	if index < 0 || index >= len(quizTypes) {
		return "", fmt.Errorf("invalid quiz type selection")
	}
	return quizTypes[index].Name, nil
}

// promptForAnswer repeatedly prompts the user for an answer until a valid one is provided.
func promptForAnswer(reader *bufio.Reader) (string, error) {
	validOptions := []string{"A", "B", "C", "D"}
	for {
		fmt.Print("Your answer (A, B, C, D): ")
		input, err := reader.ReadString('\n')
		if err != nil {
			return "", err
		}
		input = strings.TrimSpace(strings.ToUpper(input))

		for _, opt := range validOptions {
			if input == opt {
				return input, nil
			}
		}
		fmt.Printf("Invalid answer: %s. Please enter one of %v.\n", input, validOptions)
	}
}
