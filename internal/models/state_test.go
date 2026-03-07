package models

import "testing"

func TestIsValidState(t *testing.T) {
	tests := []struct {
		name  string
		state string
		want  bool
	}{
		{"valid state berlin", "berlin", true},
		{"invalid state", "hamburg", false},
		{"empty string", "", false},
		{"uppercase", "BERLIN", false},
		{"mixed case", "Berlin", false},
		{"whitespace", " berlin ", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsValidState(tt.state); got != tt.want {
				t.Errorf("IsValidState(%q) = %v, want %v", tt.state, got, tt.want)
			}
		})
	}
}

func TestValidStatesString(t *testing.T) {
	result := ValidStatesString()
	if result == "" {
		t.Error("ValidStatesString() returned empty string")
	}
	if result != "berlin" {
		t.Errorf("ValidStatesString() = %q, want %q", result, "berlin")
	}
}
