package models

import (
	"testing"
	"time"
)

// FuzzPeriodOverlaps tests the Period.Overlaps method with random inputs
func FuzzPeriodOverlaps(f *testing.F) {
	// Add seed corpus
	now := time.Now()
	seeds := []struct {
		aFromYear, aFromMonth, aFromDay int
		aToYear, aToMonth, aToDay       int
		aHasTo                          bool
		bFromYear, bFromMonth, bFromDay int
		bToYear, bToMonth, bToDay       int
		bHasTo                          bool
	}{
		{2024, 1, 1, 2024, 12, 31, true, 2024, 6, 1, 2024, 6, 30, true},
		{2024, 1, 1, 2024, 12, 31, true, 2025, 1, 1, 2025, 12, 31, true},
		{2024, 1, 1, 0, 0, 0, false, 2024, 6, 1, 2024, 6, 30, true},
		{2024, 1, 1, 0, 0, 0, false, 2023, 1, 1, 2023, 12, 31, true},
		{now.Year(), int(now.Month()), now.Day(), 0, 0, 0, false, now.Year(), int(now.Month()), now.Day(), 0, 0, 0, false},
	}

	for _, s := range seeds {
		f.Add(
			s.aFromYear, s.aFromMonth, s.aFromDay,
			s.aToYear, s.aToMonth, s.aToDay, s.aHasTo,
			s.bFromYear, s.bFromMonth, s.bFromDay,
			s.bToYear, s.bToMonth, s.bToDay, s.bHasTo,
		)
	}

	f.Fuzz(func(t *testing.T,
		aFromYear, aFromMonth, aFromDay int,
		aToYear, aToMonth, aToDay int, aHasTo bool,
		bFromYear, bFromMonth, bFromDay int,
		bToYear, bToMonth, bToDay int, bHasTo bool,
	) {
		// Sanitize inputs to valid date ranges
		aFromYear = clamp(aFromYear, 1900, 2100)
		aFromMonth = clamp(aFromMonth, 1, 12)
		aFromDay = clamp(aFromDay, 1, 28)
		aToYear = clamp(aToYear, 1900, 2100)
		aToMonth = clamp(aToMonth, 1, 12)
		aToDay = clamp(aToDay, 1, 28)

		bFromYear = clamp(bFromYear, 1900, 2100)
		bFromMonth = clamp(bFromMonth, 1, 12)
		bFromDay = clamp(bFromDay, 1, 28)
		bToYear = clamp(bToYear, 1900, 2100)
		bToMonth = clamp(bToMonth, 1, 12)
		bToDay = clamp(bToDay, 1, 28)

		aFrom := time.Date(aFromYear, time.Month(aFromMonth), aFromDay, 0, 0, 0, 0, time.UTC)
		var aTo *time.Time
		if aHasTo {
			to := time.Date(aToYear, time.Month(aToMonth), aToDay, 0, 0, 0, 0, time.UTC)
			// Ensure To is after From
			if to.Before(aFrom) {
				to = aFrom.AddDate(0, 1, 0)
			}
			aTo = &to
		}

		bFrom := time.Date(bFromYear, time.Month(bFromMonth), bFromDay, 0, 0, 0, 0, time.UTC)
		var bTo *time.Time
		if bHasTo {
			to := time.Date(bToYear, time.Month(bToMonth), bToDay, 0, 0, 0, 0, time.UTC)
			// Ensure To is after From
			if to.Before(bFrom) {
				to = bFrom.AddDate(0, 1, 0)
			}
			bTo = &to
		}

		periodA := Period{From: aFrom, To: aTo}
		periodB := Period{From: bFrom, To: bTo}

		// Test that Overlaps is symmetric
		overlapAB := periodA.Overlaps(periodB)
		overlapBA := periodB.Overlaps(periodA)

		if overlapAB != overlapBA {
			t.Errorf("Overlaps is not symmetric: A.Overlaps(B)=%v, B.Overlaps(A)=%v\nA: %+v\nB: %+v",
				overlapAB, overlapBA, periodA, periodB)
		}

		// A period should always overlap with itself
		if !periodA.Overlaps(periodA) {
			t.Errorf("Period should overlap with itself: %+v", periodA)
		}
	})
}

// FuzzPeriodIsActiveOn tests the Period.IsActiveOn method with random inputs
func FuzzPeriodIsActiveOn(f *testing.F) {
	// Add seed corpus
	f.Add(2024, 1, 1, 2024, 12, 31, true, 2024, 6, 15)
	f.Add(2024, 1, 1, 0, 0, 0, false, 2024, 6, 15)
	f.Add(2024, 1, 1, 2024, 1, 1, true, 2024, 1, 1) // Single day period

	f.Fuzz(func(t *testing.T,
		fromYear, fromMonth, fromDay int,
		toYear, toMonth, toDay int, hasTo bool,
		checkYear, checkMonth, checkDay int,
	) {
		// Sanitize inputs
		fromYear = clamp(fromYear, 1900, 2100)
		fromMonth = clamp(fromMonth, 1, 12)
		fromDay = clamp(fromDay, 1, 28)
		toYear = clamp(toYear, 1900, 2100)
		toMonth = clamp(toMonth, 1, 12)
		toDay = clamp(toDay, 1, 28)
		checkYear = clamp(checkYear, 1900, 2100)
		checkMonth = clamp(checkMonth, 1, 12)
		checkDay = clamp(checkDay, 1, 28)

		from := time.Date(fromYear, time.Month(fromMonth), fromDay, 0, 0, 0, 0, time.UTC)
		var to *time.Time
		if hasTo {
			toDate := time.Date(toYear, time.Month(toMonth), toDay, 0, 0, 0, 0, time.UTC)
			// Ensure To is after or equal to From
			if toDate.Before(from) {
				toDate = from
			}
			to = &toDate
		}

		checkDate := time.Date(checkYear, time.Month(checkMonth), checkDay, 0, 0, 0, 0, time.UTC)
		period := Period{From: from, To: to}

		isActive := period.IsActiveOn(checkDate)

		// Verify the result makes sense
		if checkDate.Before(from) && isActive {
			t.Errorf("Date before From should not be active: date=%v, period=%+v", checkDate, period)
		}

		if to != nil && checkDate.After(*to) && isActive {
			t.Errorf("Date after To should not be active: date=%v, period=%+v", checkDate, period)
		}

		// If ongoing (no To), any date >= From should be active
		if to == nil && !checkDate.Before(from) && !isActive {
			t.Errorf("Ongoing period should be active for date >= From: date=%v, period=%+v", checkDate, period)
		}
	})
}

func clamp(val, min, max int) int {
	if val < min {
		return min
	}
	if val > max {
		return max
	}
	return val
}
