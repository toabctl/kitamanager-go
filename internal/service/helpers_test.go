package service

import (
	"errors"
	"fmt"
	"testing"
	"time"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
)

// ptr returns a pointer to v.
func ptr[T any](v T) *T { return &v }

// --- test types for verifyOrgOwnership ---

type testOrgEntity struct{ orgID uint }

func (e *testOrgEntity) GetOrganizationID() uint { return e.orgID }

// --- test type for verifyRecordOwnership (implements models.PeriodRecord) ---

type testPeriodRecord struct {
	ownerID uint
	from    time.Time
	to      *time.Time
}

func (r *testPeriodRecord) GetOwnerID() uint   { return r.ownerID }
func (r *testPeriodRecord) GetFrom() time.Time { return r.from }
func (r *testPeriodRecord) GetTo() *time.Time  { return r.to }

// --- test type for validateNoOverlap ---

type testPeriodEntry struct {
	id     uint
	period models.Period
}

func testGetID(e testPeriodEntry) uint              { return e.id }
func testGetPeriod(e testPeriodEntry) models.Period { return e.period }

// ==========================================================================
// classifyStoreError
// ==========================================================================

func TestClassifyStoreError_NotFound(t *testing.T) {
	err := classifyStoreError(store.ErrNotFound, "child")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
	if err.Error() != "child not found" {
		t.Errorf("message = %q, want %q", err.Error(), "child not found")
	}
}

func TestClassifyStoreError_OtherError(t *testing.T) {
	original := errors.New("connection refused")
	err := classifyStoreError(original, "employee")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	var appErr *apperror.AppError
	if !errors.As(err, &appErr) {
		t.Fatalf("expected *apperror.AppError, got %T", err)
	}
	if appErr.Code != 500 {
		t.Errorf("Code = %d, want 500", appErr.Code)
	}
	if appErr.Message != "failed to fetch employee" {
		t.Errorf("Message = %q, want %q", appErr.Message, "failed to fetch employee")
	}
}

func TestClassifyStoreError_WrappedNotFound(t *testing.T) {
	wrapped := fmt.Errorf("wrapped: %w", store.ErrNotFound)
	err := classifyStoreError(wrapped, "section")
	// errors.Is should unwrap and detect store.ErrNotFound
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound for wrapped store.ErrNotFound, got %v", err)
	}
}

func TestClassifyStoreError_NotFoundMessage(t *testing.T) {
	err := classifyStoreError(store.ErrNotFound, "organization")
	if err.Error() != "organization not found" {
		t.Errorf("message = %q, want %q", err.Error(), "organization not found")
	}
}

// ==========================================================================
// verifyOrgOwnership
// ==========================================================================

func TestVerifyOrgOwnership_CorrectOrg(t *testing.T) {
	entity := &testOrgEntity{orgID: 5}
	err := verifyOrgOwnership(entity, 5, "child")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestVerifyOrgOwnership_WrongOrg(t *testing.T) {
	entity := &testOrgEntity{orgID: 5}
	err := verifyOrgOwnership(entity, 99, "child")
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestVerifyOrgOwnership_WrongOrg_Message(t *testing.T) {
	entity := &testOrgEntity{orgID: 5}
	err := verifyOrgOwnership(entity, 99, "employee")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "employee not found" {
		t.Errorf("message = %q, want %q", err.Error(), "employee not found")
	}
}

func TestVerifyOrgOwnership_NilInterface(t *testing.T) {
	// Passing an untyped nil (OrgOwned(nil))
	err := verifyOrgOwnership(nil, 1, "employee")
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound for nil entity, got %v", err)
	}
}

func TestVerifyOrgOwnership_ZeroOrgID(t *testing.T) {
	// OrgID 0 matching OrgID 0 should succeed
	entity := &testOrgEntity{orgID: 0}
	err := verifyOrgOwnership(entity, 0, "child")
	if err != nil {
		t.Errorf("expected nil for matching zero org IDs, got %v", err)
	}
}

// ==========================================================================
// verifyRecordOwnership
// ==========================================================================

func TestVerifyRecordOwnership_CorrectOwner(t *testing.T) {
	rec := &testPeriodRecord{ownerID: 10}
	err := verifyRecordOwnership(rec, 10, "contract")
	if err != nil {
		t.Errorf("expected nil, got %v", err)
	}
}

func TestVerifyRecordOwnership_WrongOwner(t *testing.T) {
	rec := &testPeriodRecord{ownerID: 10}
	err := verifyRecordOwnership(rec, 99, "contract")
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound, got %v", err)
	}
}

