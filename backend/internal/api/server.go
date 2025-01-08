package api

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/matheuspolitano/quiz-go/backend/internal/config"
	"github.com/matheuspolitano/quiz-go/backend/internal/memdb"
	"github.com/matheuspolitano/quiz-go/backend/internal/token"
)

// Server api type
type Server struct {
	router     *gin.Engine
	config     config.Config
	httpSvc    *http.Server
	store      *memdb.DBManager
	tokenMaker token.Maker
}

// New create new server
func New(config config.Config, store *memdb.DBManager) *Server {
	router := gin.Default()
	svc := &Server{router: router, config: config, tokenMaker: &token.JWTMaker{}, store: store}
	return svc.WithRoutes().WithServer()
}

// WithRoutes implement the routes
func (svc *Server) WithRoutes() *Server {
	apiGroup := svc.router.Group("/api")
	apiGroup.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusAccepted, gin.H{
			"message": "pong",
		})
	})
	apiGroup.POST("/login", svc.startUser)
	authRoutes := apiGroup.Group("/quiz").Use(authMiddleware(svc.tokenMaker))
	authRoutes.GET("/ping", func(ctx *gin.Context) {
		authPayload := ctx.MustGet(authorizationPayloadKey).(*token.Payload)
		ctx.JSON(http.StatusAccepted, gin.H{
			"message": authPayload.Username,
		})
	})
	authRoutes.GET("/types", svc.listAllTypeQuiz)
	authRoutes.GET("/question/:questionID", svc.getQuestion)
	authRoutes.GET("/joinQuiz/:typeQuiz", svc.joinQuiz)
	authRoutes.GET("/answer/:typeQuiz/next", svc.nextQuestion)
	authRoutes.POST("/answer/:typeQuiz/:questionID", svc.answerQuestion)
	authRoutes.GET("/answer/:typeQuiz/score", svc.generalScore)
	return svc
}

// WithServer add http server
func (svc *Server) WithServer() *Server {
	httpSvc := &http.Server{
		Addr:              fmt.Sprintf(":%s", svc.config.ApiPort),
		Handler:           svc.router,
		ReadTimeout:       time.Second * 10,
		ReadHeaderTimeout: time.Second * 5,
		IdleTimeout:       time.Second * 10,
		WriteTimeout:      time.Second * 5,
	}
	svc.httpSvc = httpSvc
	return svc
}

// Start the server
func (svc *Server) Start() error {
	return svc.httpSvc.ListenAndServe()
}

// Shutdown handles server shutdown
func (svc *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(svc.config.ApiTimeShutdown))
	defer cancel()
	if err := svc.httpSvc.Shutdown(ctx); err != nil {
		return err
	}
	return svc.httpSvc.Shutdown(ctx)
}
