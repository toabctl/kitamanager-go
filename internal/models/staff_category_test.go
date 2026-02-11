package models

import "testing"

func TestIsValidStaffCategory(t *testing.T) {
	validCategories := []string{"qualified", "supplementary", "non_pedagogical"}
	for _, cat := range validCategories {
		if !IsValidStaffCategory(cat) {
			t.Errorf("IsValidStaffCategory(%q) = false, want true", cat)
		}
	}
}

func TestIsValidStaffCategory_Invalid(t *testing.T) {
	invalidCategories := []string{"", "invalid", "teacher", "QUALIFIED", "Qualified", "fachkraft", "ergaenzungskraft", "non-pedagogical", "non pedagogical"}
	for _, cat := range invalidCategories {
		if IsValidStaffCategory(cat) {
			t.Errorf("IsValidStaffCategory(%q) = true, want false", cat)
		}
	}
}

func TestValidStaffCategoryStrings(t *testing.T) {
	strings := ValidStaffCategoryStrings()
	if len(strings) != 3 {
		t.Errorf("ValidStaffCategoryStrings() returned %d items, want 3", len(strings))
	}
	expected := map[string]bool{"qualified": true, "supplementary": true, "non_pedagogical": true}
	for _, s := range strings {
		if !expected[s] {
			t.Errorf("unexpected staff category string: %q", s)
		}
	}
}

func TestStaffCategory_Constants(t *testing.T) {
	if string(StaffCategoryQualified) != "qualified" {
		t.Errorf("StaffCategoryQualified = %q, want %q", StaffCategoryQualified, "qualified")
	}
	if string(StaffCategorySupplementary) != "supplementary" {
		t.Errorf("StaffCategorySupplementary = %q, want %q", StaffCategorySupplementary, "supplementary")
	}
	if string(StaffCategoryNonPedagogical) != "non_pedagogical" {
		t.Errorf("StaffCategoryNonPedagogical = %q, want %q", StaffCategoryNonPedagogical, "non_pedagogical")
	}
}