func TestVerifyRecordOwnership_NilRecord(t *testing.T) {
	err := verifyRecordOwnership(nil, 1, "contract")
	if !errors.Is(err, apperror.ErrNotFound) {
		t.Errorf("expected ErrNotFound for nil record, got %v", err)
	}
}

func TestVerifyRecordOwnership_WrongOwner_Message(t *testing.T) {
	rec := &testPeriodRecord{ownerID: 10}
	err := verifyRecordOwnership(rec, 99, "funding_period")
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	if err.Error() != "funding_period not found" {
		t.Errorf("message = %q, want %q", err.Error(), "funding_period not found")
	}
}

// ==========================================================================
// periodsOverlap
//
// Per models.Period doc: "both From and To are INCLUSIVE".
// Two inclusive ranges [A,B] and [C,D] overlap when A <= D AND C <= B.
// With nil To meaning "extends indefinitely into the future".
// ==========================================================================

func TestPeriodsOverlap_BothBounded_Overlapping(t *testing.T) {
	// [Jan 1 .. Mar 31] and [Mar 1 .. Jun 30] overlap
	from1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to1 := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
	from2 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	to2 := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	if !periodsOverlap(from1, &to1, from2, &to2) {
		t.Error("expected overlap")
	}
}

func TestPeriodsOverlap_BothBounded_NotOverlapping(t *testing.T) {
	// [Jan 1 .. Jan 31] and [Mar 1 .. Mar 31] do not overlap
	from1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to1 := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	from2 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	to2 := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)

	if periodsOverlap(from1, &to1, from2, &to2) {
		t.Error("expected no overlap")
	}
}

func TestPeriodsOverlap_Adjacent_NoOverlap(t *testing.T) {
	// [Jan 1 .. Jan 31] and [Feb 1 .. Feb 28] are adjacent (not overlapping)
	from1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to1 := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	from2 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)
	to2 := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)

	if periodsOverlap(from1, &to1, from2, &to2) {
		t.Error("expected no overlap for adjacent periods")
	}
}

func TestPeriodsOverlap_SameDayBoundary(t *testing.T) {
	// [Jan 1 .. Jan 31] and [Jan 31 .. Feb 28]
	// Since both From and To are INCLUSIVE, Jan 31 is in both ranges => overlap
	from1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to1 := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	from2 := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)
	to2 := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)

	if !periodsOverlap(from1, &to1, from2, &to2) {
		t.Error("expected overlap when to1 == from2 (inclusive boundaries)")
	}
}

func TestPeriodsOverlap_FirstOpenEnded(t *testing.T) {
	// [Jan 1 .. nil] and [Jun 1 .. Jun 30] - open-ended first overlaps everything after its start
	from1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	from2 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	to2 := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	if !periodsOverlap(from1, nil, from2, &to2) {
		t.Error("expected overlap with open-ended first period")
	}
}

func TestPeriodsOverlap_SecondOpenEnded(t *testing.T) {
	// [Jan 1 .. Mar 31] and [Feb 1 .. nil] - open-ended second starts within first
	from1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to1 := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
	from2 := time.Date(2024, 2, 1, 0, 0, 0, 0, time.UTC)

	if !periodsOverlap(from1, &to1, from2, nil) {
		t.Error("expected overlap with open-ended second period")
	}
}

func TestPeriodsOverlap_BothOpenEnded(t *testing.T) {
	// Two open-ended periods always overlap (both extend infinitely)
	from1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	from2 := time.Date(2025, 6, 1, 0, 0, 0, 0, time.UTC)

	if !periodsOverlap(from1, nil, from2, nil) {
		t.Error("expected overlap when both periods are open-ended")
	}
}

func TestPeriodsOverlap_FirstOpenEnded_SecondEndsBeforeFirstStarts(t *testing.T) {
	// [Jun 1 .. nil] and [Jan 1 .. Jan 31] - second ends before first starts => no overlap
	from1 := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	from2 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to2 := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)

	if periodsOverlap(from1, nil, from2, &to2) {
		t.Error("expected no overlap when open-ended period starts after bounded period ends")
	}
}

func TestPeriodsOverlap_IdenticalPeriods(t *testing.T) {
	// Identical ranges must overlap
	from := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	if !periodsOverlap(from, &to, from, &to) {
		t.Error("expected overlap for identical periods")
	}
}

