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
				PrimaryKey:   weavetest.SequenceID(1),
				Username:     "Crypt0xxx",
				Bio:          "Best hacker in the universe",
				RegisteredAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"PrimaryKey":   nil,
				"Username":     nil,
				"Bio":          nil,
				"RegisteredAt": nil,
			},
		},
		"failure missing metadata": {
			model: &User{
				PrimaryKey:   weavetest.SequenceID(1),
				Username:     "Crypt0xxx",
				Bio:          "Best hacker in the universe",
				RegisteredAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     errors.ErrMetadata,
				"PrimaryKey":   nil,
				"Username":     nil,
				"Bio":          nil,
				"RegisteredAt": nil,
			},
		},
		"success no bio": {
			model: &User{
				Metadata:     &weave.Metadata{Schema: 1},
				PrimaryKey:   weavetest.SequenceID(1),
				Username:     "Crypt0xxx",
				RegisteredAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"PrimaryKey":   nil,
				"Username":     nil,
				"Bio":          nil,
				"RegisteredAt": nil,
			},
		},
		"failure missing PrimaryKey": {
			model: &User{
				Metadata:     &weave.Metadata{Schema: 1},
				Username:     "Crypt0xxx",
				Bio:          "Best hacker in the universe",
				RegisteredAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"PrimaryKey":   errors.ErrEmpty,
				"Username":     nil,
				"Bio":          nil,
				"RegisteredAt": nil,
			},
		},
		"failure missing username": {
			model: &User{
				Metadata:     &weave.Metadata{Schema: 1},
				PrimaryKey:   weavetest.SequenceID(1),
				Bio:          "Best hacker in the universe",
				RegisteredAt: now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"PrimaryKey":   nil,
				"Username":     errors.ErrModel,
				"Bio":          nil,
				"RegisteredAt": nil,
			},
		},
		"failure missing registered at": {
			model: &User{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				Username:   "Crypt0xxx",
				Bio:        "Best hacker in the universe",
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"PrimaryKey":   nil,
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

func TestDeleteArticleTask(t *testing.T) {
	cases := map[string]struct {
		model    orm.Model
		wantErrs map[string]*errors.Error
	}{
		"success": {
			model: &DeleteArticleTask{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				ArticleID:  weavetest.SequenceID(1),
				TaskOwner:  weavetest.NewCondition().Address(),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"ArticleID":  nil,
				"TaskOwner":  nil,
			},
		},
		"failure missing metadata": {
			model: &DeleteArticleTask{
				PrimaryKey: weavetest.SequenceID(1),
				ArticleID:  weavetest.SequenceID(1),
				TaskOwner:  weavetest.NewCondition().Address(),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"PrimaryKey": nil,
				"ArticleID":  nil,
				"TaskOwner":  nil,
			},
		},
		"failure missing id": {
			model: &DeleteArticleTask{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: weavetest.SequenceID(1),
				TaskOwner: weavetest.NewCondition().Address(),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": errors.ErrEmpty,
				"ArticleID":  nil,
				"TaskOwner":  nil,
			},
		},
		"failure missing article id": {
			model: &DeleteArticleTask{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				TaskOwner:  weavetest.NewCondition().Address(),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"ArticleID":  errors.ErrEmpty,
				"TaskOwner":  nil,
			},
		},
		"failure missing task owner": {
			model: &DeleteArticleTask{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				ArticleID:  weavetest.SequenceID(1),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"ArticleID":  nil,
				"TaskOwner":  errors.ErrEmpty,
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
				PrimaryKey:  weavetest.SequenceID(1),
				Owner:       weavetest.NewCondition().Address(),
				Title:       "Best hacker's blog",
				Description: "Best description ever",
				CreatedAt:   now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
				"CreatedAt":   nil,
			},
		},
		"failure missing metadata": {
			model: &Blog{
				PrimaryKey:  weavetest.SequenceID(1),
				Owner:       weavetest.NewCondition().Address(),
				Title:       "Best hacker's blog",
				Description: "Best description ever",
				CreatedAt:   now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    errors.ErrMetadata,
				"PrimaryKey":  nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
				"CreatedAt":   nil,
			},
		},
		"failure missing PrimaryKey": {
			model: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				Owner:       weavetest.NewCondition().Address(),
				Title:       "Best hacker's blog",
				Description: "Best description ever",
				CreatedAt:   now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  errors.ErrEmpty,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
				"CreatedAt":   nil,
			},
		},
		"failure missing owner": {
			model: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				PrimaryKey:  weavetest.SequenceID(1),
				Title:       "Best hacker's blog",
				Description: "Best description ever",
				CreatedAt:   now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
				"Owner":       errors.ErrEmpty,
				"Title":       nil,
				"Description": nil,
				"CreatedAt":   nil,
			},
		},
		"failure missing title": {
			model: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				PrimaryKey:  weavetest.SequenceID(1),
				Owner:       weavetest.NewCondition().Address(),
				Description: "Best description ever",
				CreatedAt:   now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
				"Owner":       nil,
				"Title":       errors.ErrModel,
				"Description": nil,
				"CreatedAt":   nil,
			},
		},
		"failure missing description": {
			model: &Blog{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				Owner:      weavetest.NewCondition().Address(),
				Title:      "Best hacker's blog",
				CreatedAt:  now,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": errors.ErrModel,
				"CreatedAt":   nil,
			},
		},
		"failure missing created at": {
			model: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				PrimaryKey:  weavetest.SequenceID(1),
				Owner:       weavetest.NewCondition().Address(),
				Title:       "Best hacker's blog",
				Description: "Best description ever",
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
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
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				BlogID:     weavetest.SequenceID(1),
				Owner:      weavetest.NewCondition().Address(),
				Title:      "Best hacker's blog",
				Content:    "Best description ever",
				CreatedAt:  now,
				DeleteAt:   future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogID":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"successs no delete at": {
			model: &Article{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				BlogID:     weavetest.SequenceID(1),
				Owner:      weavetest.NewCondition().Address(),
				Title:      "Best hacker's blog",
				Content:    "Best description ever",
				CreatedAt:  future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogID":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing metadata": {
			model: &Article{
				PrimaryKey: weavetest.SequenceID(1),
				BlogID:     weavetest.SequenceID(1),
				Owner:      weavetest.NewCondition().Address(),
				Title:      "Best hacker's blog",
				Content:    "Best description ever",
				CreatedAt:  now,
				DeleteAt:   future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"PrimaryKey": nil,
				"BlogID":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing PrimaryKey": {
			model: &Article{
				Metadata:  &weave.Metadata{Schema: 1},
				BlogID:    weavetest.SequenceID(1),
				Owner:     weavetest.NewCondition().Address(),
				Title:     "Best hacker's blog",
				Content:   "Best description ever",
				CreatedAt: now,
				DeleteAt:  future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": errors.ErrEmpty,
				"BlogID":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing owner": {
			model: &Article{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				BlogID:     weavetest.SequenceID(1),
				Title:      "Best hacker's blog",
				Content:    "Best description ever",
				CreatedAt:  now,
				DeleteAt:   future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogID":     nil,
				"Owner":      errors.ErrEmpty,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing blog id": {
			model: &Article{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				Owner:      weavetest.NewCondition().Address(),
				Title:      "Best hacker's blog",
				Content:    "Best description ever",
				CreatedAt:  now,
				DeleteAt:   future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogID":     errors.ErrEmpty,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing title": {
			model: &Article{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				BlogID:     weavetest.SequenceID(1),
				Owner:      weavetest.NewCondition().Address(),
				Content:    "Best description ever",
				CreatedAt:  now,
				DeleteAt:   future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogID":     nil,
				"Owner":      nil,
				"Title":      errors.ErrModel,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing content": {
			model: &Article{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				BlogID:     weavetest.SequenceID(1),
				Owner:      weavetest.NewCondition().Address(),
				Title:      "Best hacker's blog",
				CreatedAt:  now,
				DeleteAt:   future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogID":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    errors.ErrModel,
				"DeleteAt":   nil,
			},
		},
		"failure missing created at": {
			model: &Article{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				BlogID:     weavetest.SequenceID(1),
				Owner:      weavetest.NewCondition().Address(),
				Title:      "Best hacker's blog",
				Content:    "Best description ever",
				DeleteAt:   future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogID":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  errors.ErrEmpty,
				"DeleteAt":   nil,
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
