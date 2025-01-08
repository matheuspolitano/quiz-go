package api

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func (svc *Server) listAllTypeQuiz(ctx *gin.Context) {
	allTypes, err := svc.store.ListAllTypes()
	if err != nil {
		SendError(ctx, "error in provide question types", err.Error(), http.StatusBadRequest)
		return
	}
	ctx.JSON(http.StatusAccepted, allTypes)
}
