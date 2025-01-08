package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/matheuspolitano/quiz-go/client/internal/models"
)

// Client wraps the configuration needed to make API calls.
type Client struct {
	BaseURL    string
	httpClient *http.Client
	token      string
}

// NewClient creates a new Client instance. The http.Client can be customized for timeouts, etc.
func NewClient(baseURL string) *Client {
	return &Client{
		BaseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second, // example timeout
		},
	}
}

// SetToken allows manually setting the token, if needed.
func (c *Client) SetToken(token string) {
	c.token = token
}

// Login sends a POST to /api/login with the given username
// and stores the received token in the client's token field.
func (c *Client) Login(username string) error {
	payload := map[string]string{"username": username}
	data, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("failed to marshal login payload: %w", err)
	}

	url := c.BaseURL + "/api/login"
	resp, err := c.httpClient.Post(url, "application/json", bytes.NewBuffer(data))
	if err != nil {
		return fmt.Errorf("failed to make login request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var tokenResp models.AccessTokenResponse
	if err := json.NewDecoder(resp.Body).Decode(&tokenResp); err != nil {
		return fmt.Errorf("decoding token response: %w", err)
	}

	c.token = tokenResp.AccessToken
	return nil
}

// GetQuizTypes retrieves the available quiz types for the user.
func (c *Client) GetQuizTypes() ([]models.QuizType, error) {
	req, err := http.NewRequest(http.MethodGet, c.BaseURL+"/api/quiz/types", nil)
	if err != nil {
		return nil, fmt.Errorf("creating GetQuizTypes request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending GetQuizTypes request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var quizTypes []models.QuizType
	if err := json.NewDecoder(resp.Body).Decode(&quizTypes); err != nil {
		return nil, fmt.Errorf("decoding quiz types: %w", err)
	}
	return quizTypes, nil
}

// JoinQuiz joins a specific quiz type for the logged-in user.
func (c *Client) JoinQuiz(quizType string) error {
	url := fmt.Sprintf("%s/api/quiz/joinQuiz/%s", c.BaseURL, quizType)
	req, err := http.NewRequest(http.MethodPost, url, nil)
	if err != nil {
		return fmt.Errorf("creating JoinQuiz request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("sending JoinQuiz request: %w", err)
	}
	defer resp.Body.Close()

	// 400 => user may already have joined this quiz.
	if resp.StatusCode == http.StatusBadRequest {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("already joined? status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	if resp.StatusCode != http.StatusAccepted && resp.StatusCode != http.StatusCreated {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	return nil
}

// GetNextQuestion fetches the next question for the user.
func (c *Client) GetNextQuestion(quizType string) (*models.Question, error) {
	url := fmt.Sprintf("%s/api/quiz/answer/%s/next", c.BaseURL, quizType)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating GetNextQuestion request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending GetNextQuestion request: %w", err)
	}
	defer resp.Body.Close()

	// If 400 with "question flow is already closed", no more questions.
	if resp.StatusCode == http.StatusBadRequest {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(string(bodyBytes))
	}
	if resp.StatusCode == http.StatusNotFound {
		return nil, errors.New("quiz has finished")
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var question models.Question
	if err := json.NewDecoder(resp.Body).Decode(&question); err != nil {
		return nil, fmt.Errorf("decoding question: %w", err)
	}
	return &question, nil
}

// SubmitAnswer sends the user’s answer to the server for a specific question.
func (c *Client) SubmitAnswer(quizType, questionID, answer string) (*models.History, error) {
	payload := map[string]string{"answer": answer}
	data, err := json.Marshal(payload)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal answer payload: %w", err)
	}

	url := fmt.Sprintf("%s/api/quiz/answer/%s/%s", c.BaseURL, quizType, questionID)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(data))
	if err != nil {
		return nil, fmt.Errorf("creating SubmitAnswer request: %w", err)
	}
	c.addAuthHeader(req)
	req.Header.Set("Content-Type", "application/json")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending SubmitAnswer request: %w", err)
	}
	defer resp.Body.Close()

	// Typically expect 202 Accepted if the answer was processed.
	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}
	var history models.History
	if err := json.NewDecoder(resp.Body).Decode(&history); err != nil {
		return nil, fmt.Errorf("decoding history: %w", err)
	}
	return &history, nil
}

// GetScore retrieves the user’s score for a given quiz type.
func (c *Client) GetScore(quizType string) (*models.ScoreResponse, error) {
	url := fmt.Sprintf("%s/api/quiz/answer/%s/score", c.BaseURL, quizType)
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("creating GetScore request: %w", err)
	}
	c.addAuthHeader(req)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("sending GetScore request: %w", err)
	}
	defer resp.Body.Close()

	// Expecting 202 Accepted
	if resp.StatusCode != http.StatusAccepted {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status %d: %s", resp.StatusCode, string(bodyBytes))
	}

	var score models.ScoreResponse
	if err := json.NewDecoder(resp.Body).Decode(&score); err != nil {
		return nil, fmt.Errorf("decoding score: %w", err)
	}
	return &score, nil
}

// addAuthHeader sets the "Authorization: Bearer {token}" header, if a token is available.
func (c *Client) addAuthHeader(req *http.Request) {
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}
}
