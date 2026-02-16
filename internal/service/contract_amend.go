package service

import (
	"time"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
)

// amendMode determines how a contract update should be handled.
type amendMode int

const (
	// amendModeInPlace means the contract started today or later — update in place.
	amendModeInPlace amendMode = iota
	// amendModeAmend means the contract started before today — close old + create new.
	amendModeAmend
)

// determineAmendMode decides whether to update in place or amend.
// Returns an error if the contract has already ended (To date is in the past).
func determineAmendMode(contractFrom time.Time, contractTo *time.Time) (amendMode, error) {
	today := models.TruncateToDate(time.Now())
	from := models.TruncateToDate(contractFrom)

	// Contract already ended → reject
	if contractTo != nil {
		to := models.TruncateToDate(*contractTo)
		if to.Before(today) {
			return 0, apperror.BadRequest("cannot update a contract that has already ended")
		}
	}

	// Contract starts today or in the future → update in place
	if !from.Before(today) {
		return amendModeInPlace, nil
	}

	// Contract started before today → amend (close old + create new)
	return amendModeAmend, nil
}