func TestPeriodsOverlap_SingleDayPeriods_SameDay(t *testing.T) {
	// Single-day periods on the same day overlap (inclusive)
	day := time.Date(2024, 5, 15, 0, 0, 0, 0, time.UTC)

	if !periodsOverlap(day, &day, day, &day) {
		t.Error("expected overlap for same single-day periods")
	}
}

func TestPeriodsOverlap_SameDayDifferentTimes(t *testing.T) {
	// Same calendar date but different times of day.
	// Per the documented contract ("both From and To are INCLUSIVE" date ranges),
	// [Jan 1 .. Jan 31] and [Jan 31 .. Feb 28] should overlap regardless of
	// the time-of-day component, because Jan 31 is in both ranges.
	from1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to1 := time.Date(2024, 1, 31, 0, 0, 0, 0, time.UTC)    // midnight
	from2 := time.Date(2024, 1, 31, 23, 0, 0, 0, time.UTC) // 11pm same day
	to2 := time.Date(2024, 2, 28, 0, 0, 0, 0, time.UTC)

	if !periodsOverlap(from1, &to1, from2, &to2) {
		t.Error("expected overlap: to1 and from2 are the same calendar date (Jan 31)")
	}
}

func TestPeriodsOverlap_Symmetry(t *testing.T) {
	// periodsOverlap(a, b) should equal periodsOverlap(b, a)
	from1 := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to1 := time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)
	from2 := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	to2 := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)

	r1 := periodsOverlap(from1, &to1, from2, &to2)
	r2 := periodsOverlap(from2, &to2, from1, &to1)
	if r1 != r2 {
		t.Errorf("overlap not symmetric: (%v, %v) vs (%v, %v)", from1, to1, from2, to2)
	}
}

// ==========================================================================
// validateNoOverlap
// ==========================================================================

func TestValidateNoOverlap_NoConflict(t *testing.T) {
	existing := []testPeriodEntry{
		{id: 1, period: models.Period{
			From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			To:   ptr(time.Date(2024, 3, 31, 0, 0, 0, 0, time.UTC)),
		}},
	}
	from := time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 8, 31, 0, 0, 0, 0, time.UTC)

	err := validateNoOverlap(existing, testGetID, testGetPeriod, from, &to, nil)
	if err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestValidateNoOverlap_ConflictDetected(t *testing.T) {
	existing := []testPeriodEntry{
		{id: 1, period: models.Period{
			From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			To:   ptr(time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)),
		}},
	}
	from := time.Date(2024, 3, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 8, 31, 0, 0, 0, 0, time.UTC)

	err := validateNoOverlap(existing, testGetID, testGetPeriod, from, &to, nil)
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
}

func TestValidateNoOverlap_ExcludesSelf(t *testing.T) {
	existing := []testPeriodEntry{
		{id: 1, period: models.Period{
			From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			To:   ptr(time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)),
		}},
	}
	// Same range as existing[0], but excluded by ID
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	excludeID := uint(1)

	err := validateNoOverlap(existing, testGetID, testGetPeriod, from, &to, &excludeID)
	if err != nil {
		t.Errorf("expected nil when excluding self, got %v", err)
	}
}

func TestValidateNoOverlap_ExcludesOnlyMatchingID(t *testing.T) {
	// Two existing periods; exclude one but still overlap with the other
	existing := []testPeriodEntry{
		{id: 1, period: models.Period{
			From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			To:   ptr(time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)),
		}},
		{id: 2, period: models.Period{
			From: time.Date(2024, 4, 1, 0, 0, 0, 0, time.UTC),
			To:   ptr(time.Date(2024, 9, 30, 0, 0, 0, 0, time.UTC)),
		}},
	}
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 6, 30, 0, 0, 0, 0, time.UTC)
	excludeID := uint(1) // excludes id=1, but id=2 still overlaps

	err := validateNoOverlap(existing, testGetID, testGetPeriod, from, &to, &excludeID)
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest from non-excluded overlap, got %v", err)
	}
}

func TestValidateNoOverlap_EmptyExisting(t *testing.T) {
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)

	err := validateNoOverlap([]testPeriodEntry{}, testGetID, testGetPeriod, from, &to, nil)
	if err != nil {
		t.Errorf("expected nil for empty existing, got %v", err)
	}
}

