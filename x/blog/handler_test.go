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

func TestCreateBlog(t *testing.T) {
	owner := weavetest.NewCondition()

	cases := map[string]struct {
		msg             weave.Msg
		owner           weave.Condition
		expected        *Blog
		wantCheckErrs   map[string]*errors.Error
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &CreateBlogMsg{
				Metadata:    &weave.Metadata{Schema: 1},
				Title:       "insanely good title",
				Description: "best description in the existence",
			},
			owner: owner,
			expected: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				ID:          weavetest.SequenceID(1),
				Owner:       owner.Address(),
				Title:       "insanely good title",
				Description: "best description in the existence",
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
			},
		},
		// TODO add metadata test
		"failure no signer": {
			msg: &CreateBlogMsg{
				Metadata:    &weave.Metadata{Schema: 1},
				Title:       "insanely good title",
				Description: "best description in the existence",
			},
			owner: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       errors.ErrEmpty,
				"Title":       nil,
				"Description": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       errors.ErrEmpty,
				"Title":       nil,
				"Description": nil,
			},
		},
		"failure missing title": {
			msg: &CreateBlogMsg{
				Metadata:    &weave.Metadata{Schema: 1},
				Description: "best description in the existence",
			},
			owner: owner,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       nil,
				"Title":       errors.ErrModel,
				"Description": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       nil,
				"Title":       errors.ErrModel,
				"Description": nil,
			},
		},
		"failure missing description": {
			msg: &CreateBlogMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Title:    "insanely good title",
			},
			owner:    owner,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": errors.ErrModel,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"ID":          nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": errors.ErrModel,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			auth := &weavetest.Auth{
				Signer: tc.owner,
			}

			rt := app.NewRouter()
			RegisterRoutes(rt, auth)
			kv := store.MemStore()
			bucket := NewBlogBucket()

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
				var stored Blog
				err := bucket.One(kv, res.Data, &stored)
				assert.Nil(t, err)

				// ensure createdAt is after test execution starting time
				createdAt := stored.CreatedAt
				weave.InTheFuture(ctx, createdAt.Time())

				// avoid registered at missing error
				tc.expected.CreatedAt = createdAt

				assert.Nil(t, err)
				assert.Equal(t, tc.expected, &stored)
			}
		})
	}
}

