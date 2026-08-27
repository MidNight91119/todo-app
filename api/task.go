package api

import (
	"errors"
	"log"
	"net/http"
	"strconv"

	db "github.com/MidNight91119/todo-app/db/sqlc"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type createTaskRequest struct {
	Title string `json:"title" binding:"required,max=255"`
}

func (s *Server) createTask(ctx *gin.Context) {
	var req createTaskRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Println("bad request: ", err)
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

func (s *Server) listTasks(ctx *gin.Context) {
	tasks, err := s.queries.ListTasks(ctx.Request.Context())
	if err != nil {
		log.Println("create task failed: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	if tasks == nil {
		tasks = []db.Task{}
	}

	ctx.JSON(http.StatusOK, tasks)
}

func (s *Server) getTask(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Println("bad request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	task, err := s.queries.GetTask(ctx.Request.Context(), id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("task does not exist: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		log.Println("internal server error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

func (s *Server) deleteTask(ctx *gin.Context) {
	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Println("bad request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	err = s.queries.DeleteTask(ctx.Request.Context(), id)
	if err != nil {
		log.Println("internal server error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.Status(http.StatusNoContent)
}

type updateTaskRequest struct {
	Title string `json:"title" binding:"required,max=255"`
}

func (s *Server) updateTask(ctx *gin.Context) {
	var req updateTaskRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Println("bad request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Println("bad request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	arg := db.UpdateTaskParams{
		Title: req.Title,
	}
	arg.ID = id

	task, err := s.queries.UpdateTask(ctx.Request.Context(), arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("task does not exist: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		log.Println("internal server error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}

type completeTaskRequest struct {
	Completed bool `json:"completed"`
}

func (s *Server) completeTask(ctx *gin.Context) {
	var req completeTaskRequest

	err := ctx.ShouldBindJSON(&req)
	if err != nil {
		log.Println("bad request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	id, err := strconv.ParseInt(ctx.Param("id"), 10, 64)
	if err != nil {
		log.Println("bad request: ", err)
		ctx.JSON(http.StatusBadRequest, gin.H{"error": "bad request"})
		return
	}

	arg := db.CompleteTaskParams{
		Completed: req.Completed,
		ID:        id,
	}

	task, err := s.queries.CompleteTask(ctx.Request.Context(), arg)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			log.Println("task does not exist: ", err)
			ctx.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}
		log.Println("internal server error: ", err)
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		return
	}

	ctx.JSON(http.StatusOK, task)
}