func TestValidateNoOverlap_OpenEndedExisting(t *testing.T) {
	existing := []testPeriodEntry{
		{id: 1, period: models.Period{
			From: time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC),
			To:   nil, // open-ended
		}},
	}
	from := time.Date(2025, 1, 1, 0, 0, 0, 0, time.UTC)
	to := time.Date(2025, 6, 30, 0, 0, 0, 0, time.UTC)

	err := validateNoOverlap(existing, testGetID, testGetPeriod, from, &to, nil)
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest for overlap with open-ended existing, got %v", err)
	}
}

func TestValidateNoOverlap_NewOpenEndedConflicts(t *testing.T) {
	existing := []testPeriodEntry{
		{id: 1, period: models.Period{
			From: time.Date(2024, 6, 1, 0, 0, 0, 0, time.UTC),
			To:   ptr(time.Date(2024, 12, 31, 0, 0, 0, 0, time.UTC)),
		}},
	}
	from := time.Date(2024, 1, 1, 0, 0, 0, 0, time.UTC)
	// new period is open-ended
	err := validateNoOverlap(existing, testGetID, testGetPeriod, from, nil, nil)
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest for open-ended new period overlapping existing, got %v", err)
	}
}

// ==========================================================================
// validateRequiredName
// ==========================================================================

func TestValidateRequiredName_ValidName(t *testing.T) {
	result, err := validateRequiredName("  Sunshine Kita  ")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != "Sunshine Kita" {
		t.Errorf("result = %q, want %q", result, "Sunshine Kita")
	}
}

func TestValidateRequiredName_WhitespaceOnly(t *testing.T) {
	_, err := validateRequiredName("   ")
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest for whitespace-only, got %v", err)
	}
}

func TestValidateRequiredName_Empty(t *testing.T) {
	_, err := validateRequiredName("")
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest for empty string, got %v", err)
	}
}

func TestValidateRequiredName_Trimming(t *testing.T) {
	result, err := validateRequiredName("\t Hello \n")
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if result != "Hello" {
		t.Errorf("result = %q, want %q", result, "Hello")
	}
}

// ==========================================================================
// toResponseList
// ==========================================================================

func TestToResponseList_Converts(t *testing.T) {
	items := []int{1, 2, 3}
	result := toResponseList(items, func(i *int) string {
		return fmt.Sprintf("item-%d", *i)
	})
	if len(result) != 3 {
		t.Fatalf("len = %d, want 3", len(result))
	}
	expected := []string{"item-1", "item-2", "item-3"}
	for i, want := range expected {
		if result[i] != want {
			t.Errorf("result[%d] = %q, want %q", i, result[i], want)
		}
	}
}

func TestToResponseList_EmptySlice(t *testing.T) {
	result := toResponseList([]int{}, func(i *int) string {
		return fmt.Sprintf("%d", *i)
	})
	if len(result) != 0 {
		t.Errorf("len = %d, want 0", len(result))
	}
}

func TestToResponseList_NilSlice(t *testing.T) {
	var items []int
	result := toResponseList(items, func(i *int) string {
		return fmt.Sprintf("%d", *i)
	})
	// make([]R, 0) returns empty (not nil) slice
	if len(result) != 0 {
		t.Errorf("len = %d, want 0", len(result))
	}
}

// ==========================================================================
// applyPersonUpdates
// ==========================================================================

