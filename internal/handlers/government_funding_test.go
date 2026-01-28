package handlers

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/service"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

func TestGovernmentFundingHandler_CRUD(t *testing.T) {
	db := setupTestDB(t)
	fundingStore := store.NewGovernmentFundingStore(db)
	svc := service.NewGovernmentFundingService(fundingStore)
	handler := NewGovernmentFundingHandler(svc)

	r := setupTestRouter()
	r.GET("/fundings", handler.List)
	r.GET("/fundings/:id", handler.Get)
	r.POST("/fundings", handler.Create)
	r.PUT("/fundings/:id", handler.Update)
	r.DELETE("/fundings/:id", handler.Delete)

	// Test Create
	t.Run("Create", func(t *testing.T) {
		body := models.GovernmentFundingCreateRequest{Name: "Berlin Kita Funding", State: "berlin"}
		w := performRequest(r, "POST", "/fundings", body)

		if w.Code != http.StatusCreated {
			t.Errorf("expected status %d, got %d: %s", http.StatusCreated, w.Code, w.Body.String())
		}

		var result models.GovernmentFunding
		parseResponse(t, w, &result)
		if result.Name != "Berlin Kita Funding" {
			t.Errorf("expected name 'Berlin Kita Funding', got '%s'", result.Name)
		}
		if result.State != "berlin" {
			t.Errorf("expected state 'berlin', got '%s'", result.State)
		}
	})

	// Test List
	t.Run("List", func(t *testing.T) {
		w := performRequest(r, "GET", "/fundings", nil)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response models.PaginatedResponse[models.GovernmentFunding]
		parseResponse(t, w, &response)
		if len(response.Data) != 1 {
			t.Errorf("expected 1 funding, got %d", len(response.Data))
		}
	})

	// Test Get
	t.Run("Get", func(t *testing.T) {
		w := performRequest(r, "GET", "/fundings/1", nil)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d", http.StatusOK, w.Code)
		}

		var result models.GovernmentFunding
		parseResponse(t, w, &result)
		if result.Name != "Berlin Kita Funding" {
			t.Errorf("expected name 'Berlin Kita Funding', got '%s'", result.Name)
		}
	})

	// Test Update
	t.Run("Update", func(t *testing.T) {
		name := "Berlin Updated"
		body := models.GovernmentFundingUpdateRequest{Name: &name}
		w := performRequest(r, "PUT", "/fundings/1", body)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var result models.GovernmentFunding
		parseResponse(t, w, &result)
		if result.Name != "Berlin Updated" {
			t.Errorf("expected name 'Berlin Updated', got '%s'", result.Name)
		}
	})

	// Test Delete
	t.Run("Delete", func(t *testing.T) {
		w := performRequest(r, "DELETE", "/fundings/1", nil)

		if w.Code != http.StatusNoContent {
			t.Errorf("expected status %d, got %d", http.StatusNoContent, w.Code)
		}
	})
}

