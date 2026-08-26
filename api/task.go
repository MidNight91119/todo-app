package api

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type createTaskRequest struct {
	Title string `json:"title" binding:"required,max=255"`
}

func (s *Server) createTask(ctx *gin.Context) {
	var req createTaskRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	task, err := s.queries.CreateTask(ctx.Request.Context(), req.Title)
	if err != nil {
		log.Println("create task failed: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusCreated, task)
}
