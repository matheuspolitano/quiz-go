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
	ErrNoQuestions          = errors.New("no questions available for this question type")
	ErrAllQuestionsAnswered = errors.New("all questions have been answered in this flow")
)

func (db *DBManager) AddQuestionFlow(userID, typeQuestionID string) (*models.QuestionFlow, error) {
	db.globalMu.Lock()
	defer db.globalMu.Unlock()

	_, err := db.typeQuestionRepo.FindByID(typeQuestionID)
	if err != nil {
		return nil, fmt.Errorf("AddQuestionFlow: typeQuestion does not exist: %s", err.Error())
	}

	flowID := utils.CombineIDs(userID, typeQuestionID)
	_, err = db.questionsFlowRepo.FindByID(flowID)
	if err == nil {
		return nil, fmt.Errorf("AddQuestionFlow: question flow already exists for user %s and type %s", userID, typeQuestionID)
	}
	if !errors.Is(err, ErrNotFound) {
		return nil, err
	}

	newFlow := &models.QuestionFlow{
		UserID:         userID,
		TypeQuestionID: typeQuestionID,
		CreatedAt:      time.Now(),
		AccuracyRate:   1.0,
	}

	if saveErr := db.questionsFlowRepo.Save(&newFlow); saveErr != nil {
		return nil, fmt.Errorf("AddQuestionFlow: failed to save new flow: %s", saveErr.Error())
	}

	user, uErr := db.userProgressRepo.FindByID(userID)
	if uErr != nil {
		return nil, fmt.Errorf("AddQuestionFlow: failed to find user: %s", uErr.Error())
	}
	user.QuestionsFlowsID = append(user.QuestionsFlowsID, newFlow.GetID())
	if saveUserErr := db.userProgressRepo.Save(&user); saveUserErr != nil {
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
	tQuestion, err := db.typeQuestionRepo.FindByID(qFlow.TypeQuestionID)
	if err != nil {
		return nil, fmt.Errorf("NextQuestion: typeQuestion not found: %s", err.Error())
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
	_ = db.questionsFlowRepo.Save(&qFlow)

	return nil, ErrAllQuestionsAnswered
}

// AddAnswer stores an answer for a specific question in the flow.
func (db *DBManager) AddAnswer(questionFlowID, questionID, userAnswer string) error {
	db.globalMu.Lock()
	defer db.globalMu.Unlock()

	qFlow, err := db.questionsFlowRepo.FindByID(questionFlowID)
	if err != nil {
		return fmt.Errorf("AddAnswer: question flow not found: %s", err.Error())
	}
	if !qFlow.ClosedAt.IsZero() {
		return ErrFlowClosed
	}

	typeQ, err := db.typeQuestionRepo.FindByID(qFlow.TypeQuestionID)
	if err != nil {
		return fmt.Errorf("AddAnswer: invalid typeQuestion: %s", err.Error())
	}
	found := false
	for _, qID := range typeQ.QuestionsID {
		if qID == questionID {
			found = true
			break
		}
	}
	if !found {
		return fmt.Errorf("AddAnswer: question %s not part of typeQuestion %s", questionID, typeQ.ID)
	}

	questionObj, err := db.questionRepo.FindByID(questionID)
	if err != nil {
		return fmt.Errorf("AddAnswer: cannot find question %s: %w", questionID, err)
	}

	newHist := &models.History{
		ID:         uuid.NewString(),
		UserID:     qFlow.UserID,
		QuestionID: questionID,
		Answer:     userAnswer,
		IsRight:    questionObj.Answer == userAnswer,
		CreatedAt:  time.Now(),
	}
	if err := db.historyRepo.Save(&newHist); err != nil {
		return fmt.Errorf("AddAnswer: failed to save new History: %s", err.Error())
	}

	qFlow.History = append(qFlow.History, newHist.ID)

	correctCount := 0
	for _, histID := range qFlow.History {
		h, herr := db.historyRepo.FindByID(histID)
		if herr == nil && h.IsRight {
			correctCount++
		}
	}
	totalAnswers := len(qFlow.History)
	if totalAnswers > 0 {
		qFlow.AccuracyRate = float32(correctCount) / float32(totalAnswers)
	} else {
		qFlow.AccuracyRate = 1.0
	}

	if err := db.questionsFlowRepo.Save(&qFlow); err != nil {
		return fmt.Errorf("AddAnswer: failed to update question flow with new history: %s", err.Error())
	}

	return nil
}
