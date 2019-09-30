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

func TestChangeOwner(t *testing.T) {
	owner := weavetest.NewCondition()
	newOwner := weavetest.NewCondition()

	blogID := weavetest.SequenceID(1)

	blog := &Blog{
		Metadata:    &weave.Metadata{Schema: 1},
		ID:          blogID,
		Owner:       owner.Address(),
		Title:       "insanely good title",
		Description: "best description in the existence",
		CreatedAt:   weave.AsUnixTime(time.Now()),
	}

	cases := map[string]struct {
		msg             weave.Msg
		owner           weave.Condition
		expected        *Blog
		wantCheckErrs   map[string]*errors.Error
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &ChangeBlogOwnerMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   blogID,
				NewOwner: newOwner.Address(),
			},
			owner: owner,
			expected: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				ID:          blogID,
				Owner:       newOwner.Address(),
				Title:       "insanely good title",
				Description: "best description in the existence",
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogID":   nil,
				"NewOwner": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogID":   nil,
				"NewOwner": nil,
			},
		},
		// TODO add metadata test
		"failure signer does not own the blog": {
			msg: &ChangeBlogOwnerMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   blogID,
				NewOwner: newOwner.Address(),
			},
			owner:    weavetest.NewCondition(),
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogID":   nil,
				"NewOwner": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogID":   nil,
				"NewOwner": nil,
			},
		},
		"failure invalid owner": {
			msg: &ChangeBlogOwnerMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   blogID,
				NewOwner: []byte{0, 0},
			},
			owner:    owner,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogID":   nil,
				"NewOwner": errors.ErrInput,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogID":   nil,
				"NewOwner": errors.ErrInput,
			},
		},
		"failure missing blog id": {
			msg: &ChangeBlogOwnerMsg{
				Metadata: &weave.Metadata{Schema: 1},
				NewOwner: newOwner.Address(),
			},
			owner:    owner,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogID":   errors.ErrEmpty,
				"NewOwner": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogID":   errors.ErrEmpty,
				"NewOwner": nil,
			},
		},
		"failure invalid blog id": {
			msg: &ChangeBlogOwnerMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogID:   []byte{0, 0},
				NewOwner: newOwner.Address(),
			},
			owner:    owner,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogID":   errors.ErrInput,
				"NewOwner": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogID":   errors.ErrInput,
				"NewOwner": nil,
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

			err := bucket.Put(kv, blog)
			assert.Nil(t, err)

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

func TestCronDeleteArticle(t *testing.T) {
	now := weave.AsUnixTime(time.Now())
	future := now.Add(time.Hour)

	articleID := weavetest.SequenceID(1)
	article := &Article{
		Metadata:     &weave.Metadata{Schema: 1},
		ID:           articleID,
		BlogID:       weavetest.SequenceID(1),
		Owner:        weavetest.NewCondition().Address(),
		Title:        "Best hacker's blog",
		Content:      "Best description ever",
		CommentCount: 1,
		LikeCount:    2,
		CreatedAt:    now,
		DeleteAt:     future,
	}

	notExistingArticleID := weavetest.SequenceID(2)

	cases := map[string]struct {
		msg             weave.Msg
		expected        *Article
		wantCheckErrs   map[string]*errors.Error
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &DeleteArticleMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: articleID,
			},
			expected: article,
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
		"success article already deleted": {
			msg: &DeleteArticleMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: notExistingArticleID,
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
			auth := &weavetest.Auth{}

			// initalize environment
			rt := app.NewRouter()
			RegisterCronRoutes(rt, auth)
			kv := store.MemStore()

			// initalize article bucket and save articles
			articleBucket := NewArticleBucket()
			err := articleBucket.Put(kv, article)
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
				if err := articleBucket.Has(kv, tc.msg.(*DeleteArticleMsg).ArticleID); err != nil && !errors.ErrNotFound.Is(err) {
					t.Fatal("article still exists")
				}
			}
		})
	}
}

