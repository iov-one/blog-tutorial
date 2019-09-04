package custom

import (
	"context"
	"testing"
	"time"

	"github.com/iov-one/weave"
	"github.com/iov-one/weave/app"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/migration"
	"github.com/iov-one/weave/store"
	"github.com/iov-one/weave/weavetest"
	"github.com/iov-one/weave/weavetest/assert"
)

func TestCreateTimedState(t *testing.T) {
	meta := &weave.Metadata{Schema: 1}
	now := weave.AsUnixTime(time.Now())
	future := now.Add(time.Hour)
	past := now.Add(-1 * time.Hour)

	cases := map[string]struct {
		msg             weave.Msg
		expected        *TimedState
		wantCheckErrs   map[string]*errors.Error
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &CreateTimedStateMsg{
				Metadata:       meta,
				InnerStateEnum: InnerStateEnum_CaseOne,
				Str:            "cstm_str",
				Byte:           []byte{0, 1},
				DeleteAt:       future,
			},
			expected: &TimedState{
				Metadata:       meta,
				InnerStateEnum: InnerStateEnum_CaseOne,
				Str:            "cstm_str",
				Byte:           []byte{0, 1},
				DeleteAt:       future,
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            nil,
				"Byte":           nil,
				"DeleteAt":       nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            nil,
				"Byte":           nil,
				"DeleteAt":       nil,
			},
		},
		"success no delete at": {
			msg: &CreateTimedStateMsg{
				Metadata:       meta,
				InnerStateEnum: InnerStateEnum_CaseOne,
				Str:            "cstm_str",
				Byte:           []byte{0, 1},
			},
			expected: &TimedState{
				Metadata:       meta,
				InnerStateEnum: InnerStateEnum_CaseOne,
				Str:            "cstm_str",
				Byte:           []byte{0, 1},
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            nil,
				"Byte":           nil,
				"DeleteAt":       nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            nil,
				"Byte":           nil,
				"DeleteAt":       nil,
			},
		},
		"failure delete at is in the past": {
			msg: &CreateTimedStateMsg{
				Metadata:       meta,
				InnerStateEnum: InnerStateEnum_CaseOne,
				Str:            "cstm_str",
				Byte:           []byte{0, 1},
				DeleteAt:       past,
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            nil,
				"Byte":           nil,
				"DeleteAt":       errors.ErrInput,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            nil,
				"Byte":           nil,
				"DeleteAt":       errors.ErrInput,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			auth := &weavetest.Auth{}

			rt := app.NewRouter()
			RegisterRoutes(rt, auth, &weavetest.Cron{})

			kv := store.MemStore()
			bucket := NewTimedStateBucket()
			migration.MustInitPkg(kv, packageName)

			tx := &weavetest.Tx{Msg: tc.msg}

			ctx := weave.WithBlockTime(context.Background(), time.Now().Round(time.Second))

			if _, err := rt.Check(ctx, kv, tx); err != nil {
				for field, wantErr := range tc.wantCheckErrs {
					assert.FieldError(t, err, field, wantErr)
				}
			}

			res, err := rt.Deliver(ctx, kv, tx)
			if err != nil {
				for field, wantErr := range tc.wantDeliverErrs {
					assert.FieldError(t, err, field, wantErr)
				}
			}

			if tc.expected != nil {
				var stored TimedState
				err := bucket.One(kv, res.Data, &stored)

				assert.Nil(t, err)
				assert.Equal(t, tc.expected, &stored)
			}
		})
	}
}

func TestDeleteTimedState(t *testing.T) {
	meta := &weave.Metadata{Schema: 1}
	timeStateID := weavetest.SequenceID(1)

	cases := map[string]struct {
		msg             weave.Msg
		shouldDeliver   bool
		wantCheckErrs   map[string]*errors.Error
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &DeleteTimedStateMsg{
				Metadata:     meta,
				TimedStateID: timeStateID,
			},
			shouldDeliver: true,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"TimedStateID": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"TimedStateID": nil,
			},
		},
		"failure invalid timed state id": {
			msg: &DeleteTimedStateMsg{
				Metadata:     meta,
				TimedStateID: []byte{0, 1},
			},
			shouldDeliver: false,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"TimedStateID": errors.ErrInput,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"TimedStateID": errors.ErrInput,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			auth := &weavetest.Auth{}

			rt := app.NewRouter()
			RegisterCronRoutes(rt, auth)
			kv := store.MemStore()
			bucket := NewTimedStateBucket()
			migration.MustInitPkg(kv, packageName)

			stored := &TimedState{
				Metadata:       &weave.Metadata{Schema: 1},
				InnerStateEnum: InnerStateEnum_CaseOne,
				Str:            "cstm_string",
				Byte:           []byte{0, 1},
			}

			_, err := bucket.Put(kv, timeStateID, stored)
			if err != nil {
				t.Fatalf("unexpected error: %+v", err)
			}

			tx := &weavetest.Tx{Msg: tc.msg}

			if _, err := rt.Check(context.TODO(), kv, tx); err != nil {
				for field, wantErr := range tc.wantCheckErrs {
					assert.FieldError(t, err, field, wantErr)
				}
			}

			_, err = rt.Deliver(context.TODO(), kv, tx)

			for field, wantErr := range tc.wantDeliverErrs {
				assert.FieldError(t, err, field, wantErr)
			}

			if tc.shouldDeliver {
				err = bucket.Has(kv, timeStateID)
				assert.IsErr(t, errors.ErrNotFound, err)
			}
		})
	}
}

func TestCreateState(t *testing.T) {
	meta := &weave.Metadata{Schema: 1}
	now := weave.AsUnixTime(time.Now())
	address := weavetest.NewCondition().Address()

	cases := map[string]struct {
		msg             weave.Msg
		expected        *State
		wantCheckErrs   map[string]*errors.Error
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &CreateStateMsg{
				Metadata:   meta,
				InnerState: &InnerState{St1: 1, St2: 2},
				Address:    address,
			},
			expected: &State{
				Metadata:   meta,
				InnerState: &InnerState{St1: 1, St2: 2},
				Address:    address,
				CreatedAt:  now,
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"InnerState": nil,
				"Address":    nil,
				"CreatedAt":  nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"InnerState": nil,
				"Address":    nil,
				"CreatedAt":  nil,
			},
		},
		"failure empty message": {
			msg: &CreateStateMsg{},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"InnerState": errors.ErrEmpty,
				"Address":    errors.ErrEmpty,
				"CreatedAt":  nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"InnerState": errors.ErrEmpty,
				"Address":    errors.ErrEmpty,
				"CreatedAt":  nil,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			auth := &weavetest.Auth{}

			h := NewCreateStateHandler(auth)
			kv := store.MemStore()
			bucket := NewStateBucket()
			migration.MustInitPkg(kv, packageName)

			tx := &weavetest.Tx{Msg: tc.msg}

			if _, err := h.Check(context.TODO(), kv, tx); err != nil {
				for field, wantErr := range tc.wantCheckErrs {
					assert.FieldError(t, err, field, wantErr)
				}
			}

			res, err := h.Deliver(context.TODO(), kv, tx)
			for field, wantErr := range tc.wantDeliverErrs {
				assert.FieldError(t, err, field, wantErr)
			}

			if tc.expected != nil {
				err := bucket.Has(kv, res.Data)
				assert.Nil(t, err)
			}
		})
	}
}
