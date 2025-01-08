package memdb

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/matheuspolitano/quiz-go/backend/internal/models"
	"github.com/matheuspolitano/quiz-go/backend/internal/utils"
)

var (
	ErrFlowClosed           = errors.New("question flow is already closed")
	ErrUsernameAlreadyExist = errors.New("username already exist")
	ErrNoQuestions          = errors.New("no questions available for this question type")
	ErrAllQuestionsAnswered = errors.New("all questions have been answered in this flow")
)

func (db *DBManager) GetScoreUser(userID, quizType string) (*models.QuestionFlow, float32, error) {
	questionFlow, err := db.questionsFlowRepo.FindByID(utils.CombineIDs(userID, quizType))
	if err != nil {
		return nil, 0, err
	}

	allQuestionsFlow, err := db.questionsFlowRepo.ListAll()
	if err != nil {
		return nil, 0, err
	}
	var total float32
	var length float32
	for _, item := range allQuestionsFlow {
		if !item.ClosedAt.IsZero() && item.TypeQuizName == questionFlow.TypeQuizName {
			total += float32(item.AccuracyRate)
			length++
		}

	}

	return questionFlow, total / length, nil
}

func (db *DBManager) GetQuestion(id string) (*models.Question, error) {
	question, err := db.questionRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return question, nil
}

func (db *DBManager) ListAllTypes() ([]*models.TypeQuiz, error) {
	TypeQuizs, err := db.TypeQuizRepo.ListAll()
	if err != nil {
		return nil, err
	}
	return TypeQuizs, nil
}

