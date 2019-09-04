package custom

import (
	"time"

	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/migration"
	"github.com/iov-one/weave/x"
)

const (
	packageName        = "custom"
	newStateCost int64 = 100
)

// RegisterQuery registers buckets for querying.
func RegisterQuery(qr weave.QueryRouter) {
	NewTimedStateBucket().Register("customTimedStates", qr)
	NewStateBucket().Register("customStates", qr)
}

// RegisterRoutes registers handlers for message processing.
func RegisterRoutes(r weave.Registry, auth x.Authenticator, scheduler weave.Scheduler) {
	r = migration.SchemaMigratingRegistry(packageName, r)

	r.Handle(&CreateTimedStateMsg{}, NewCreateTimedStateHandler(auth, scheduler))
	r.Handle(&CreateStateMsg{}, NewCreateStateHandler(auth))
}

// RegisterCronRoutes registers routes that are not exposed to
// routers
func RegisterCronRoutes(
	r weave.Registry,
	auth x.Authenticator,
) {
	r.Handle(&DeleteTimedStateMsg{}, newDeleteTimedStateHandler(auth))
}

// ------------------- CreateTimedState HANDLERS -------------------

// CreateTimedStateHandler will handle creating custom indexed state buckets
type CreateTimedStateHandler struct {
	auth      x.Authenticator
	b         *TimedStateBucket
	scheduler weave.Scheduler
}

var _ weave.Handler = CreateTimedStateHandler{}

// NewCreateTimedStateHandler creates a handler
func NewCreateTimedStateHandler(auth x.Authenticator, scheduler weave.Scheduler) weave.Handler {
	return CreateTimedStateHandler{
		auth:      auth,
		b:         NewTimedStateBucket(),
		scheduler: scheduler,
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CreateTimedStateHandler) validate(ctx weave.Context, db weave.KVStore, tx weave.Tx) (*CreateTimedStateMsg, error) {
	var msg CreateTimedStateMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, errors.Wrap(err, "load msg")
	}

	if msg.DeleteAt != 0 && weave.InThePast(ctx, msg.DeleteAt.Time()) {
		return nil, errors.AppendField(nil, "DeleteAt", errors.ErrInput)
	}

	return &msg, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateTimedStateHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newStateCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h CreateTimedStateHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	msg, err := h.validate(ctx, store, tx)

	if err != nil {
		return nil, err
	}

	timedState := &TimedState{
		Metadata:       &weave.Metadata{},
		InnerStateEnum: msg.InnerStateEnum,
		Str:            msg.Str,
		Byte:           msg.Byte,
		DeleteAt:       msg.DeleteAt,
	}

	key, err := h.b.Put(store, nil, timedState)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store indexed state")
	}

	return &weave.DeliverResult{Data: key}, nil
}

// DeleteTimedStateHandler will handle deleting timed state
type DeleteTimedStateHandler struct {
	auth x.Authenticator
	b    *TimedStateBucket
}

var _ weave.Handler = DeleteTimedStateHandler{}

// newDeleteTimedStateHandler creates a handler
func newDeleteTimedStateHandler(auth x.Authenticator) weave.Handler {
	return DeleteTimedStateHandler{
		auth: auth,
		b:    NewTimedStateBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h DeleteTimedStateHandler) validate(ctx weave.Context, db weave.KVStore, tx weave.Tx) (*DeleteTimedStateMsg, error) {
	var msg DeleteTimedStateMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, errors.Wrap(err, "load msg")
	}

	return &msg, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h DeleteTimedStateHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newStateCost}, nil
}

// Deliver delete state if all preconditions are met
func (h DeleteTimedStateHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	msg, err := h.validate(ctx, store, tx)

	if err != nil {
		return nil, err
	}
	err = h.b.Delete(store, msg.TimedStateID)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store indexed state")
	}

	return &weave.DeliverResult{}, nil
}

// ------------------- State HANDLERS -------------------

// CreateStateHandler will handle creating custom state buckets
type CreateStateHandler struct {
	auth x.Authenticator
	b    *StateBucket
}

var _ weave.Handler = CreateStateHandler{}

// NewCreateStateHandler creates a handler
func NewCreateStateHandler(auth x.Authenticator) weave.Handler {
	return CreateStateHandler{
		auth: auth,
		b:    NewStateBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CreateStateHandler) validate(ctx weave.Context, db weave.KVStore, tx weave.Tx) (*CreateStateMsg, error) {
	var msg CreateStateMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, errors.Wrap(err, "load msg")
	}

	return &msg, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateStateHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newStateCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h CreateStateHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	msg, err := h.validate(ctx, store, tx)

	if err != nil {
		return nil, err
	}

	now := weave.AsUnixTime(time.Now())

	state := &State{
		Metadata:   &weave.Metadata{},
		InnerState: msg.InnerState,
		Address:    msg.Address,
		CreatedAt:  now,
	}

	res, err := h.b.Put(store, nil, state)

	if err != nil {
		return nil, err
	}

	return &weave.DeliverResult{Data: res}, err
}
