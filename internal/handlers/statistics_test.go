package handlers

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

func TestStatisticsHandler_GetStaffingHours_Success(t *testing.T) {
	db := setupTestDB(t)
	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	svc := createStatisticsService(db)
	handler := NewStatisticsHandler(svc)

	r := setupTestRouter()
	r.GET("/organizations/:orgId/statistics/staffing-hours", handler.GetStaffingHours)

	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/statistics/staffing-hours?from=2024-01-01&to=2024-03-01", org.ID), nil)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var response models.StaffingHoursResponse
	parseResponse(t, w, &response)
	if len(response.DataPoints) == 0 {
		t.Error("expected data points")
	}
}

func TestStatisticsHandler_GetStaffingHours_WithQueryParams(t *testing.T) {
	db := setupTestDB(t)
	org := createTestOrganization(t, db, "Test Org")
	db.Model(org).Update("state", "berlin")

	svc := createStatisticsService(db)
	handler := NewStatisticsHandler(svc)

	r := setupTestRouter()
	r.GET("/organizations/:orgId/statistics/staffing-hours", handler.GetStaffingHours)

	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/statistics/staffing-hours?from=2024-06-01&to=2024-09-01&section_id=1", org.ID), nil)
	if w.Code != http.StatusOK {
		t.Errorf("expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var response models.StaffingHoursResponse
	parseResponse(t, w, &response)
	// With custom from/to spanning 4 months (Jun, Jul, Aug, Sep), expect 4 data points
	if len(response.DataPoints) != 4 {
		t.Errorf("expected 4 data points, got %d", len(response.DataPoints))
	}
}

func TestStatisticsHandler_GetStaffingHours_InvalidOrgId(t *testing.T) {
	db := setupTestDB(t)

	svc := createStatisticsService(db)
	handler := NewStatisticsHandler(svc)

	r := setupTestRouter()
	r.GET("/organizations/:orgId/statistics/staffing-hours", handler.GetStaffingHours)

	w := performRequest(r, "GET", "/organizations/abc/statistics/staffing-hours", nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400, got %d: %s", w.Code, w.Body.String())
	}
}

func TestStatisticsHandler_GetStaffingHours_InvalidDates(t *testing.T) {
	db := setupTestDB(t)
	org := createTestOrganization(t, db, "Test Org")

	svc := createStatisticsService(db)
	handler := NewStatisticsHandler(svc)

	r := setupTestRouter()
	r.GET("/organizations/:orgId/statistics/staffing-hours", handler.GetStaffingHours)

	// Invalid from date
	w := performRequest(r, "GET", fmt.Sprintf("/organizations/%d/statistics/staffing-hours?from=not-a-date", org.ID), nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid from date, got %d: %s", w.Code, w.Body.String())
	}

	// Invalid to date
	w = performRequest(r, "GET", fmt.Sprintf("/organizations/%d/statistics/staffing-hours?to=2024-13-99", org.ID), nil)
	if w.Code != http.StatusBadRequest {
		t.Errorf("expected 400 for invalid to date, got %d: %s", w.Code, w.Body.String())
	}
}
