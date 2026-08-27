package handler

import (
	service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
)

type ExerciseEntryHandler struct {
	service *service.ExerciseEntryService
}

func NewExerciseEntryHandler(service *service.ExerciseEntryService) *ExerciseEntryHandler {
	return &ExerciseEntryHandler{service: service}
}

// ExerciseEntries
/*

func (h *ExerciseHandler) Create(c *gin.Context) {
	var req models.ExerciseEntryCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	res, err := service.EntryCreate(req)
	if err != nil {
		switch {
		case errors.Is(err, ):
			c.JSON(http.Stratus, gin.H{"error": })
		}
	}

	c.JSON(http.StatusCreated, res)
}

func (h *ExerciseHandler) List(c *gin.Context) {
	var req models.ExerciseEntryListRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	res, err := service.EntryList(req)
	if err != nil {
		switch {
		case errors.Is(err, ):
			c.JSON(http.Stratus, gin.H{"error": })
		}
	}

	c.JSON(http.StatusOK, res)
}

func (h *ExerciseEntryHandler) Get(c *gin.Context) {
	var req models.ExerciseEntryGetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	res, err == service.Get(req)
	if err != nil {
		switch {
		case errors.Is(err, ):
			c.JSON(http.Stratus, gin.H{"error": })
		}
	}

	c.JSON(http.StatusOK, res)
}

func (h *ExerciseEntryHandler) Patch(c *gin.Context) {
	var req models.ExerciseEntryPatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	res, err == service.Patch(req)
	if err != nil {
		switch {
		case errors.Is(err, ):
			c.JSON(http.Stratus, gin.H{"error": })
		}
	}

	c.JSON(http.StatusOK, res)
}

func (h *ExerciseEntryHandler) Delete(c *gin.Context) {
	var req models.ExerciseEntryDeleteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	err == service.Delete(req)
	if err != nil {
		switch {
		case errors.Is(err, ):
			c.JSON(http.Stratus, gin.H{"error": })
		}
	}

	c.JSON(http.StatusNoContent)
}*/
