// Code generated by ent, DO NOT EDIT.

package ent

import (
	"bloock-managed-api/internal/platform/repository/sql/ent/localkey"
	"context"
	"errors"
	"fmt"

	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/bloock/bloock-sdk-go/v2/entity/key"
	"github.com/google/uuid"
)

// LocalKeyCreate is the builder for creating a LocalKey entity.
type LocalKeyCreate struct {
	config
	mutation *LocalKeyMutation
	hooks    []Hook
}

// SetLocalKey sets the "local_key" field.
func (lkc *LocalKeyCreate) SetLocalKey(kk *key.LocalKey) *LocalKeyCreate {
	lkc.mutation.SetLocalKey(kk)
	return lkc
}

// SetKeyType sets the "key_type" field.
func (lkc *LocalKeyCreate) SetKeyType(s string) *LocalKeyCreate {
	lkc.mutation.SetKeyType(s)
	return lkc
}

// SetID sets the "id" field.
func (lkc *LocalKeyCreate) SetID(u uuid.UUID) *LocalKeyCreate {
	lkc.mutation.SetID(u)
	return lkc
}

// SetNillableID sets the "id" field if the given value is not nil.
func (lkc *LocalKeyCreate) SetNillableID(u *uuid.UUID) *LocalKeyCreate {
	if u != nil {
		lkc.SetID(*u)
	}
	return lkc
}

// Mutation returns the LocalKeyMutation object of the builder.
func (lkc *LocalKeyCreate) Mutation() *LocalKeyMutation {
	return lkc.mutation
}

// Save creates the LocalKey in the database.
func (lkc *LocalKeyCreate) Save(ctx context.Context) (*LocalKey, error) {
	lkc.defaults()
	return withHooks(ctx, lkc.sqlSave, lkc.mutation, lkc.hooks)
}

// SaveX calls Save and panics if Save returns an error.
func (lkc *LocalKeyCreate) SaveX(ctx context.Context) *LocalKey {
	v, err := lkc.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lkc *LocalKeyCreate) Exec(ctx context.Context) error {
	_, err := lkc.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lkc *LocalKeyCreate) ExecX(ctx context.Context) {
	if err := lkc.Exec(ctx); err != nil {
		panic(err)
	}
}

// defaults sets the default values of the builder before save.
func (lkc *LocalKeyCreate) defaults() {
	if _, ok := lkc.mutation.ID(); !ok {
		v := localkey.DefaultID()
		lkc.mutation.SetID(v)
	}
}

// check runs all checks and user-defined validators on the builder.
func (lkc *LocalKeyCreate) check() error {
	if _, ok := lkc.mutation.LocalKey(); !ok {
		return &ValidationError{Name: "local_key", err: errors.New(`ent: missing required field "LocalKey.local_key"`)}
	}
	if _, ok := lkc.mutation.KeyType(); !ok {
		return &ValidationError{Name: "key_type", err: errors.New(`ent: missing required field "LocalKey.key_type"`)}
	}
	return nil
}

func (lkc *LocalKeyCreate) sqlSave(ctx context.Context) (*LocalKey, error) {
	if err := lkc.check(); err != nil {
		return nil, err
	}
	_node, _spec := lkc.createSpec()
	if err := sqlgraph.CreateNode(ctx, lkc.driver, _spec); err != nil {
		if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	if _spec.ID.Value != nil {
		if id, ok := _spec.ID.Value.(*uuid.UUID); ok {
			_node.ID = *id
		} else if err := _node.ID.Scan(_spec.ID.Value); err != nil {
			return nil, err
		}
	}
	lkc.mutation.id = &_node.ID
	lkc.mutation.done = true
	return _node, nil
}

func (lkc *LocalKeyCreate) createSpec() (*LocalKey, *sqlgraph.CreateSpec) {
	var (
		_node = &LocalKey{config: lkc.config}
		_spec = sqlgraph.NewCreateSpec(localkey.Table, sqlgraph.NewFieldSpec(localkey.FieldID, field.TypeUUID))
	)
	if id, ok := lkc.mutation.ID(); ok {
		_node.ID = id
		_spec.ID.Value = &id
	}
	if value, ok := lkc.mutation.LocalKey(); ok {
		_spec.SetField(localkey.FieldLocalKey, field.TypeJSON, value)
		_node.LocalKey = value
	}
	if value, ok := lkc.mutation.KeyType(); ok {
		_spec.SetField(localkey.FieldKeyType, field.TypeString, value)
		_node.KeyType = value
	}
	return _node, _spec
}

// LocalKeyCreateBulk is the builder for creating many LocalKey entities in bulk.
type LocalKeyCreateBulk struct {
	config
	err      error
	builders []*LocalKeyCreate
}

// Save creates the LocalKey entities in the database.
func (lkcb *LocalKeyCreateBulk) Save(ctx context.Context) ([]*LocalKey, error) {
	if lkcb.err != nil {
		return nil, lkcb.err
	}
	specs := make([]*sqlgraph.CreateSpec, len(lkcb.builders))
	nodes := make([]*LocalKey, len(lkcb.builders))
	mutators := make([]Mutator, len(lkcb.builders))
	for i := range lkcb.builders {
		func(i int, root context.Context) {
			builder := lkcb.builders[i]
			builder.defaults()
			var mut Mutator = MutateFunc(func(ctx context.Context, m Mutation) (Value, error) {
				mutation, ok := m.(*LocalKeyMutation)
				if !ok {
					return nil, fmt.Errorf("unexpected mutation type %T", m)
				}
				if err := builder.check(); err != nil {
					return nil, err
				}
				builder.mutation = mutation
				var err error
				nodes[i], specs[i] = builder.createSpec()
				if i < len(mutators)-1 {
					_, err = mutators[i+1].Mutate(root, lkcb.builders[i+1].mutation)
				} else {
					spec := &sqlgraph.BatchCreateSpec{Nodes: specs}
					// Invoke the actual operation on the latest mutation in the chain.
					if err = sqlgraph.BatchCreate(ctx, lkcb.driver, spec); err != nil {
						if sqlgraph.IsConstraintError(err) {
							err = &ConstraintError{msg: err.Error(), wrap: err}
						}
					}
				}
				if err != nil {
					return nil, err
				}
				mutation.id = &nodes[i].ID
				mutation.done = true
				return nodes[i], nil
			})
			for i := len(builder.hooks) - 1; i >= 0; i-- {
				mut = builder.hooks[i](mut)
			}
			mutators[i] = mut
		}(i, ctx)
	}
	if len(mutators) > 0 {
		if _, err := mutators[0].Mutate(ctx, lkcb.builders[0].mutation); err != nil {
			return nil, err
		}
	}
	return nodes, nil
}

// SaveX is like Save, but panics if an error occurs.
func (lkcb *LocalKeyCreateBulk) SaveX(ctx context.Context) []*LocalKey {
	v, err := lkcb.Save(ctx)
	if err != nil {
		panic(err)
	}
	return v
}

// Exec executes the query.
func (lkcb *LocalKeyCreateBulk) Exec(ctx context.Context) error {
	_, err := lkcb.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (lkcb *LocalKeyCreateBulk) ExecX(ctx context.Context) {
	if err := lkcb.Exec(ctx); err != nil {
		panic(err)
	}
}
