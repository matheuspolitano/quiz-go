package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/matheuspolitano/quiz-go/backend/internal/models"
	"github.com/matheuspolitano/quiz-go/backend/internal/token"
	"github.com/matheuspolitano/quiz-go/backend/internal/utils"
)

type userAnswerRequest struct {
	Answer string `json:"answer" binding:"required"`
}

type generalScore struct {
	UserQuiz             *models.QuestionFlow `json:"user_quiz"`
	GeneralAccuracyRates float32              `json:"general_accuracy_rates"`
}

func (svc *Server) getQuestion(ctx *gin.Context) {
	id := ctx.Param("questionID")
	question, err := svc.store.GetQuestion(id)
	if err != nil {
		SendError(ctx, "", err.Error(), http.StatusNotFound)
	}
	ctx.JSON(http.StatusAccepted, question)
}

func (svc *Server) joinQuiz(ctx *gin.Context) {
	typeQuiz := ctx.Param("typeQuiz")
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	questionFlow, err := svc.store.AddQuestionFlow(authPayload.Username, typeQuiz)
	if err != nil {
		SendError(ctx, "", err.Error(), http.StatusBadRequest)
	}
	ctx.JSON(http.StatusAccepted, questionFlow)
}
func (svc *Server) nextQuestion(ctx *gin.Context) {
	typeQuiz := ctx.Param("typeQuiz")
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	question, err := svc.store.NextQuestion(utils.CombineIDs(authPayload.Username, typeQuiz))
	if err != nil {
		SendError(ctx, "", err.Error(), http.StatusNotFound)
	}
	ctx.JSON(http.StatusAccepted, question)
}

func (svc *Server) generalScore(ctx *gin.Context) {
	typeQuiz := ctx.Param("typeQuiz")
	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
	question, generalAccuracyRates, err := svc.store.GetScoreUser(authPayload.Username, typeQuiz)
	if err != nil {
		SendError(ctx, "", err.Error(), http.StatusNotFound)
	}
	ctx.JSON(http.StatusAccepted, &generalScore{
		UserQuiz:             question,
		GeneralAccuracyRates: generalAccuracyRates,
	})
}

func (svc *Server) answerQuestion(ctx *gin.Context) {
	id := ctx.Param("typeQuiz")
	questionID := ctx.Param("questionID")

	var req userAnswerRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendError(ctx, "error in bind body", err.Error(), http.StatusBadRequest)
		return
	}

	authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)

	history, err := svc.store.AddAnswer(utils.CombineIDs(authPayload.Username, id), questionID, req.Answer)
	if err != nil {
		SendError(ctx, "", err.Error(), http.StatusNotFound)
	}
	ctx.JSON(http.StatusAccepted, history)
}