func TestCreateArticle(t *testing.T) {
	blogOwner := weavetest.NewCondition()
	signer := weavetest.NewCondition()

	now := weave.AsUnixTime(time.Now())
	past := now.Add(-1 * 5 * time.Hour)
	future := now.Add(time.Hour)

	ownedBlog := &Blog{
		Metadata:    &weave.Metadata{Schema: 1},
		ID:          weavetest.SequenceID(1),
		Owner:       signer.Address(),
		Title:       "Best hacker's blog",
		Description: "Best description ever",
		CreatedAt:   past,
	}
	notOwnedBlog := &Blog{
		Metadata:    &weave.Metadata{Schema: 1},
		ID:          weavetest.SequenceID(2),
		Owner:       blogOwner.Address(),
		Title:       "Worst hacker's blog",
		Description: "Worst description ever",
		CreatedAt:   past,
	}

	cases := map[string]struct {
		msg             weave.Msg
		signer          weave.Condition
		expected        *Article
		wantCheckErrs   map[string]*errors.Error
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   ownedBlog.ID,
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			signer: signer,
			expected: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				BlogID:       ownedBlog.ID,
				Owner:        signer.Address(),
				Title:        "insanely good title",
				Content:      "best content in the existence",
				CommentCount: 0,
				LikeCount:    0,
				CreatedAt:    now,
				DeleteAt:     future,
			},
			wantCheckErrs: map[string]*errors.Error{
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
			wantDeliverErrs: map[string]*errors.Error{
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
		"success no delete at": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   ownedBlog.ID,
				Title:    "insanely good title",
				Content:  "best content in the existence",
			},
			signer: signer,
			expected: &Article{
				Metadata:     &weave.Metadata{Schema: 1},
				ID:           weavetest.SequenceID(1),
				BlogID:       ownedBlog.ID,
				Owner:        signer.Address(),
				Title:        "insanely good title",
				Content:      "best content in the existence",
				CommentCount: 0,
				LikeCount:    0,
				CreatedAt:    now,
			},
			wantCheckErrs: map[string]*errors.Error{
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
			wantDeliverErrs: map[string]*errors.Error{
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
		"failure signer not authorized": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   ownedBlog.ID,
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			signer:   weavetest.NewCondition(),
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
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
			wantDeliverErrs: map[string]*errors.Error{
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
		// TODO add metadata test
		"failure missing blog id": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
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
			wantDeliverErrs: map[string]*errors.Error{
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
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   ownedBlog.ID,
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
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
			wantDeliverErrs: map[string]*errors.Error{
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
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   ownedBlog.ID,
				Title:    "insanely good title",
				DeleteAt: future,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
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
			wantDeliverErrs: map[string]*errors.Error{
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
		"failure blog is not owned by signer": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   notOwnedBlog.ID,
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
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
			wantDeliverErrs: map[string]*errors.Error{
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
		"failure missing signer": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   ownedBlog.ID,
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
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
			wantDeliverErrs: map[string]*errors.Error{
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
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			auth := &weavetest.Auth{
				Signer: tc.signer,
			}

			// initalize environment
			rt := app.NewRouter()
			RegisterRoutes(rt, auth)
			kv := store.MemStore()

			// initalize blog bucket and save blogs
			blogBucket := NewBlogBucket()
			err := blogBucket.Put(kv, ownedBlog)
			assert.Nil(t, err)

			err = blogBucket.Put(kv, notOwnedBlog)
			assert.Nil(t, err)

			// initialize article bucket
			articleBucket := NewArticleBucket()

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
				var stored Article
				err := articleBucket.One(kv, res.Data, &stored)
				assert.Nil(t, err)

				// ensure createdAt is after test execution starting time
				createdAt := stored.CreatedAt
				weave.InTheFuture(ctx, createdAt.Time())

				// avoid registered at missing error
				tc.expected.CreatedAt = createdAt

				assert.Nil(t, err)
				assert.Equal(t, tc.expected, &stored)
			}
		})
	}
}

func TestDeleteArticle(t *testing.T) {
	bob := weavetest.NewCondition()
	signer := weavetest.NewCondition()

	now := weave.AsUnixTime(time.Now())
	future := now.Add(time.Hour)

	ownedArticleID := weavetest.SequenceID(1)
	ownedArticle := &Article{
		Metadata:     &weave.Metadata{Schema: 1},
		ID:           ownedArticleID,
		BlogID:       weavetest.SequenceID(1),
		Owner:        signer.Address(),
		Title:        "Best hacker's blog",
		Content:      "Best description ever",
		CommentCount: 1,
		LikeCount:    2,
		CreatedAt:    now,
		DeleteAt:     future,
	}

	notOwnedArticleID := weavetest.SequenceID(2)
	notOwnedArticle := &Article{
		Metadata:     &weave.Metadata{Schema: 1},
		ID:           notOwnedArticleID,
		BlogID:       weavetest.SequenceID(2),
		Owner:        bob.Address(),
		Title:        "Worst hacker's blog",
		Content:      "Worst description ever",
		CommentCount: 1,
		LikeCount:    2,
		CreatedAt:    now,
		DeleteAt:     future,
	}

	cases := map[string]struct {
		msg             weave.Msg
		signer          weave.Condition
		expected        *Article
		wantCheckErrs   map[string]*errors.Error
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &DeleteArticleMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: ownedArticleID,
			},
			signer:   signer,
			expected: ownedArticle,
			wantCheckErrs: map[string]*errors.Error{
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
			wantDeliverErrs: map[string]*errors.Error{
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
		"failure unauthorized": {
			msg: &DeleteArticleMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: notOwnedArticleID,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
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
			wantDeliverErrs: map[string]*errors.Error{
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
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			auth := &weavetest.Auth{
				Signer: tc.signer,
			}

			// initalize environment
			rt := app.NewRouter()
			RegisterRoutes(rt, auth)
			kv := store.MemStore()

			// initalize article bucket and save articles
			articleBucket := NewArticleBucket()
			err := articleBucket.Put(kv, ownedArticle)
			assert.Nil(t, err)

			err = articleBucket.Put(kv, notOwnedArticle)
			assert.Nil(t, err)

			tx := &weavetest.Tx{Msg: tc.msg}

			ctx := weave.WithBlockTime(context.Background(), time.Now().Round(time.Second))

			if _, err := rt.Check(ctx, kv, tx); err != nil {
				for field, wantErr := range tc.wantCheckErrs {
					assert.FieldError(t, err, field, wantErr)
				}
			}

			_, err = rt.Deliver(ctx, kv, tx)
			if err != nil {
				for field, wantErr := range tc.wantDeliverErrs {
					assert.FieldError(t, err, field, wantErr)
				}
			}

			if tc.expected != nil {
				if err := articleBucket.Has(kv, tc.msg.(*DeleteArticleMsg).ArticleID); err == nil {
					t.Fatalf("got %+v", err)
				}
			}
		})
	}
}
