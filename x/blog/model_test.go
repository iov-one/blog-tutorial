package blog

import (
	"testing"
	"time"

	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/orm"
	"github.com/iov-one/weave/weavetest"
	"github.com/iov-one/weave/weavetest/assert"
)

func TestValidateUser(t *testing.T) {
	now := weave.AsUnixTime(time.Now())

	cases := map[string]struct {
		model    orm.Model
		wantErrs map[string]*errors.Error
	}{
		"success": {
			model: &User{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				Username:     "Crypt0xxx",
				Bio:          "Best hacker in the universe",
				RegisteredAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"Username":     nil,
				"Bio":          nil,
				"RegisteredAt": nil,
			},
		},
		"failure, no ID": {
			model: &User{
				Metadata:     &weave.Metadata{Schema: 1},
				Username:     "Crypt0xxx",
				Bio:          "Best hacker in the universe",
				RegisteredAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           errors.ErrEmpty,
				"Username":     nil,
				"Bio":          nil,
				"RegisteredAt": nil,
			},
		},
		"failure, missing username": {
			model: &User{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				Bio:          "Best hacker in the universe",
				RegisteredAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"Username":     errors.ErrModel,
				"Bio":          nil,
				"RegisteredAt": nil,
			},
		},
		"failure, missing registered at": {
			model: &User{
				Metadata: &weave.Metadata{Schema: 1},
				ID:       weavetest.SequenceID(1),
				Username: "Crypt0xxx",
				Bio:      "Best hacker in the universe",
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"Username":     nil,
				"Bio":          nil,
				"RegisteredAt": errors.ErrEmpty,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			err := tc.model.Validate()
			for field, wantErr := range tc.wantErrs {
				assert.FieldError(t, err, field, wantErr)
			}
		})
	}
}
