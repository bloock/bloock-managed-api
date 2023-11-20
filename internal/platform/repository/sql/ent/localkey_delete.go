// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent/localkey"
	"github.com/bloock/bloock-managed-api/internal/platform/repository/sql/ent/predicate"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
)

// LocalKeyDelete is the builder for deleting a LocalKey entity.
type LocalKeyDelete struct {
	config
	hooks    []Hook
	mutation *LocalKeyMutation
}

// Where appends a list predicates to the LocalKeyDelete builder.
func (lkd *LocalKeyDelete) Where(ps ...predicate.LocalKey) *LocalKeyDelete {
	lkd.mutation.Where(ps...)
	return lkd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (lkd *LocalKeyDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, lkd.sqlExec, lkd.mutation, lkd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (lkd *LocalKeyDelete) ExecX(ctx context.Context) int {
	n, err := lkd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (lkd *LocalKeyDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(localkey.Table, sqlgraph.NewFieldSpec(localkey.FieldID, field.TypeUUID))
	if ps := lkd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, lkd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	lkd.mutation.done = true
	return affected, err
}

// LocalKeyDeleteOne is the builder for deleting a single LocalKey entity.
type LocalKeyDeleteOne struct {
	lkd *LocalKeyDelete
}

// Where appends a list predicates to the LocalKeyDelete builder.
func (lkdo *LocalKeyDeleteOne) Where(ps ...predicate.LocalKey) *LocalKeyDeleteOne {
	lkdo.lkd.mutation.Where(ps...)
	return lkdo
}

// Exec executes the deletion query.
func (lkdo *LocalKeyDeleteOne) Exec(ctx context.Context) error {
	n, err := lkdo.lkd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{localkey.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (lkdo *LocalKeyDeleteOne) ExecX(ctx context.Context) {
	if err := lkdo.Exec(ctx); err != nil {
		panic(err)
	}
}
