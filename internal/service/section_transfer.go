package service

import (
	"context"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

type transferAction int

const (
	transferNone    transferAction = iota // no active contract
	transferUpdate                        // same-day: update section on existing contract
	transferReplace                       // close existing, create new
)

// decideSectionTransfer determines how to handle a section change on a contract.
// If the contract started today, we update it in place (same-day correction).
// Otherwise we close the old contract and create a new one.
func decideSectionTransfer(contractFrom time.Time) transferAction {
	today := time.Now().UTC().Truncate(24 * time.Hour)
	if contractFrom.Truncate(24 * time.Hour).Equal(today) {
		return transferUpdate
	}
	return transferReplace
}

// sameSectionID returns true if both section IDs are equal (including both nil).
func sameSectionID(a, b *uint) bool {
	if a == nil && b == nil {
		return true
	}
	if a == nil || b == nil {
		return false
	}
	return *a == *b
}

// handleSectionTransfer is a generic helper that creates a contract transition when a person
// (child or employee) moves between sections. It captures the shared logic of:
//  1. Fetching the current contract
//  2. Checking if the section actually changed
//  3. Either updating in place (same-day) or closing + creating a new contract (replace)
//
// Parameters:
//   - contracts: the generic ContractStorer for GetCurrentContract
//   - updateContract: store method to update a contract (type-specific, on the parent store)
//   - createContract: store method to create a contract (type-specific, on the parent store)
//   - getBase: extracts the embedded BaseContract from the concrete contract type
//   - personID: the child or employee ID
//   - newSectionID: the target section
//   - newContractFactory: builds a new contract of the concrete type for the transferReplace case
func handleSectionTransfer[T models.HasPeriod](
	ctx context.Context,
	contracts store.ContractStorer[T],
	updateContract func(context.Context, *T) error,
	createContract func(context.Context, *T) error,
	getBase func(*T) *models.BaseContract,
	personID uint,
	newSectionID *uint,
	newContractFactory func(old *T, today time.Time) *T,
) error {
	contract, err := contracts.GetCurrentContract(ctx, personID)
	if err != nil {
		return apperror.InternalWrap(err, "failed to fetch current contract")
	}
	if contract == nil {
		return nil // no active contract, nothing to transition
	}

	base := getBase(contract)

	// If the contract already has the same section, nothing to do
	if sameSectionID(base.SectionID, newSectionID) {
		return nil
	}

	today := time.Now().UTC().Truncate(24 * time.Hour)

	switch decideSectionTransfer(base.From) {
	case transferUpdate:
		// Same-day: just update the section on the existing contract
		base.SectionID = newSectionID
		base.Section = nil
		return updateContract(ctx, contract)
	case transferReplace:
		// Close old contract (end = yesterday) and create a new one starting today
		yesterday := today.AddDate(0, 0, -1)
		base.To = &yesterday
		base.Section = nil
		if err := updateContract(ctx, contract); err != nil {
			return apperror.InternalWrap(err, "failed to close existing contract")
		}

		newContract := newContractFactory(contract, today)
		return createContract(ctx, newContract)
	}
	return nil
}
