package api

import (
	"net/http"

	db "github.com/MidNight91119/todo-app/db/sqlc"
	"github.com/gin-gonic/gin"
)

type Server struct {
	queries *db.Queries
	router  *gin.Engine
}

func NewServer(queries *db.Queries) *Server {
	server := &Server{
		queries: queries,
	}

	router := gin.Default()
	server.router = router

	router.GET("/ping", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
