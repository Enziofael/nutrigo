package handler

import (
	"errors"
	"net/http"

	repository "github.com/Enziofael/nutrigo/backend/internal/repository/v1"
	service "github.com/Enziofael/nutrigo/backend/internal/service/v1"
	models "github.com/Enziofael/nutrigo/shared/models/v1"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{
		service: service,
	}
}

// Users

func (h *UserHandler) Create(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query: " + err.Error()})
		return
	}

	include, err := req.GetIncludeQuery().Parse(req.Include)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid include parameter: " + err.Error()})
		return
	}

	resp, err := h.service.Create(c.Request.Context(), req, include)
	if err != nil {
		switch {
		//Validation error
		case errors.Is(err, models.ErrValidation):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//Already exists
		case errors.Is(err, repository.ErrAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		//Internal
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, resp)
}

func (h *UserHandler) List(c *gin.Context) {
	var req models.UserListRequest
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query: " + err.Error()})
		return
	}

	include, err := req.GetIncludeQuery().Parse(req.Include)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid include parameter: " + err.Error()})
		return
	}

	resp, err := h.service.List(c.Request.Context(), req, include)
	if err != nil {
		switch {
		//Validation error
		case errors.Is(err, models.ErrValidation):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//Internal
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Get(c *gin.Context) {
	var req models.UserGetRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uri: " + err.Error()})
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query: " + err.Error()})
		return
	}

	include, err := req.GetIncludeQuery().Parse(req.Include)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid include parameter: " + err.Error()})
		return
	}

	resp, err := h.service.Get(c.Request.Context(), req, include)
	if err != nil {
		switch {
		//Validation error
		case errors.Is(err, models.ErrValidation):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//Not found
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		//Internal
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, resp)
}

func (h *UserHandler) Patch(c *gin.Context) {
	var req models.UserPatchRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uri: " + err.Error()})
		return
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}
	if err := c.ShouldBindQuery(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid query: " + err.Error()})
		return
	}

	include, err := req.GetIncludeQuery().Parse(req.Include)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid include parameter: " + err.Error()})
		return
	}

	res, err := h.service.Patch(c.Request.Context(), req, include)
	if err != nil {
		switch {
		//Validation error
		case errors.Is(err, models.ErrValidation):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		//Not found
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		//Conflict
		case errors.Is(err, repository.ErrRequestConflict):
			c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		//Internal
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusOK, res)
}

func (h *UserHandler) Delete(c *gin.Context) {
	var req models.UserDeleteRequest
	if err := c.ShouldBindUri(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid uri: " + err.Error()})
		return
	}

	err := h.service.Delete(c.Request.Context(), req)
	if err != nil {
		switch {
		case errors.Is(err, models.ErrValidation):
			//Ошибки валидации
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

/*

// =========================================================
// CRUD HANDLERS
// =========================================================

func (h *UserHandler) Create(c *gin.Context) {
	var req models.UserCreateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	user, err := h.service.Create(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrAlreadyExists):
			c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) Delete(c *gin.Context) {
	tgID, err := parseTgID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := models.UserDeleteRequest{TgID: tgID}
	err = h.service.Delete(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.Status(http.StatusNoContent)
}

func (h *UserHandler) Get(c *gin.Context) {
	tgID, err := parseTgID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	req := models.UserGetRequest{TgID: tgID}
	user, err := h.service.Get(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) Patch(c *gin.Context) {
	var req models.UserPatchRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	tgID, err := parseTgID(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.TgID = tgID

	user, err := h.service.Patch(c.Request.Context(), req)

	if err != nil {
		switch {
		case errors.Is(err, repository.ErrNotFound):
			c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		case errors.Is(err, repository.ErrInvalidRequest):
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, user)
}

*/
