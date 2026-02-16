package service

import (
	"context"

	"github.com/eenemeene/kitamanager-go/internal/apperror"
	"github.com/eenemeene/kitamanager-go/internal/models"
	"github.com/eenemeene/kitamanager-go/internal/store"
	"github.com/eenemeene/kitamanager-go/internal/validation"
)

// personList fetches all entities with pagination.
func personList[T any, R any](
	ctx context.Context,
	findAll func(ctx context.Context, limit, offset int) ([]T, int64, error),
	toResponse func(*T) R,
	resourceName string,
	limit, offset int,
) ([]R, int64, error) {
	items, total, err := findAll(ctx, limit, offset)
	if err != nil {
		return nil, 0, apperror.InternalWrap(err, "failed to fetch "+resourceName)
	}
	return toResponseList(items, toResponse), total, nil
}

// personGetByID fetches an entity by ID with org ownership check.
func personGetByID[T OrgOwned, R any](
	ctx context.Context,
	findByID func(ctx context.Context, id uint) (*T, error),
	toResponse func(*T) R,
	id, orgID uint,
	resourceName string,
) (*R, error) {
	entity, err := findByID(ctx, id)
	if err != nil {
		return nil, apperror.NotFound(resourceName)
	}
	if (*entity).GetOrganizationID() != orgID {
		return nil, apperror.NotFound(resourceName)
	}
	resp := toResponse(entity)
	return &resp, nil
}

// personCreate validates person fields and creates entity.
func personCreate[T any, R any](
	ctx context.Context,
	fields *validation.PersonCreateFields,
	buildEntity func(person models.Person) *T,
	createFn func(ctx context.Context, entity *T) error,
	toResponse func(*T) R,
	orgID uint,
	resourceName string,
) (*R, error) {
	person, err := validation.ValidatePersonCreate(fields)
	if err != nil {
		return nil, err
	}

	entity := buildEntity(models.Person{
		OrganizationID: orgID,
		FirstName:      person.FirstName,
		LastName:       person.LastName,
		Gender:         person.Gender,
		Birthdate:      person.Birthdate,
	})

	if err := createFn(ctx, entity); err != nil {
		return nil, apperror.InternalWrap(err, "failed to create "+resourceName)
	}

	resp := toResponse(entity)
	return &resp, nil
}

// personUpdate validates and applies person field updates.
func personUpdate[T OrgOwned, R any](
	ctx context.Context,
	findByID func(ctx context.Context, id uint) (*T, error),
	getPerson func(*T) *models.Person,
	updateFn func(ctx context.Context, entity *T) error,
	toResponse func(*T) R,
	id, orgID uint,
	fields personUpdateFields,
	resourceName string,
) (*R, error) {
	entity, err := findByID(ctx, id)
	if err != nil {
		return nil, apperror.NotFound(resourceName)
	}
	if (*entity).GetOrganizationID() != orgID {
		return nil, apperror.NotFound(resourceName)
	}

	if err := applyPersonUpdates(getPerson(entity), fields); err != nil {
		return nil, err
	}

	if err := updateFn(ctx, entity); err != nil {
		return nil, apperror.InternalWrap(err, "failed to update "+resourceName)
	}

	// Reload to get fresh associations
	entity, err = findByID(ctx, id)
	if err != nil {
		return nil, apperror.InternalWrap(err, "failed to reload "+resourceName+" after update")
	}

	resp := toResponse(entity)
	return &resp, nil
}

// personDelete validates ownership and deletes within a transaction.
func personDelete[T OrgOwned](
	ctx context.Context,
	transactor store.Transactor,
	findByID func(ctx context.Context, id uint) (*T, error),
	deleteFn func(ctx context.Context, id uint) error,
	id, orgID uint,
	resourceName string,
) error {
	return transactor.InTransaction(ctx, func(txCtx context.Context) error {
		entity, err := findByID(txCtx, id)
		if err != nil {
			return apperror.NotFound(resourceName)
		}
		if (*entity).GetOrganizationID() != orgID {
			return apperror.NotFound(resourceName)
		}
		if err := deleteFn(txCtx, id); err != nil {
			return apperror.InternalWrap(err, "failed to delete "+resourceName)
		}
		return nil
	})
}
