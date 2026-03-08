package models

import "testing"

// FuzzMatchesAge tests GovernmentFundingProperty.MatchesAge with random age ranges.
func FuzzMatchesAge(f *testing.F) {
	f.Add(3, true, 0, true, 6)   // within range
	f.Add(7, true, 0, true, 6)   // above max
	f.Add(3, false, 0, false, 0) // no bounds
	f.Add(3, true, 5, false, 0)  // below min
	f.Add(0, true, 0, true, 0)   // exact single age

	f.Fuzz(func(t *testing.T, age int, hasMin bool, minAge int, hasMax bool, maxAge int) {
		// Clamp to reasonable values
		age = clamp(age, -10, 200)
		minAge = clamp(minAge, 0, 100)
		maxAge = clamp(maxAge, 0, 100)

		prop := &GovernmentFundingProperty{}
		if hasMin {
			v := minAge
			prop.MinAge = &v
		}
		if hasMax {
			v := maxAge
			prop.MaxAge = &v
		}

		result := prop.MatchesAge(age)

		// Both nil → always true
		if !hasMin && !hasMax {
			if !result {
				t.Errorf("no age bounds should always match, got false for age %d", age)
			}
			return
		}

		// age < minAge → false
		if hasMin && age < minAge && result {
			t.Errorf("age %d < minAge %d should not match", age, minAge)
		}

		// age > maxAge → false
		if hasMax && age > maxAge && result {
			t.Errorf("age %d > maxAge %d should not match", age, maxAge)
		}

		// minAge <= age <= maxAge → true (when both set and range is valid)
		if hasMin && hasMax && minAge <= maxAge && age >= minAge && age <= maxAge && !result {
			t.Errorf("age %d in [%d, %d] should match", age, minAge, maxAge)
		}

		// Only min set, age >= minAge → true
		if hasMin && !hasMax && age >= minAge && !result {
			t.Errorf("age %d >= minAge %d should match (no max)", age, minAge)
		}

		// Only max set, age <= maxAge → true
		if !hasMin && hasMax && age <= maxAge && !result {
			t.Errorf("age %d <= maxAge %d should match (no min)", age, maxAge)
		}
	})
}
