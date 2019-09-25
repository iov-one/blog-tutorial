package blog

import (
	"context"
	"testing"
	"time"

	"github.com/iov-one/weave"
	"github.com/iov-one/weave/app"
	"github.com/iov-one/weave/errors"

	"github.com/iov-one/weave/store"
	"github.com/iov-one/weave/weavetest"
	"github.com/iov-one/weave/weavetest/assert"
)

func TestCreateUser(t *testing.T) {
	cases := map[string]struct {
		msg             weave.Msg
		expected        *User
		wantCheckErrs   map[string]*errors.Error
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &CreateUserMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Username: "Crpto0X",
				Bio:      "Best hacker in the universe",
			},
			expected: &User{
				Metadata: &weave.Metadata{Schema: 1},
				ID:       weavetest.SequenceID(1),
				Username: "Crpto0X",
				Bio:      "Best hacker in the universe",
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"Username": nil,
				"Bio":      nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"Username": nil,
				"Bio":      nil,
			},
		},
		"success missing bio": {
			msg: &CreateUserMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Username: "Crpto0X",
			},
			expected: &User{
				Metadata: &weave.Metadata{Schema: 1},
				ID:       weavetest.SequenceID(1),
				Username: "Crpto0X",
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"Username": nil,
				"Bio":      nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"Username": nil,
				"Bio":      nil,
			},
		},
		// TODO add missing metadata test
		"failure missing username": {
			msg: &CreateUserMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Bio:      "Best hacker in the universe",
			},
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"Username": errors.ErrModel,
				"Bio":      nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"Username": errors.ErrModel,
				"Bio":      nil,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			auth := &weavetest.Auth{}

			rt := app.NewRouter()
			RegisterRoutes(rt, auth)
			kv := store.MemStore()
			bucket := NewUserBucket()

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
				var stored User
				err := bucket.One(kv, res.Data, &stored)
				assert.Nil(t, err)

				// ensure registeredAt is after test creation time
				registeredAt := stored.RegisteredAt
				weave.InTheFuture(ctx, registeredAt.Time())

				// avoid registered at missing error
				tc.expected.RegisteredAt = registeredAt

				assert.Nil(t, err)
				assert.Equal(t, tc.expected, &stored)
			}
		})
	}
}
