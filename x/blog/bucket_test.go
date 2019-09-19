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

func TestBlogIDIndexer(t *testing.T) {
	now := weave.AsUnixTime(time.Now())
	future := now.Add(time.Hour)

	blogID := weavetest.SequenceID(1)

	article := &Article{
		Metadata:     &weave.Metadata{Schema: 1},
		ID:           weavetest.SequenceID(1),
		BlogID:       blogID,
		Title:        "Best hacker's blog",
		Content:      "Best description ever",
		CommentCount: 1,
		LikeCount:    2,
		CreatedAt:    now,
		DeleteAt:     future,
	}

	cases := map[string]struct {
		obj      orm.Object
		expected []byte
		wantErr  *errors.Error
	}{
		"success": {
			obj:      orm.NewSimpleObj(nil, article),
			expected: blogID,
			wantErr:  nil,
		},
		"failure, obj is nil": {
			obj:      nil,
			expected: nil,
			wantErr:  nil,
		},
		"not article": {
			obj:      orm.NewSimpleObj(nil, new(Blog)),
			expected: nil,
			wantErr:  errors.ErrState,
		},
	}

	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			index, err := blogIDIndexer(tc.obj)

			if !tc.wantErr.Is(err) {
				t.Fatalf("unexpected error: %+v", err)
			}
			assert.Equal(t, tc.expected, index)
		})
	}
}
