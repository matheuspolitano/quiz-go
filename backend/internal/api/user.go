package api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matheuspolitano/quiz-go/backend/internal/memdb"
)

type createUserRequest struct {
	Username string `json:"username" binding:"required"`
}

type loginUserResponse struct {
	AccessToken string `json:"access_token"`
}

func (server *Server) startUser(ctx *gin.Context) {
	var req createUserRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		SendError(ctx, "error in bind body", err.Error(), http.StatusBadRequest)
		return
	}

	_, err := server.store.CreateUser(req.Username)
	if err != nil && err != memdb.ErrUsernameAlreadyExist {
		SendError(ctx, "", err.Error(), http.StatusBadRequest)
		return
	}

	token, _, err := server.tokenMaker.CreateToken(req.Username, "regular", time.Minute*60*24*365)
	if err != nil {
		SendError(ctx, "", err.Error(), http.StatusBadRequest)
		return
	}
	userResponse := &loginUserResponse{
		AccessToken: token,
	}
	ctx.JSON(http.StatusCreated, userResponse)

}