func TestCreateComment(t *testing.T) {
	now := weave.AsUnixTime(time.Now())

	commentOwner := weavetest.NewCondition()
	commentCount := int64(0)

	existingArticleID := weavetest.SequenceID(1)
	notExistingArticleID := weavetest.SequenceID(2)

	article := &Article{
		Metadata:     &weave.Metadata{Schema: 1},
		ID:           existingArticleID,
		BlogID:       weavetest.SequenceID(1),
		Owner:        weavetest.NewCondition().Address(),
		Title:        "Best hacker's blog",
		Content:      "Best description ever",
		CommentCount: commentCount,
		LikeCount:    2,
		CreatedAt:    now,
	}

	cases := map[string]struct {
		msg             weave.Msg
		expected        *Comment
		newCommentCount int64
		wantCheckErrs   map[string]*errors.Error
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &CreateCommentMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: existingArticleID,
				Content:   "best content in the existence",
			},
			expected: &Comment{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				ArticleID: existingArticleID,
				Owner:     commentOwner.Address(),
				Content:   "best content in the existence",
			},
			newCommentCount: commentCount + 1,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"Content":   nil,
				"CreatedAt": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"Content":   nil,
				"CreatedAt": nil,
			},
		},
		"failure article does not exist": {
			msg: &CreateCommentMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: notExistingArticleID,
				Content:   "best content in the existence",
			},
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"Content":   nil,
				"CreatedAt": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"Content":   nil,
				"CreatedAt": nil,
			},
		},
		// TODO add metadata test
		"failure content is missing": {
			msg: &CreateCommentMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: existingArticleID,
			},
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"Content":   errors.ErrModel,
				"CreatedAt": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"Content":   errors.ErrModel,
				"CreatedAt": nil,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			auth := &weavetest.Auth{
				Signer: commentOwner,
			}

			rt := app.NewRouter()
			RegisterRoutes(rt, auth)
			kv := store.MemStore()

			articleBucket := NewArticleBucket()

			err := articleBucket.Put(kv, article)
			assert.Nil(t, err)

			commentBucket := NewCommentBucket()

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
				var storedComment Comment
				err := commentBucket.One(kv, res.Data, &storedComment)
				assert.Nil(t, err)
				tc.expected.CreatedAt = storedComment.CreatedAt
				assert.Equal(t, tc.expected, &storedComment)

				// assert if comment count is increased on article
				var storedArticle Article
				err = articleBucket.One(kv, existingArticleID, &storedArticle)
				assert.Nil(t, err)

				if storedArticle.CommentCount != tc.newCommentCount {
					t.Errorf("expected %d but got %d", tc.newCommentCount, article.CommentCount)
				}
			}
		})
	}
}

func TestCreateLike(t *testing.T) {
	existingArticleID := weavetest.SequenceID(1)
	notExistingArticleID := weavetest.SequenceID(2)

	likeOwner := weavetest.NewCondition()

	now := weave.AsUnixTime(time.Now())

	likeCount := int64(1)
	article := &Article{
		Metadata:     &weave.Metadata{Schema: 1},
		ID:           existingArticleID,
		BlogID:       weavetest.SequenceID(1),
		Owner:        weavetest.NewCondition().Address(),
		Title:        "Best hacker's blog",
		Content:      "Best description ever",
		CommentCount: 1,
		LikeCount:    likeCount,
		CreatedAt:    now,
	}

	cases := map[string]struct {
		msg             weave.Msg
		expected        *Like
		wantCheckErrs   map[string]*errors.Error
		newLikeCount    int64
		wantDeliverErrs map[string]*errors.Error
	}{
		"success": {
			msg: &CreateLikeMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: existingArticleID,
			},
			expected: &Like{
				Metadata:  &weave.Metadata{Schema: 1},
				ID:        weavetest.SequenceID(1),
				ArticleID: existingArticleID,
				Owner:     likeOwner.Address(),
			},
			newLikeCount: likeCount + 1,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"CreatedAt": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"CreatedAt": nil,
			},
		},
		// TODO add missing metadata test
		"failure missing article id": {
			msg: &CreateLikeMsg{
				Metadata: &weave.Metadata{Schema: 1},
			},
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": errors.ErrEmpty,
				"Owner":     nil,
				"CreatedAt": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": errors.ErrEmpty,
				"Owner":     nil,
				"CreatedAt": nil,
			},
		},
		"failure article does not exist": {
			msg: &CreateLikeMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleID: notExistingArticleID,
			},
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"CreatedAt": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":  nil,
				"ID":        nil,
				"ArticleID": nil,
				"Owner":     nil,
				"CreatedAt": nil,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			auth := &weavetest.Auth{
				Signer: likeOwner,
			}

			rt := app.NewRouter()
			RegisterRoutes(rt, auth)
			kv := store.MemStore()

			articleBucket := NewArticleBucket()

			err := articleBucket.Put(kv, article)
			assert.Nil(t, err)

			likeBucket := NewLikeBucket()

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
				var like Like
				err := likeBucket.One(kv, res.Data, &like)
				assert.Nil(t, err)
				tc.expected.CreatedAt = like.CreatedAt
				assert.Equal(t, tc.expected, &like)

				// test if articles like count is increased
				var article Article
				err = articleBucket.One(kv, existingArticleID, &article)
				assert.Nil(t, err)

				if article.LikeCount != tc.newLikeCount {
					t.Errorf("expected %d but got %d", tc.newLikeCount, article.LikeCount)
				}
			}
		})
	}
}
