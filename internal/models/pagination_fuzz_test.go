package models

import "testing"

// FuzzPaginationValidateAndDefaults tests PaginationParams validation and defaults.
func FuzzPaginationValidateAndDefaults(f *testing.F) {
	f.Add(1, 20)
	f.Add(0, 0)
	f.Add(-1, 10)
	f.Add(1, -5)
	f.Add(1, 101)
	f.Add(5, 50)

	f.Fuzz(func(t *testing.T, page, limit int) {
		p := PaginationParams{Page: page, Limit: limit}

		err := p.Validate()

		// Negative page → error
		if page < 0 && err == nil {
			t.Errorf("negative page %d should fail validation", page)
		}
		// Negative limit → error
		if limit < 0 && err == nil {
			t.Errorf("negative limit %d should fail validation", limit)
		}
		// Limit > 100 → error
		if limit > 100 && err == nil {
			t.Errorf("limit %d > 100 should fail validation", limit)
		}
		// Valid input → no error
		if page >= 0 && limit >= 0 && limit <= 100 && err != nil {
			t.Errorf("valid params (page=%d, limit=%d) should not error, got: %v", page, limit, err)
		}

		// After SetDefaults: page >= 1, limit >= 1
		p.SetDefaults()
		if p.Page < 1 {
			t.Errorf("after SetDefaults, page should be >= 1, got %d", p.Page)
		}
		if p.Limit < 1 {
			t.Errorf("after SetDefaults, limit should be >= 1, got %d", p.Limit)
		}

		// Offset == (page-1) * limit
		expectedOffset := (p.Page - 1) * p.Limit
		if p.Offset() != expectedOffset {
			t.Errorf("Offset() = %d, expected %d (page=%d, limit=%d)", p.Offset(), expectedOffset, p.Page, p.Limit)
		}
	})
}