func TestApplyPersonUpdates_AllFieldsNil(t *testing.T) {
	person := &models.Person{
		FirstName: "Original",
		LastName:  "Name",
		Gender:    "male",
		Birthdate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	err := applyPersonUpdates(person, personUpdateFields{})
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if person.FirstName != "Original" {
		t.Errorf("FirstName changed to %q", person.FirstName)
	}
	if person.LastName != "Name" {
		t.Errorf("LastName changed to %q", person.LastName)
	}
	if person.Gender != "male" {
		t.Errorf("Gender changed to %q", person.Gender)
	}
	expected := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	if !person.Birthdate.Equal(expected) {
		t.Errorf("Birthdate changed to %v", person.Birthdate)
	}
}

func TestApplyPersonUpdates_AllFieldsSet(t *testing.T) {
	person := &models.Person{
		FirstName: "Old",
		LastName:  "Name",
		Gender:    "male",
		Birthdate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	fields := personUpdateFields{
		FirstName: ptr("New"),
		LastName:  ptr("Person"),
		Gender:    ptr("female"),
		Birthdate: ptr("1995-06-15"),
	}
	err := applyPersonUpdates(person, fields)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if person.FirstName != "New" {
		t.Errorf("FirstName = %q, want %q", person.FirstName, "New")
	}
	if person.LastName != "Person" {
		t.Errorf("LastName = %q, want %q", person.LastName, "Person")
	}
	if person.Gender != "female" {
		t.Errorf("Gender = %q, want %q", person.Gender, "female")
	}
	expected := time.Date(1995, 6, 15, 0, 0, 0, 0, time.UTC)
	if !person.Birthdate.Equal(expected) {
		t.Errorf("Birthdate = %v, want %v", person.Birthdate, expected)
	}
}

func TestApplyPersonUpdates_InvalidGender(t *testing.T) {
	person := &models.Person{Gender: "male"}
	fields := personUpdateFields{
		Gender: ptr("invalid"),
	}
	err := applyPersonUpdates(person, fields)
	if err == nil {
		t.Fatal("expected error for invalid gender")
	}
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
	// Gender should not have changed on validation failure
	if person.Gender != "male" {
		t.Errorf("Gender changed to %q despite error", person.Gender)
	}
}

func TestApplyPersonUpdates_InvalidBirthdateFormat(t *testing.T) {
	original := time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC)
	person := &models.Person{Birthdate: original}
	fields := personUpdateFields{
		Birthdate: ptr("not-a-date"),
	}
	err := applyPersonUpdates(person, fields)
	if err == nil {
		t.Fatal("expected error for invalid birthdate")
	}
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
	// Birthdate should not have changed on validation failure
	if !person.Birthdate.Equal(original) {
		t.Errorf("Birthdate changed to %v despite error", person.Birthdate)
	}
}

func TestApplyPersonUpdates_WhitespaceFirstName(t *testing.T) {
	person := &models.Person{FirstName: "Keep"}
	fields := personUpdateFields{
		FirstName: ptr("   "),
	}
	err := applyPersonUpdates(person, fields)
	if err == nil {
		t.Fatal("expected error for whitespace-only first name")
	}
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
	// FirstName should not have changed on validation failure
	if person.FirstName != "Keep" {
		t.Errorf("FirstName changed to %q despite error", person.FirstName)
	}
}

func TestApplyPersonUpdates_WhitespaceLastName(t *testing.T) {
	person := &models.Person{LastName: "Keep"}
	fields := personUpdateFields{
		LastName: ptr(""),
	}
	err := applyPersonUpdates(person, fields)
	if err == nil {
		t.Fatal("expected error for empty last name")
	}
	if !errors.Is(err, apperror.ErrBadRequest) {
		t.Errorf("expected ErrBadRequest, got %v", err)
	}
}

func TestApplyPersonUpdates_TrimsNames(t *testing.T) {
	person := &models.Person{FirstName: "Old", LastName: "Name"}
	fields := personUpdateFields{
		FirstName: ptr("  Alice  "),
		LastName:  ptr("  Smith  "),
	}
	err := applyPersonUpdates(person, fields)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if person.FirstName != "Alice" {
		t.Errorf("FirstName = %q, want %q", person.FirstName, "Alice")
	}
	if person.LastName != "Smith" {
		t.Errorf("LastName = %q, want %q", person.LastName, "Smith")
	}
}

func TestApplyPersonUpdates_PartialUpdate_OnlyFirstName(t *testing.T) {
	person := &models.Person{
		FirstName: "Old",
		LastName:  "Keep",
		Gender:    "male",
		Birthdate: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	fields := personUpdateFields{
		FirstName: ptr("New"),
	}
	err := applyPersonUpdates(person, fields)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}
	if person.FirstName != "New" {
		t.Errorf("FirstName = %q, want %q", person.FirstName, "New")
	}
	if person.LastName != "Keep" {
		t.Errorf("LastName changed to %q, should have stayed %q", person.LastName, "Keep")
	}
	if person.Gender != "male" {
		t.Errorf("Gender changed to %q, should have stayed %q", person.Gender, "male")
	}
}

func TestApplyPersonUpdates_DiverseGender(t *testing.T) {
	person := &models.Person{Gender: "male"}
	fields := personUpdateFields{
		Gender: ptr("diverse"),
	}
	err := applyPersonUpdates(person, fields)
	if err != nil {
		t.Fatalf("expected no error for diverse gender, got %v", err)
	}
	if person.Gender != "diverse" {
		t.Errorf("Gender = %q, want %q", person.Gender, "diverse")
	}
}
