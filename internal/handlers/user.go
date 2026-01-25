package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/service"
)

type UserHandler struct {
	service *service.UserService
}

func NewUserHandler(service *service.UserService) *UserHandler {
	return &UserHandler{service: service}
}

// List godoc
// @Summary List all users
// @Description Get a paginated list of all users
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(20) maximum(100)
// @Success 200 {object} models.PaginatedResponse[models.UserResponse]
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [get]
func (h *UserHandler) List(c *gin.Context) {
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

	users, total, err := h.service.List(c.Request.Context(), params.Limit, params.Offset())
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, models.NewPaginatedResponse(users, params.Page, params.Limit, total))
}

// Get godoc
// @Summary Get user by ID
// @Description Get a single user by their ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /api/v1/users/{id} [get]
func (h *UserHandler) Get(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	user, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Create godoc
// @Summary Create a new user
// @Description Create a new user account
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body models.UserCreate true "User data"
// @Success 201 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users [post]
func (h *UserHandler) Create(c *gin.Context) {
	var req models.UserCreate
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest(err.Error()))
		return
	}

	userEmail, _ := c.Get("userEmail")
	createdBy, _ := userEmail.(string)

	user, err := h.service.Create(c.Request.Context(), &req, createdBy)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, user)
}

// Update godoc
// @Summary Update a user
// @Description Update an existing user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body models.UserUpdate true "User data"
// @Success 200 {object} models.UserResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [put]
func (h *UserHandler) Update(c *gin.Context) {
	id, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	var req models.UserUpdate
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest(err.Error()))
		return
	}

	user, err := h.service.Update(c.Request.Context(), id, &req)
	if err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, user)
}

// Delete godoc
// @Summary Delete a user
// @Description Delete a user by ID
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id} [delete]
func (h *UserHandler) Delete(c *gin.Context) {
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

// AddToGroupRequest represents the request body for adding a user to a group
type AddToGroupRequest struct {
	GroupID uint `json:"group_id" binding:"required" example:"1"`
}

// AddToGroup godoc
// @Summary Add user to group
// @Description Add a user to a group. User must be a member of the group's organization.
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body AddToGroupRequest true "Group ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 403 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id}/groups [post]
func (h *UserHandler) AddToGroup(c *gin.Context) {
	userID, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	var req AddToGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest(err.Error()))
		return
	}

	if err := h.service.AddToGroup(c.Request.Context(), userID, req.GroupID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user added to group"})
}

// RemoveFromGroup godoc
// @Summary Remove user from group
// @Description Remove a user from a group
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param gid path int true "Group ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id}/groups/{gid} [delete]
func (h *UserHandler) RemoveFromGroup(c *gin.Context) {
	userID, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	groupID, err := parseID(c, "gid")
	if err != nil {
		respondError(c, err)
		return
	}

	if err := h.service.RemoveFromGroup(c.Request.Context(), userID, groupID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}

// AddToOrganizationRequest represents the request body for adding a user to an organization
type AddToOrganizationRequest struct {
	OrganizationID uint `json:"organization_id" binding:"required" example:"1"`
}

// AddToOrganization godoc
// @Summary Add user to organization
// @Description Add a user to an organization
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param request body AddToOrganizationRequest true "Organization ID"
// @Success 200 {object} MessageResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id}/organizations [post]
func (h *UserHandler) AddToOrganization(c *gin.Context) {
	userID, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	var req AddToOrganizationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		respondError(c, apperror.BadRequest(err.Error()))
		return
	}

	if err := h.service.AddToOrganization(c.Request.Context(), userID, req.OrganizationID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user added to organization"})
}

// RemoveFromOrganization godoc
// @Summary Remove user from organization
// @Description Remove a user from an organization
// @Tags users
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "User ID"
// @Param oid path int true "Organization ID"
// @Success 204 "No Content"
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /api/v1/users/{id}/organizations/{oid} [delete]
func (h *UserHandler) RemoveFromOrganization(c *gin.Context) {
	userID, err := parseID(c, "id")
	if err != nil {
		respondError(c, err)
		return
	}

	orgID, err := parseID(c, "oid")
	if err != nil {
		respondError(c, err)
		return
	}

	if err := h.service.RemoveFromOrganization(c.Request.Context(), userID, orgID); err != nil {
		respondError(c, err)
		return
	}

	c.JSON(http.StatusNoContent, nil)
}
