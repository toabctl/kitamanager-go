package validation

import (
	"fmt"
	"html"
	"strconv"
	"strings"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/models"
)

// IsWhitespaceOnly returns true if string is empty or contains only whitespace
func IsWhitespaceOnly(s string) bool {
	return strings.TrimSpace(s) == ""
}

// ValidateBirthdate ensures birthdate is not in the future
func ValidateBirthdate(birthdate time.Time) error {
	if birthdate.After(time.Now()) {
		return fmt.Errorf("birthdate cannot be in the future")
	}
	return nil
}

// ValidatePeriod ensures From < To when To is provided (contracts must have positive duration)
func ValidatePeriod(from time.Time, to *time.Time) error {
	if to != nil {
		if from.After(*to) {
			return fmt.Errorf("from date must be before to date")
		}
		if from.Equal(*to) {
			return fmt.Errorf("contract must have a positive duration (from and to dates cannot be equal)")
		}
	}
	return nil
}

// SanitizeHTML escapes HTML to prevent XSS
func SanitizeHTML(s string) string {
	return html.EscapeString(s)
}

// MaxWeeklyHours is the maximum number of hours in a week
const MaxWeeklyHours = 168.0

// ValidateWeeklyHours validates hours per week
func ValidateWeeklyHours(hours float64, fieldName string) error {
	if hours < 0 {
		return fmt.Errorf("%s cannot be negative", fieldName)
	}
	if hours > MaxWeeklyHours {
		return fmt.Errorf("%s cannot exceed %.0f hours per week", fieldName, MaxWeeklyHours)
	}
	return nil
}

// ValidateSalary validates salary in cents (must be non-negative)
func ValidateSalary(salary int) error {
	if salary < 0 {
		return fmt.Errorf("salary cannot be negative")
	}
	return nil
}

// CalculateAgeOnDate calculates the age in complete years on a given reference date.
// The age is the number of complete years from birthdate to referenceDate.
func CalculateAgeOnDate(birthdate, referenceDate time.Time) int {
	years := referenceDate.Year() - birthdate.Year()
	// Check if birthday hasn't occurred yet this year
	if referenceDate.Month() < birthdate.Month() ||
		(referenceDate.Month() == birthdate.Month() && referenceDate.Day() < birthdate.Day()) {
		years--
	}
	if years < 0 {
		return 0
	}
	return years
}

// ValidateEmployeeContractProperty validates a property value against its schema.
// Returns nil for unknown property names (allows extensibility).
func ValidateEmployeeContractProperty(name, value string) error {
	schema, known := models.PropertySchemas[name]
	if !known {
		// Allow unknown properties for extensibility
		return nil
	}

	switch schema.Type {
	case models.PropertyTypeFloat:
		f, err := strconv.ParseFloat(value, 64)
		if err != nil {
			return fmt.Errorf("property %s must be a valid number", name)
		}
		if schema.MinFloat != nil && f < *schema.MinFloat {
			return fmt.Errorf("property %s must be at least %v", name, *schema.MinFloat)
		}
		if schema.MaxFloat != nil && f > *schema.MaxFloat {
			return fmt.Errorf("property %s must be at most %v", name, *schema.MaxFloat)
		}

	case models.PropertyTypeInt:
		i, err := strconv.Atoi(value)
		if err != nil {
			return fmt.Errorf("property %s must be a valid integer", name)
		}
		if schema.MinInt != nil && i < *schema.MinInt {
			return fmt.Errorf("property %s must be at least %d", name, *schema.MinInt)
		}

	case models.PropertyTypeBool:
		if value != "true" && value != "false" {
			return fmt.Errorf("property %s must be 'true' or 'false'", name)
		}

	case models.PropertyTypeString:
		if len(schema.AllowedValues) > 0 {
			valid := false
			for _, allowed := range schema.AllowedValues {
				if value == allowed {
					valid = true
					break
				}
			}
			if !valid {
				return fmt.Errorf("property %s must be one of: %v", name, schema.AllowedValues)
			}
		}
	}

	return nil
}