func (db *DBManager) CreateUser(username string) (*models.User, error) {
	_, err := db.userProgressRepo.FindByID(username)
	if err == nil {
		return nil, ErrUsernameAlreadyExist
	}
	user := &models.User{
		Username:         username,
		CreatedAt:        time.Now(),
		QuestionsFlowsID: make([]string, 0),
	}
	err = db.userProgressRepo.Save(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func (db *DBManager) AddQuestionFlow(userID, TypeQuizName string) (*models.QuestionFlow, error) {
	db.globalMu.Lock()
	defer db.globalMu.Unlock()

	_, err := db.TypeQuizRepo.FindByID(TypeQuizName)
	if err != nil {
		return nil, fmt.Errorf("AddQuestionFlow: TypeQuiz does not exist: %s", err.Error())
	}

	flowID := utils.CombineIDs(userID, TypeQuizName)
	_, err = db.questionsFlowRepo.FindByID(flowID)
	if err == nil {
		return nil, fmt.Errorf("AddQuestionFlow: question flow already exists for user %s and type %s", userID, TypeQuizName)
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	newFlow := &models.QuestionFlow{
		UserID:       userID,
		TypeQuizName: TypeQuizName,
		CreatedAt:    time.Now(),
		AccuracyRate: 1.0,
		ClosedAt:     time.Time{},
		History:      make([]string, 0),
	}
	user, uErr := db.userProgressRepo.FindByID(userID)
	if uErr != nil {
		return nil, fmt.Errorf("AddQuestionFlow: failed to find user: %s", uErr.Error())
	}
	if saveErr := db.questionsFlowRepo.Save(newFlow); saveErr != nil {
		return nil, fmt.Errorf("AddQuestionFlow: failed to save new flow: %s", saveErr.Error())
	}

	user.QuestionsFlowsID = append(user.QuestionsFlowsID, newFlow.GetID())
	if saveUserErr := db.userProgressRepo.Save(user); saveUserErr != nil {
		return nil, fmt.Errorf("AddQuestionFlow: failed to update user flows: %s", saveUserErr.Error())
	}

	return newFlow, nil
}

// NextQuestion retrieves the next question for an existing QuestionFlow.
// Returns (Question, nil) when a question is found,
// returns an error (ErrFlowClosed, ErrNoQuestions, ErrAllQuestionsAnswered, etc.) otherwise.
func (db *DBManager) NextQuestion(questionFlowID string) (*models.Question, error) {
	db.globalMu.Lock()
	defer db.globalMu.Unlock()

	qFlow, err := db.questionsFlowRepo.FindByID(questionFlowID)
	if err != nil {
		return nil, fmt.Errorf("NextQuestion: question flow not found: %s", err.Error())
	}
	if !qFlow.ClosedAt.IsZero() {
		return nil, ErrFlowClosed
	}
	tQuestion, err := db.TypeQuizRepo.FindByID(qFlow.TypeQuizName)
	if err != nil {
		return nil, fmt.Errorf("NextQuestion: TypeQuiz not found: %s", err.Error())
	}
	if len(tQuestion.QuestionsID) == 0 {
		return nil, ErrNoQuestions
	}

	answeredQuestionIDs := make(map[string]bool)
	for _, histID := range qFlow.History {
		histEntry, hErr := db.historyRepo.FindByID(histID)
		if hErr != nil {
			continue
		}
		answeredQuestionIDs[histEntry.QuestionID] = true
	}

	for _, qID := range tQuestion.QuestionsID {
		if !answeredQuestionIDs[qID] {
			nextQ, qErr := db.questionRepo.FindByID(qID)
			if qErr != nil {
				continue
			}
			return nextQ, nil
		}
	}

	var lastAnswerTime time.Time
	for _, histID := range qFlow.History {
		histEntry, _ := db.historyRepo.FindByID(histID)
		if histEntry.CreatedAt.After(lastAnswerTime) {
			lastAnswerTime = histEntry.CreatedAt
		}
	}
	if lastAnswerTime.IsZero() {
		lastAnswerTime = time.Now()
	}

	qFlow.ClosedAt = lastAnswerTime
	_ = db.questionsFlowRepo.Save(qFlow)

	return nil, ErrAllQuestionsAnswered
}

// AddAnswer stores an answer for a specific question in the flow.
func (db *DBManager) AddAnswer(questionFlowID, questionID, userAnswer string) (*models.History, error) {
	db.globalMu.Lock()
	defer db.globalMu.Unlock()

	qFlow, err := db.questionsFlowRepo.FindByID(questionFlowID)
	if err != nil {
		return nil, fmt.Errorf("AddAnswer: question flow not found: %s", err.Error())
	}
	if !qFlow.ClosedAt.IsZero() {
		return nil, ErrFlowClosed
	}

	typeQ, err := db.TypeQuizRepo.FindByID(qFlow.TypeQuizName)
	if err != nil {
		return nil, fmt.Errorf("AddAnswer: invalid TypeQuiz: %s", err.Error())
	}
	found := false
	for _, qID := range typeQ.QuestionsID {
		if qID == questionID {
			found = true
			break
		}
	}
	if !found {
		return nil, fmt.Errorf("AddAnswer: question %s not part of TypeQuiz %s", questionID, typeQ.Name)
	}

	answeredQuestionIDs := make(map[string]bool)
	for _, histID := range qFlow.History {
		histEntry, hErr := db.historyRepo.FindByID(histID)
		if hErr != nil {
			continue
		}
		answeredQuestionIDs[histEntry.QuestionID] = true
	}

	if answeredQuestionIDs[questionID] {
		return nil, fmt.Errorf("AddAnswer: question already answer. Use the next to get the question without answer")
	}

	questionObj, err := db.questionRepo.FindByID(questionID)
	if err != nil {
		return nil, fmt.Errorf("AddAnswer: cannot find question %s: %w", questionID, err)
	}

	newHist := &models.History{
		ID:             uuid.NewString(),
		UserID:         qFlow.UserID,
		QuestionID:     questionID,
		Answer:         userAnswer,
		ExpectedAnswer: questionObj.Answer,
		CreatedAt:      time.Now(),
	}
	if err := db.historyRepo.Save(newHist); err != nil {
		return nil, fmt.Errorf("AddAnswer: failed to save new History: %s", err.Error())
	}

	qFlow.History = append(qFlow.History, newHist.ID)

	correctCount := 0
	for _, histID := range qFlow.History {
		h, herr := db.historyRepo.FindByID(histID)
		if herr == nil && h.ExpectedAnswer == h.Answer {
			correctCount++
		}
	}
	totalAnswers := len(qFlow.History)
	if totalAnswers > 0 {
		qFlow.AccuracyRate = float32(correctCount) / float32(totalAnswers)
	} else {
		qFlow.AccuracyRate = 1.0
	}

	if err := db.questionsFlowRepo.Save(qFlow); err != nil {
		return nil, fmt.Errorf("AddAnswer: failed to update question flow with new history: %s", err.Error())
	}

	return newHist, nil
}
