package blog

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/x"
)

const (
	packageName       = "blog"
	newUserCost int64 = 100
)

// RegisterQuery registers buckets for querying.
func RegisterQuery(qr weave.QueryRouter) {
}

// RegisterRoutes registers handlers for message processing.
func RegisterRoutes(r weave.Registry, auth x.Authenticator) {
	//r = migration.SchemaMigratingRegistry(packageName, r)
	r.Handle(&CreateUserMsg{}, NewCreateUserHandler(auth))
}

// ------------------- CreateUserHandler -------------------

// CreateUserHandler will handle CreateUserMsg
type CreateUserHandler struct {
	auth x.Authenticator
	b    *UserBucket
}

var _ weave.Handler = CreateUserHandler{}

// NewCreateUserHandler creates a create user message handler
func NewCreateUserHandler(auth x.Authenticator) weave.Handler {
	return CreateUserHandler{
		auth: auth,
		b:    NewUserBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CreateUserHandler) validate(ctx weave.Context, db weave.KVStore, tx weave.Tx) (*CreateUserMsg, error) {
	var msg CreateUserMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, errors.Wrap(err, "load msg")
	}

	return &msg, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateUserHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newUserCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h CreateUserHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	msg, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	blockTime, err := weave.BlockTime(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "no block time in header")
	}
	now := weave.AsUnixTime(blockTime)

	user := &User{
		Metadata:     msg.Metadata,
		Username:     msg.Username,
		Bio:          msg.Bio,
		RegisteredAt: now,
	}

	err = h.b.Put(store, user)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store user")
	}

	// Returns generated user ID as response
	return &weave.DeliverResult{Data: user.ID}, nil
}
