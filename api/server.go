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
	router.POST("/tasks", server.createTask)
	router.GET("/tasks", server.listTasks)
	router.GET("/tasks/:id", server.getTask)
	router.PATCH("/tasks/:id", server.updateTask)
	router.PATCH("/tasks/:id/complete", server.completeTask)
	router.DELETE("/tasks/:id", server.deleteTask)

	return server
}

func (s *Server) Start(address string) error {
	return s.router.Run(address)
}
