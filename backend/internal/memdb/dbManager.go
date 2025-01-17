package memdb

import (
	"fmt"
	"sync"

	"github.com/matheuspolitano/quiz-go/backend/internal/models"
)

type DBManager struct {
	userProgressRepo  *Repository[*models.User]
	historyRepo       *Repository[*models.History]
	questionRepo      *Repository[*models.Question]
	TypeQuizRepo      *Repository[*models.TypeQuiz]
	questionsFlowRepo *Repository[*models.QuestionFlow]

	globalMu sync.Mutex
}

func NewDBManager() (*DBManager, error) {
	userRepo, err := NewRepositoryDefault[*models.User]("users")
	if err != nil {
		return nil, fmt.Errorf("failed to create user progress repo: %v", err)
	}

	historyRepo, err := NewRepositoryDefault[*models.History]("history")
	if err != nil {
		return nil, fmt.Errorf("failed to create history repo: %v", err)
	}

	questionRepo, err := NewRepositoryDefault[*models.Question]("questions")
	if err != nil {
		return nil, fmt.Errorf("failed to create question repo: %v", err)
	}

	TypeQuizRepo, err := NewRepositoryDefault[*models.TypeQuiz]("typesQuiz")
	if err != nil {
		return nil, fmt.Errorf("failed to create type question repo: %v", err)
	}

	questionsFlowRepo, err := NewRepositoryDefault[*models.QuestionFlow]("questionsFlows")
	if err != nil {
		return nil, fmt.Errorf("failed to create question repo: %v", err)
	}

	return &DBManager{
		userProgressRepo:  userRepo,
		historyRepo:       historyRepo,
		questionRepo:      questionRepo,
		TypeQuizRepo:      TypeQuizRepo,
		questionsFlowRepo: questionsFlowRepo,
	}, nil
}
