package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	// models imported for swaggo type resolution
	_ "github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/service"
)

// GovernmentFundingBillHandler handles government funding bill upload endpoints.
type GovernmentFundingBillHandler struct {
	service *service.GovernmentFundingBillService
}

// NewGovernmentFundingBillHandler creates a new GovernmentFundingBillHandler.
func NewGovernmentFundingBillHandler(service *service.GovernmentFundingBillService) *GovernmentFundingBillHandler {
	return &GovernmentFundingBillHandler{service: service}
}

// UploadISBJ godoc
// @Summary Upload ISBJ government funding bill
// @Description Parse an ISBJ Senatsabrechnung Excel file and return funding bill data enriched with matched child/contract info
// @Tags government-funding-bills
// @Accept multipart/form-data
// @Produce json
// @Security BearerAuth
// @Param orgId path int true "Organization ID"
// @Param file formData file true "ISBJ Senatsabrechnung Excel file (.xlsx)"
// @Success 200 {object} models.GovernmentFundingBillResponse
// @Failure 400 {object} models.ErrorResponse
// @Failure 401 {object} models.ErrorResponse
// @Failure 500 {object} models.ErrorResponse
// @Router /api/v1/organizations/{orgId}/government-funding-bills/isbj [post]
func (h *GovernmentFundingBillHandler) UploadISBJ(c *gin.Context) {
	orgID, ok := parseOrgID(c)
	if !ok {
		return
	}

	fileHeader, err := c.FormFile("file")
	if err != nil {
		respondError(c, apperror.BadRequest("file is required"))
		return
	}

	file, err := fileHeader.Open()
	if err != nil {
		respondError(c, apperror.BadRequest("failed to read uploaded file"))
		return
	}
	defer file.Close()

	result, err := h.service.ProcessISBJ(c.Request.Context(), orgID, file)
	if err != nil {
		respondError(c, apperror.BadRequest(err.Error()))
		return
	}

	c.JSON(http.StatusOK, result)
}
