package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/service"
)

type ChildHandler struct {
	service *service.ChildService
}

func NewChildHandler(service *service.ChildService) *ChildHandler {
	return &ChildHandler{service: service}
}

// List godoc
// @Summary List all children
// @Description Get a paginated list of all children
// @Tags children
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20) maximum(100)
// @Success 200 {object} models.PaginatedResponse[models.Child]
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/children [get]
func (h *ChildHandler) List(c *gin.Context) {
	var params models.PaginationParams
	if err := c.ShouldBindQuery(&params); err != nil {
		respondError(c, apperror.BadRequest("invalid pagination parameters"))
		return
	}
	if err := params.Validate(); err != nil {
		respondError(c, apperror.BadRequest(err.Error()))
		return
	}
	params.SetDefaults()

	children, total, err := h.service.List(c.Request.Context(), params.Limit, params.Offset())
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.NewPaginatedResponse(children, params.Page, params.Limit, total))
}

// Get godoc
// @Summary Get child by ID
// @Description Get a single child by their ID
// @Tags children
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Child ID"
// @Success 200 {object} models.Child
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/children/{id} [get]
func (h *ChildHandler) Get(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	child, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, child)
}

// Create godoc
// @Summary Create a new child
// @Description Create a new child
// @Tags children
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.ChildCreate true "Child data"
// @Success 201 {object} models.Child
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/children [post]
func (h *ChildHandler) Create(c *gin.Context) {
	var req models.ChildCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest(err.Error()))
		return
	}

	child, err := h.service.Create(c.Request.Context(), &req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, child)
}

// Update godoc
// @Summary Update a child
// @Description Update an existing child by ID
// @Tags children
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Child ID"
// @Param request body models.ChildUpdate true "Child data"
// @Success 200 {object} models.Child
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/children/{id} [put]
func (h *ChildHandler) Update(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	var req models.ChildUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest(err.Error()))
		return
	}

	child, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, child)
}

// Delete godoc
// @Summary Delete a child
// @Description Delete a child by ID
// @Tags children
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Child ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/children/{id} [delete]
func (h *ChildHandler) Delete(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// ListContracts godoc
// @Summary List child contracts
// @Description Get all contracts for a child
// @Tags children
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Child ID"
// @Success 200 {array} models.ChildContract
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/children/{id}/contracts [get]
func (h *ChildHandler) ListContracts(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	contracts, err := h.service.ListContracts(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contracts)
}

// GetCurrentContract godoc
// @Summary Get current child contract
// @Description Get the currently active contract for a child
// @Tags children
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Child ID"
// @Success 200 {object} models.ChildContract
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/children/{id}/contracts/current [get]
func (h *ChildHandler) GetCurrentContract(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	contract, err := h.service.GetCurrentContract(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, contract)
}

// CreateContract godoc
// @Summary Create child contract
// @Description Create a new contract for a child
// @Tags children
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Child ID"
// @Param request body models.ChildContractCreate true "Contract data"
// @Success 201 {object} models.ChildContract
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse "Contract overlaps with existing"
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/children/{id}/contracts [post]
func (h *ChildHandler) CreateContract(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	var req models.ChildContractCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest(err.Error()))
		return
	}

	contract, err := h.service.CreateContract(c.Request.Context(), id, &req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, contract)
}

// DeleteContract godoc
// @Summary Delete child contract
// @Description Delete a contract by ID
// @Tags children
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Child ID"
// @Param contractId path int true "Contract ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/children/{id}/contracts/{contractId} [delete]
func (h *ChildHandler) DeleteContract(c *gin.Context) {
	contractID, err := parseID(c, "contractId")
	if err != nil {
		respondError(c, err)
		return
	}

	if err := h.service.DeleteContract(c.Request.Context(), contractID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
