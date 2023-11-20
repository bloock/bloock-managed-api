// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent/certification"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// CertificationDelete is the builder for deleting a Certification entity.
type CertificationDelete struct {
	config
	hooks    []Hook
	mutation *CertificationMutation
}

// Where appends a list predicates to the CertificationDelete builder.
func (cd *CertificationDelete) Where(ps ...predicate.Certification) *CertificationDelete {
	cd.mutation.Where(ps...)
	return cd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (cd *CertificationDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, cd.sqlExec, cd.mutation, cd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (cd *CertificationDelete) ExecX(ctx context.Context) int {
	n, err := cd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (cd *CertificationDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(certification.Table, sqlgraph.NewFieldSpec(certification.FieldID, field.TypeUUID))
	if ps := cd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, cd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	cd.mutation.done = true
	return affected, err
}

// CertificationDeleteOne is the builder for deleting a single Certification entity.
type CertificationDeleteOne struct {
	cd *CertificationDelete
}

// Where appends a list predicates to the CertificationDelete builder.
func (cdo *CertificationDeleteOne) Where(ps ...predicate.Certification) *CertificationDeleteOne {
	cdo.cd.mutation.Where(ps...)
	return cdo
}

// Exec executes the deletion query.
func (cdo *CertificationDeleteOne) Exec(ctx context.Context) error {
	n, err := cdo.cd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{certification.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (cdo *CertificationDeleteOne) ExecX(ctx context.Context) {
	if err := cdo.Exec(ctx); err != nil {
		panic(err)
	}
}
