package models

import "testing"

// FuzzContractPropertiesHasValue tests HasValue with random key/value combinations.
func FuzzContractPropertiesHasValue(f *testing.F) {
	f.Add("care_type", "ganztag", "care_type", "ganztag", true)   // scalar match
	f.Add("care_type", "ganztag", "care_type", "halbzeit", false) // scalar no match
	f.Add("supp", "ndh", "supp", "ndh", true)                     // array match
	f.Add("supp", "ndh", "supp", "mss", false)                    // array miss
	f.Add("key1", "val1", "other", "val1", false)                 // wrong key

	f.Fuzz(func(t *testing.T, setKey, setValue, queryKey, queryValue string, useArray bool) {
		// Build properties with either scalar or array value
		props := ContractProperties{}
		if useArray {
			props[setKey] = []string{setValue}
		} else {
			props[setKey] = setValue
		}

		result := props.HasValue(queryKey, queryValue)

		if queryKey != setKey {
			// Different key → always false
			if result {
				t.Errorf("HasValue(%q, %q) should be false for key %q", queryKey, queryValue, setKey)
			}
			return
		}

		// Same key — check value match
		if queryValue == setValue && !result {
			t.Errorf("HasValue(%q, %q) should be true (stored as array=%v)", queryKey, queryValue, useArray)
		}
		if queryValue != setValue && result {
			t.Errorf("HasValue(%q, %q) should be false when stored value is %q", queryKey, queryValue, setValue)
		}

		// Nil properties → always false
		var nilProps ContractProperties
		if nilProps.HasValue(queryKey, queryValue) {
			t.Error("nil properties should always return false")
		}
	})
}

// FuzzContractPropertiesMergeDefaults tests MergeDefaults merge precedence.
func FuzzContractPropertiesMergeDefaults(f *testing.F) {
	f.Add("k1", "orig", "k1", "default") // same key, original wins
	f.Add("k1", "orig", "k2", "default") // different keys, both present
	f.Add("", "orig", "", "default")     // empty keys

	f.Fuzz(func(t *testing.T, origKey, origVal, defKey, defVal string) {
		original := ContractProperties{origKey: origVal}
		defaults := ContractProperties{defKey: defVal}

		merged := original.MergeDefaults(defaults)

		// Original keys are never overwritten
		if v, ok := merged[origKey]; ok {
			if str, ok := v.(string); ok && str != origVal {
				t.Errorf("original key %q was overwritten: got %q, want %q", origKey, str, origVal)
			}
		}

		// Default keys not in original are present after merge
		if origKey != defKey {
			if v, ok := merged[defKey]; !ok {
				t.Errorf("default key %q missing after merge", defKey)
			} else if str, ok := v.(string); ok && str != defVal {
				t.Errorf("default key %q has wrong value: got %q, want %q", defKey, str, defVal)
			}
		}

		// MergeDefaults(nil) returns original unchanged
		mergedNil := original.MergeDefaults(nil)
		for k, v := range original {
			if mergedNil[k] != v {
				t.Errorf("MergeDefaults(nil) changed key %q", k)
			}
		}

		// nil.MergeDefaults(defaults) contains all defaults
		var nilProps ContractProperties
		mergedFromNil := nilProps.MergeDefaults(defaults)
		for k, v := range defaults {
			if mergedFromNil[k] != v {
				t.Errorf("nil.MergeDefaults missing key %q", k)
			}
		}
	})
}
