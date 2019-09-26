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
		// TODO add missing metadata test
		"success no bio": {
			model: &User{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				Username:     "Crypt0xxx",
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
		"failure missing ID": {
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
		"failure missing username": {
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
		"failure missing registered at": {
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

func TestValidateBlog(t *testing.T) {
	now := weave.AsUnixTime(time.Now())

	cases := map[string]struct {
		model    orm.Model
		wantErrs map[string]*errors.Error
	}{
		"success": {
			model: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				ID:          weavetest.SequenceID(1),
				Owner:       weavetest.NewCondition().Address(),
				Title:       "Best hacker's blog",
				Description: "Best description ever",
				CreatedAt:   now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
				"CreatedAt":   nil,
			},
		},
		// TODO add missing metadata test
		"failure missing ID": {
			model: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				Owner:       weavetest.NewCondition().Address(),
				Title:       "Best hacker's blog",
				Description: "Best description ever",
				CreatedAt:   now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          errors.ErrEmpty,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
				"CreatedAt":   nil,
			},
		},
		"failure missing owner": {
			model: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				ID:          weavetest.SequenceID(1),
				Title:       "Best hacker's blog",
				Description: "Best description ever",
				CreatedAt:   now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       errors.ErrEmpty,
				"Title":       nil,
				"Description": nil,
				"CreatedAt":   nil,
			},
		},
		"failure missing title": {
			model: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				ID:          weavetest.SequenceID(1),
				Owner:       weavetest.NewCondition().Address(),
				Description: "Best description ever",
				CreatedAt:   now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       nil,
				"Title":       errors.ErrModel,
				"Description": nil,
				"CreatedAt":   nil,
			},
		},
		"failure missing description": {
			model: &Blog{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
				Title:     "Best hacker's blog",
				CreatedAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": errors.ErrModel,
				"CreatedAt":   nil,
			},
		},
		"failure missing created at": {
			model: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				ID:          weavetest.SequenceID(1),
				Owner:       weavetest.NewCondition().Address(),
				Title:       "Best hacker's blog",
				Description: "Best description ever",
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
				"CreatedAt":   errors.ErrEmpty,
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

func TestValidateArticle(t *testing.T) {
	now := weave.AsUnixTime(time.Now())
	future := now.Add(time.Hour)

	cases := map[string]struct {
		model    orm.Model
		wantErrs map[string]*errors.Error
	}{
		"success": {
			model: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				BlogID:       weavetest.SequenceID(1),
				Owner:        weavetest.NewCondition().Address(),
				Title:        "Best hacker's blog",
				Content:      "Best description ever",
				CommentCount: 1,
				LikeCount:    2,
				CreatedAt:    now,
				DeleteAt:     future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"BlogID":       nil,
				"Owner":        nil,
				"Title":        nil,
				"Content":      nil,
				"CommentCount": nil,
				"LikeCount":    nil,
				"CreatedAt":    nil,
				"DeleteAt":     nil,
			},
		},
		"successs no delete at": {
			model: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				BlogID:       weavetest.SequenceID(1),
				Owner:        weavetest.NewCondition().Address(),
				Title:        "Best hacker's blog",
				Content:      "Best description ever",
				CommentCount: 1,
				LikeCount:    2,
				CreatedAt:    future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"BlogID":       nil,
				"Owner":        nil,
				"Title":        nil,
				"Content":      nil,
				"CommentCount": nil,
				"LikeCount":    nil,
				"CreatedAt":    nil,
				"DeleteAt":     nil,
			},
		},
		// TODO add missing metadata test
		"failure missing ID": {
			model: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				BlogID:       weavetest.SequenceID(1),
				Owner:        weavetest.NewCondition().Address(),
				Title:        "Best hacker's blog",
				Content:      "Best description ever",
				CommentCount: 1,
				LikeCount:    2,
				CreatedAt:    now,
				DeleteAt:     future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           errors.ErrEmpty,
				"BlogID":       nil,
				"Owner":        nil,
				"Title":        nil,
				"Content":      nil,
				"CommentCount": nil,
				"LikeCount":    nil,
				"CreatedAt":    nil,
				"DeleteAt":     nil,
			},
		},
		"failure missing owner": {
			model: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				BlogID:       weavetest.SequenceID(1),
				Title:        "Best hacker's blog",
				Content:      "Best description ever",
				CommentCount: 1,
				LikeCount:    2,
				CreatedAt:    now,
				DeleteAt:     future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"BlogID":       nil,
				"Owner":        errors.ErrEmpty,
				"Title":        nil,
				"Content":      nil,
				"CommentCount": nil,
				"LikeCount":    nil,
				"CreatedAt":    nil,
				"DeleteAt":     nil,
			},
		},
		"failure missing blog id": {
			model: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				Owner:        weavetest.NewCondition().Address(),
				Title:        "Best hacker's blog",
				Content:      "Best description ever",
				CommentCount: 1,
				LikeCount:    2,
				CreatedAt:    now,
				DeleteAt:     future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"BlogID":       errors.ErrEmpty,
				"Owner":        nil,
				"Title":        nil,
				"Content":      nil,
				"CommentCount": nil,
				"LikeCount":    nil,
				"CreatedAt":    nil,
				"DeleteAt":     nil,
			},
		},
		"failure missing title": {
			model: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				BlogID:       weavetest.SequenceID(1),
				Owner:        weavetest.NewCondition().Address(),
				Content:      "Best description ever",
				CommentCount: 1,
				LikeCount:    2,
				CreatedAt:    now,
				DeleteAt:     future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"BlogID":       nil,
				"Owner":        nil,
				"Title":        errors.ErrModel,
				"Content":      nil,
				"CommentCount": nil,
				"LikeCount":    nil,
				"CreatedAt":    nil,
				"DeleteAt":     nil,
			},
		},
		"failure missing content": {
			model: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				BlogID:       weavetest.SequenceID(1),
				Owner:        weavetest.NewCondition().Address(),
				Title:        "Best hacker's blog",
				CommentCount: 1,
				LikeCount:    2,
				CreatedAt:    now,
				DeleteAt:     future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"BlogID":       nil,
				"Owner":        nil,
				"Title":        nil,
				"Content":      errors.ErrModel,
				"CommentCount": nil,
				"LikeCount":    nil,
				"CreatedAt":    nil,
				"DeleteAt":     nil,
			},
		},
		"failure negative comment count": {
			model: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				BlogID:       weavetest.SequenceID(1),
				Owner:        weavetest.NewCondition().Address(),
				Title:        "Best hacker's blog",
				Content:      "Best description ever",
				CommentCount: -100,
				LikeCount:    2,
				CreatedAt:    now,
				DeleteAt:     future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"BlogID":       nil,
				"Owner":        nil,
				"Title":        nil,
				"Content":      nil,
				"CommentCount": errors.ErrModel,
				"LikeCount":    nil,
				"CreatedAt":    nil,
				"DeleteAt":     nil,
			},
		},
		"failure negative like count": {
			model: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				BlogID:       weavetest.SequenceID(1),
				Owner:        weavetest.NewCondition().Address(),
				Title:        "Best hacker's blog",
				Content:      "Best description ever",
				CommentCount: 1,
				LikeCount:    -100,
				CreatedAt:    now,
				DeleteAt:     future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"BlogID":       nil,
				"Owner":        nil,
				"Title":        nil,
				"Content":      nil,
				"CommentCount": nil,
				"LikeCount":    errors.ErrModel,
				"CreatedAt":    nil,
				"DeleteAt":     nil,
			},
		},
		"failure missing created at": {
			model: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				BlogID:       weavetest.SequenceID(1),
				Owner:        weavetest.NewCondition().Address(),
				Title:        "Best hacker's blog",
				Content:      "Best description ever",
				CommentCount: 1,
				LikeCount:    2,
				DeleteAt:     future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"ID":           nil,
				"BlogID":       nil,
				"Owner":        nil,
				"Title":        nil,
				"Content":      nil,
				"CommentCount": nil,
				"LikeCount":    nil,
				"CreatedAt":    errors.ErrEmpty,
				"DeleteAt":     nil,
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

func TestValidateComment(t *testing.T) {
	now := weave.AsUnixTime(time.Now())

	cases := map[string]struct {
		model    orm.Model
		wantErrs map[string]*errors.Error
	}{
		"success": {
			model: &Comment{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				ArticleID: weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
				Content:   "Best comment ever",
				CreatedAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"Content":   nil,
				"CreatedAt": nil,
			},
		},
		// TODO add missing metadata test
		"failure missing id": {
			model: &Comment{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
				Content:   "Best comment ever",
				CreatedAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        errors.ErrEmpty,
				"ArticleID": nil,
				"Owner":     nil,
				"Content":   nil,
				"CreatedAt": nil,
			},
		},
		"failure missing article id": {
			model: &Comment{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
				Content:   "Best comment ever",
				CreatedAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": errors.ErrEmpty,
				"Owner":     nil,
				"Content":   nil,
				"CreatedAt": nil,
			},
		},
		"failure missing owner": {
			model: &Comment{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				ArticleID: weavetest.SequenceID(1),
				Content:   "Best comment ever",
				CreatedAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     errors.ErrEmpty,
				"Content":   nil,
				"CreatedAt": nil,
			},
		},
		"failure missing content": {
			model: &Comment{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				ArticleID: weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
				CreatedAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"Content":   errors.ErrModel,
				"CreatedAt": nil,
			},
		},
		"failure missing created at": {
			model: &Comment{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				ArticleID: weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
				Content:   "Best comment ever",
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"Content":   nil,
				"CreatedAt": errors.ErrEmpty,
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

func TestValidateLike(t *testing.T) {
	now := weave.AsUnixTime(time.Now())

	cases := map[string]struct {
		model    orm.Model
		wantErrs map[string]*errors.Error
	}{
		"success": {
			model: &Like{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				ArticleID: weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
				CreatedAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"CreatedAt": nil,
			},
		},
		// TODO add missing metadata test
		"failure missing id": {
			model: &Like{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
				CreatedAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        errors.ErrEmpty,
				"ArticleID": nil,
				"Owner":     nil,
				"CreatedAt": nil,
			},
		},
		"failure missing article id": {
			model: &Like{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
				CreatedAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": errors.ErrEmpty,
				"Owner":     nil,
				"CreatedAt": nil,
			},
		},
		"failure missing owner": {
			model: &Like{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				ArticleID: weavetest.SequenceID(1),
				CreatedAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     errors.ErrEmpty,
				"CreatedAt": nil,
			},
		},
		"failure missing created at": {
			model: &Like{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				ArticleID: weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"CreatedAt": errors.ErrEmpty,
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