func TestGovernmentFundingHandler_CreatePeriod_NoOverlap(t *testing.T) {
	db := setupTestDB(t)
	fundingStore := store.NewGovernmentFundingStore(db)
	svc := service.NewGovernmentFundingService(fundingStore)
	handler := NewGovernmentFundingHandler(svc)

	// Create test funding
	funding := &models.GovernmentFunding{Name: "Test Funding", State: "berlin"}
	db.Create(funding)

	r := setupTestRouter()
	r.POST("/fundings/:id/periods", handler.CreatePeriod)

	tests := []struct {
		name           string
		existingFrom   string
		existingTo     *string
		newFrom        string
		newTo          *string
		expectedStatus int
		description    string
	}{
		{
			name:           "no overlap: new period after existing",
			existingFrom:   "2024-01-01",
			existingTo:     strPtr("2024-06-30"),
			newFrom:        "2024-07-01",
			newTo:          strPtr("2024-12-31"),
			expectedStatus: http.StatusCreated,
			description:    "New period starts after existing ends",
		},
		{
			name:           "no overlap: new period before existing",
			existingFrom:   "2024-07-01",
			existingTo:     strPtr("2024-12-31"),
			newFrom:        "2024-01-01",
			newTo:          strPtr("2024-06-30"),
			expectedStatus: http.StatusCreated,
			description:    "New period ends before existing starts",
		},
		{
			name:           "overlap: new period inside existing",
			existingFrom:   "2024-01-01",
			existingTo:     strPtr("2024-12-31"),
			newFrom:        "2024-03-01",
			newTo:          strPtr("2024-06-30"),
			expectedStatus: http.StatusBadRequest,
			description:    "New period is entirely within existing period",
		},
		{
			name:           "overlap: new period spans existing",
			existingFrom:   "2024-03-01",
			existingTo:     strPtr("2024-06-30"),
			newFrom:        "2024-01-01",
			newTo:          strPtr("2024-12-31"),
			expectedStatus: http.StatusBadRequest,
			description:    "New period completely covers existing period",
		},
		{
			name:           "overlap: partial overlap at start",
			existingFrom:   "2024-06-01",
			existingTo:     strPtr("2024-12-31"),
			newFrom:        "2024-01-01",
			newTo:          strPtr("2024-07-31"),
			expectedStatus: http.StatusBadRequest,
			description:    "New period overlaps at the start of existing",
		},
		{
			name:           "overlap: partial overlap at end",
			existingFrom:   "2024-01-01",
			existingTo:     strPtr("2024-06-30"),
			newFrom:        "2024-05-01",
			newTo:          strPtr("2024-12-31"),
			expectedStatus: http.StatusBadRequest,
			description:    "New period overlaps at the end of existing",
		},
		{
			name:           "overlap: existing has no end date",
			existingFrom:   "2024-01-01",
			existingTo:     nil,
			newFrom:        "2024-06-01",
			newTo:          strPtr("2024-12-31"),
			expectedStatus: http.StatusBadRequest,
			description:    "Existing period is ongoing, new period overlaps",
		},
		{
			name:           "overlap: new has no end date",
			existingFrom:   "2024-06-01",
			existingTo:     strPtr("2024-12-31"),
			newFrom:        "2024-01-01",
			newTo:          nil,
			expectedStatus: http.StatusBadRequest,
			description:    "New period is ongoing and overlaps existing",
		},
		{
			name:           "no overlap: new period before ongoing",
			existingFrom:   "2024-07-01",
			existingTo:     nil,
			newFrom:        "2024-01-01",
			newTo:          strPtr("2024-06-30"),
			expectedStatus: http.StatusCreated,
			description:    "New period ends before ongoing period starts",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Clean up periods from previous test
			db.Where("government_funding_id = ?", funding.ID).Delete(&models.GovernmentFundingPeriod{})

			// Create existing period
			existingFrom, _ := time.Parse("2006-01-02", tt.existingFrom)
			existingPeriod := &models.GovernmentFundingPeriod{
				GovernmentFundingID: funding.ID,
				From:                existingFrom,
			}
			if tt.existingTo != nil {
				to, _ := time.Parse("2006-01-02", *tt.existingTo)
				existingPeriod.To = &to
			}
			db.Create(existingPeriod)

			// Try to create new period
			newFrom, _ := time.Parse("2006-01-02", tt.newFrom)
			body := map[string]interface{}{
				"from": newFrom.Format(time.RFC3339),
			}
			if tt.newTo != nil {
				newTo, _ := time.Parse("2006-01-02", *tt.newTo)
				body["to"] = newTo.Format(time.RFC3339)
			}

			w := performRequest(r, "POST", "/fundings/1/periods", body)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d: %s",
					tt.description, tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestGovernmentFundingHandler_UpdatePeriod_NoOverlap(t *testing.T) {
	db := setupTestDB(t)
	fundingStore := store.NewGovernmentFundingStore(db)
	svc := service.NewGovernmentFundingService(fundingStore)
	handler := NewGovernmentFundingHandler(svc)

	// Create test funding
	funding := &models.GovernmentFunding{Name: "Test Funding", State: "berlin"}
	db.Create(funding)

	r := setupTestRouter()
	r.PUT("/fundings/:id/periods/:periodId", handler.UpdatePeriod)

	t.Run("update period to overlap with another", func(t *testing.T) {
		// Clean up
		db.Where("government_funding_id = ?", funding.ID).Delete(&models.GovernmentFundingPeriod{})

		// Create two non-overlapping periods
		from1, _ := time.Parse("2006-01-02", "2024-01-01")
		to1, _ := time.Parse("2006-01-02", "2024-06-30")
		period1 := &models.GovernmentFundingPeriod{GovernmentFundingID: funding.ID, From: from1, To: &to1}
		db.Create(period1)

		from2, _ := time.Parse("2006-01-02", "2024-07-01")
		to2, _ := time.Parse("2006-01-02", "2024-12-31")
		period2 := &models.GovernmentFundingPeriod{GovernmentFundingID: funding.ID, From: from2, To: &to2}
		db.Create(period2)

		// Try to update period2 to overlap with period1
		newFrom, _ := time.Parse("2006-01-02", "2024-05-01")
		body := map[string]interface{}{
			"from": newFrom.Format(time.RFC3339),
		}

		w := performRequest(r, "PUT", "/fundings/1/periods/"+itoa(int(period2.ID)), body)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d for overlapping update, got %d: %s",
				http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("update period without causing overlap", func(t *testing.T) {
		// Clean up
		db.Where("government_funding_id = ?", funding.ID).Delete(&models.GovernmentFundingPeriod{})

		// Create two non-overlapping periods
		from1, _ := time.Parse("2006-01-02", "2024-01-01")
		to1, _ := time.Parse("2006-01-02", "2024-06-30")
		period1 := &models.GovernmentFundingPeriod{GovernmentFundingID: funding.ID, From: from1, To: &to1}
		db.Create(period1)

		from2, _ := time.Parse("2006-01-02", "2024-07-01")
		to2, _ := time.Parse("2006-01-02", "2024-12-31")
		period2 := &models.GovernmentFundingPeriod{GovernmentFundingID: funding.ID, From: from2, To: &to2}
		db.Create(period2)

		// Update period2's end date (no overlap)
		newTo, _ := time.Parse("2006-01-02", "2025-06-30")
		body := map[string]interface{}{
			"to": newTo.Format(time.RFC3339),
		}

		w := performRequest(r, "PUT", "/fundings/1/periods/"+itoa(int(period2.ID)), body)

		if w.Code != http.StatusOK {
			t.Errorf("expected status %d for valid update, got %d: %s",
				http.StatusOK, w.Code, w.Body.String())
		}
	})
}

func TestGovernmentFundingHandler_Property_AgeRange(t *testing.T) {
	db := setupTestDB(t)
	fundingStore := store.NewGovernmentFundingStore(db)
	svc := service.NewGovernmentFundingService(fundingStore)
	handler := NewGovernmentFundingHandler(svc)

	// Create test funding and period
	funding := &models.GovernmentFunding{Name: "Test Funding", State: "berlin"}
	db.Create(funding)
	period := &models.GovernmentFundingPeriod{
		GovernmentFundingID: funding.ID,
		From:                time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	db.Create(period)

	r := setupTestRouter()
	r.POST("/fundings/:id/periods/:periodId/properties", handler.CreateProperty)

	tests := []struct {
		name           string
		minAge         *int
		maxAge         *int
		expectedStatus int
		description    string
	}{
		{
			name:           "valid range 0-2",
			minAge:         intPtr(0),
			maxAge:         intPtr(2),
			expectedStatus: http.StatusCreated,
			description:    "Children from birth up to but not including 2nd birthday",
		},
		{
			name:           "valid range 3-7",
			minAge:         intPtr(3),
			maxAge:         intPtr(7),
			expectedStatus: http.StatusCreated,
			description:    "Children from 3rd birthday up to but not including 7th birthday",
		},
		{
			name:           "no age filter",
			minAge:         nil,
			maxAge:         nil,
			expectedStatus: http.StatusCreated,
			description:    "Property applies to all ages",
		},
		{
			name:           "invalid: min equals max",
			minAge:         intPtr(2),
			maxAge:         intPtr(2),
			expectedStatus: http.StatusBadRequest,
			description:    "Empty range - no children would qualify",
		},
		{
			name:           "invalid: min greater than max",
			minAge:         intPtr(5),
			maxAge:         intPtr(3),
			expectedStatus: http.StatusBadRequest,
			description:    "Inverted range is invalid",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := map[string]interface{}{
				"name":        "test-property",
				"payment":     100000,
				"requirement": 0.1,
			}
			if tt.minAge != nil {
				body["min_age"] = *tt.minAge
			}
			if tt.maxAge != nil {
				body["max_age"] = *tt.maxAge
			}

			w := performRequest(r, "POST", fmt.Sprintf("/fundings/%d/periods/%d/properties", funding.ID, period.ID), body)

			if w.Code != tt.expectedStatus {
				t.Errorf("%s: expected status %d, got %d: %s",
					tt.description, tt.expectedStatus, w.Code, w.Body.String())
			}
		})
	}
}

func TestGovernmentFundingHandler_Get_PeriodsLimit(t *testing.T) {
	db := setupTestDB(t)
	fundingStore := store.NewGovernmentFundingStore(db)
	svc := service.NewGovernmentFundingService(fundingStore)
	handler := NewGovernmentFundingHandler(svc)

	// Create test funding with multiple periods
	funding := &models.GovernmentFunding{Name: "Test Funding", State: "berlin"}
	db.Create(funding)

	// Create 3 periods (oldest to newest)
	for i := 0; i < 3; i++ {
		from := time.Date(2024, time.Month(i*4+1), 1, 0, 0, 0, 0, time.UTC)
		to := time.Date(2024, time.Month(i*4+4), 1, 0, 0, 0, 0, time.UTC)
		period := &models.GovernmentFundingPeriod{
			GovernmentFundingID: funding.ID,
			From:                from,
			To:                  &to,
		}
		db.Create(period)
	}

	r := setupTestRouter()
	r.GET("/fundings/:id", handler.Get)

	t.Run("default returns only latest period", func(t *testing.T) {
		w := performRequest(r, "GET", fmt.Sprintf("/fundings/%d", funding.ID), nil)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var result service.GovernmentFundingWithDetailsResponse
		parseResponse(t, w, &result)

		if len(result.Periods) != 1 {
			t.Errorf("expected 1 period by default, got %d", len(result.Periods))
		}
		if result.TotalPeriods != 3 {
			t.Errorf("expected total_periods=3, got %d", result.TotalPeriods)
		}
	})

	t.Run("periods_limit=0 returns all periods", func(t *testing.T) {
		w := performRequest(r, "GET", fmt.Sprintf("/fundings/%d?periods_limit=0", funding.ID), nil)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var result service.GovernmentFundingWithDetailsResponse
		parseResponse(t, w, &result)

		if len(result.Periods) != 3 {
			t.Errorf("expected 3 periods with limit=0, got %d", len(result.Periods))
		}
		if result.TotalPeriods != 3 {
			t.Errorf("expected total_periods=3, got %d", result.TotalPeriods)
		}
	})

	t.Run("periods_limit=2 returns 2 periods", func(t *testing.T) {
		w := performRequest(r, "GET", fmt.Sprintf("/fundings/%d?periods_limit=2", funding.ID), nil)

		if w.Code != http.StatusOK {
			t.Fatalf("expected status %d, got %d: %s", http.StatusOK, w.Code, w.Body.String())
		}

		var result service.GovernmentFundingWithDetailsResponse
		parseResponse(t, w, &result)

		if len(result.Periods) != 2 {
			t.Errorf("expected 2 periods with limit=2, got %d", len(result.Periods))
		}
		if result.TotalPeriods != 3 {
			t.Errorf("expected total_periods=3, got %d", result.TotalPeriods)
		}
	})

	t.Run("negative periods_limit returns 400", func(t *testing.T) {
		w := performRequest(r, "GET", fmt.Sprintf("/fundings/%d?periods_limit=-1", funding.ID), nil)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d for negative limit, got %d: %s",
				http.StatusBadRequest, w.Code, w.Body.String())
		}
	})

	t.Run("invalid periods_limit returns 400", func(t *testing.T) {
		w := performRequest(r, "GET", fmt.Sprintf("/fundings/%d?periods_limit=abc", funding.ID), nil)

		if w.Code != http.StatusBadRequest {
			t.Errorf("expected status %d for invalid limit, got %d: %s",
				http.StatusBadRequest, w.Code, w.Body.String())
		}
	})
}

func itoa(i int) string {
	return fmt.Sprintf("%d", i)
}

func intPtr(i int) *int {
	return &i
}

func strPtr(s string) *string {
	return &s
}
