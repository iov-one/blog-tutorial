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
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				Username:   "Crpto0X",
				Bio:        "Best hacker in the universe",
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
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				Username:   "Crpto0X",
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
		"failure missing metadata": {
			msg: &CreateUserMsg{
				Username: "Crpto0X",
			},
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": errors.ErrMetadata,
				"Username": nil,
				"Bio":      nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": errors.ErrMetadata,
				"Username": nil,
				"Bio":      nil,
			},
		},
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

			scheduler := &weavetest.Cron{}
			RegisterRoutes(rt, auth, scheduler)

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
				if res.Data == nil {
					t.Fatalf("no data")
				}
				var stored User
				err := bucket.ByID(kv, res.Data, &stored)
				if err != nil {
					t.Fatalf("unexpected error: %+v", err)
				}

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
				PrimaryKey:  weavetest.SequenceID(1),
				Owner:       owner.Address(),
				Title:       "insanely good title",
				Description: "best description in the existence",
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
			},
		},
		"failure missing metadata": {
			msg: &CreateBlogMsg{
				Metadata:    &weave.Metadata{Schema: 1},
				Title:       "insanely good title",
				Description: "best description in the existence",
			},
			owner: owner,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":    errors.ErrMetadata,
				"PrimaryKey":  nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":    errors.ErrMetadata,
				"PrimaryKey":  nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": nil,
			},
		},
		"failure no signer": {
			msg: &CreateBlogMsg{
				Metadata:    &weave.Metadata{Schema: 1},
				Title:       "insanely good title",
				Description: "best description in the existence",
			},
			owner: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
				"Owner":       errors.ErrEmpty,
				"Title":       nil,
				"Description": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
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
				"PrimaryKey":  nil,
				"Owner":       nil,
				"Title":       errors.ErrModel,
				"Description": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
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
				"PrimaryKey":  nil,
				"Owner":       nil,
				"Title":       nil,
				"Description": errors.ErrModel,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":    nil,
				"PrimaryKey":  nil,
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

			scheduler := &weavetest.Cron{}
			RegisterRoutes(rt, auth, scheduler)

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
				err := bucket.ByID(kv, res.Data, &stored)
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
		PrimaryKey:  blogID,
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
				BlogKey:   blogID,
				NewOwner: newOwner.Address(),
			},
			owner: owner,
			expected: &Blog{
				Metadata:    &weave.Metadata{Schema: 1},
				PrimaryKey:  blogID,
				Owner:       newOwner.Address(),
				Title:       "insanely good title",
				Description: "best description in the existence",
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogKey":   nil,
				"NewOwner": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogKey":   nil,
				"NewOwner": nil,
			},
		},
		"failure missing metadata": {
			msg: &ChangeBlogOwnerMsg{
				BlogKey:   blogID,
				NewOwner: newOwner.Address(),
			},
			owner:    owner,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": errors.ErrMetadata,
				"BlogKey":   nil,
				"NewOwner": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": errors.ErrMetadata,
				"BlogKey":   nil,
				"NewOwner": nil,
			},
		},
		"failure signer does not own the blog": {
			msg: &ChangeBlogOwnerMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogKey:   blogID,
				NewOwner: newOwner.Address(),
			},
			owner:    weavetest.NewCondition(),
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogKey":   nil,
				"NewOwner": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogKey":   nil,
				"NewOwner": nil,
			},
		},
		"failure invalid owner": {
			msg: &ChangeBlogOwnerMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogKey:   blogID,
				NewOwner: []byte{0, 0},
			},
			owner:    owner,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogKey":   nil,
				"NewOwner": errors.ErrInput,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogKey":   nil,
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
				"BlogKey":   errors.ErrEmpty,
				"NewOwner": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogKey":   errors.ErrEmpty,
				"NewOwner": nil,
			},
		},
		"failure invalid blog id": {
			msg: &ChangeBlogOwnerMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogKey:   []byte{0, 0},
				NewOwner: newOwner.Address(),
			},
			owner:    owner,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogKey":   errors.ErrInput,
				"NewOwner": nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata": nil,
				"BlogKey":   errors.ErrInput,
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

			scheduler := &weavetest.Cron{}
			RegisterRoutes(rt, auth, scheduler)

			kv := store.MemStore()
			bucket := NewBlogBucket()

			err := bucket.Save(kv, blog)
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
				err := bucket.ByID(kv, res.Data, &stored)
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
		PrimaryKey:  weavetest.SequenceID(1),
		Owner:       signer.Address(),
		Title:       "Best hacker's blog",
		Description: "Best description ever",
		CreatedAt:   now,
	}
	notOwnedBlog := &Blog{
		Metadata:    &weave.Metadata{Schema: 1},
		PrimaryKey:  weavetest.SequenceID(2),
		Owner:       blogOwner.Address(),
		Title:       "Worst hacker's blog",
		Description: "Worst description ever",
		CreatedAt:   now,
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
				BlogKey:   ownedBlog.PrimaryKey,
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			signer: signer,
			expected: &Article{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				BlogKey:     ownedBlog.PrimaryKey,
				Owner:      signer.Address(),
				Title:      "insanely good title",
				Content:    "best content in the existence",
				CreatedAt:  now,
				DeleteAt:   future,
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"success no delete at": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogKey:   ownedBlog.PrimaryKey,
				Title:    "insanely good title",
				Content:  "best content in the existence",
			},
			signer: signer,
			expected: &Article{
				Metadata:   &weave.Metadata{Schema: 1},
				PrimaryKey: weavetest.SequenceID(1),
				BlogKey:     ownedBlog.PrimaryKey,
				Owner:      signer.Address(),
				Title:      "insanely good title",
				Content:    "best content in the existence",
				CreatedAt:  now,
			},
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing metadata": {
			msg: &CreateArticleMsg{
				BlogKey:   ownedBlog.PrimaryKey,
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure signer not authorized": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogKey:   ownedBlog.PrimaryKey,
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			signer:   weavetest.NewCondition(),
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
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
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     errors.ErrEmpty,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     errors.ErrEmpty,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing title": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogKey:   ownedBlog.PrimaryKey,
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      errors.ErrModel,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      errors.ErrModel,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing content": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogKey:   ownedBlog.PrimaryKey,
				Title:    "insanely good title",
				DeleteAt: future,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    errors.ErrModel,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"PrimaryKey":   nil,
				"BlogKey":       nil,
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
				BlogKey:   notOwnedBlog.PrimaryKey,
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing signer": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogKey:   ownedBlog.PrimaryKey,
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: future,
			},
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure delete at in the past": {
			msg: &CreateArticleMsg{
				Metadata: &weave.Metadata{Schema: 1},
				BlogKey:   ownedBlog.PrimaryKey,
				Title:    "insanely good title",
				Content:  "best content in the existence",
				DeleteAt: past,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
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

			scheduler := &weavetest.Cron{}
			RegisterRoutes(rt, auth, scheduler)

			kv := store.MemStore()

			// initalize blog bucket and save blogs
			blogBucket := NewBlogBucket()
			err := blogBucket.Save(kv, ownedBlog)
			assert.Nil(t, err)

			err = blogBucket.Save(kv, notOwnedBlog)
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
				err := articleBucket.ByID(kv, res.Data, &stored)
				if err != nil {
					t.Fatalf("unexpected error: %+v", err)
				}

				// ensure createdAt is after test execution starting time
				createdAt := stored.CreatedAt
				weave.InTheFuture(ctx, createdAt.Time())

				// avoid registered at missing error
				tc.expected.CreatedAt = createdAt
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
		Metadata:   &weave.Metadata{Schema: 1},
		PrimaryKey: ownedArticleID,
		BlogKey:     weavetest.SequenceID(1),
		Owner:      signer.Address(),
		Title:      "Best hacker's blog",
		Content:    "Best description ever",
		CreatedAt:  now,
		DeleteAt:   future,
	}

	notOwnedArticleID := weavetest.SequenceID(2)
	notOwnedArticle := &Article{
		Metadata:   &weave.Metadata{Schema: 1},
		PrimaryKey: notOwnedArticleID,
		BlogKey:     weavetest.SequenceID(2),
		Owner:      bob.Address(),
		Title:      "Worst hacker's blog",
		Content:    "Worst description ever",
		CreatedAt:  now,
		DeleteAt:   future,
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
				ArticleKey: ownedArticleID,
			},
			signer:   signer,
			expected: ownedArticle,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing metadata": {
			msg: &DeleteArticleMsg{
				ArticleKey: notOwnedArticleID,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure unauthorized": {
			msg: &DeleteArticleMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleKey: notOwnedArticleID,
			},
			signer:   signer,
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
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

			scheduler := &weavetest.Cron{}
			RegisterRoutes(rt, auth, scheduler)

			kv := store.MemStore()

			// initalize article bucket and save articles
			articleBucket := NewArticleBucket()
			err := articleBucket.Save(kv, ownedArticle)
			assert.Nil(t, err)

			err = articleBucket.Save(kv, notOwnedArticle)
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
				if err := articleBucket.Has(kv, tc.msg.(*DeleteArticleMsg).ArticleKey); err == nil {
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
		Metadata:   &weave.Metadata{Schema: 1},
		PrimaryKey: articleID,
		BlogKey:     weavetest.SequenceID(1),
		Owner:      weavetest.NewCondition().Address(),
		Title:      "Best hacker's blog",
		Content:    "Best description ever",
		CreatedAt:  now,
		DeleteAt:   future,
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
				ArticleKey: articleID,
			},
			expected: article,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"failure missing metadata": {
			msg: &DeleteArticleMsg{
				ArticleKey: articleID,
			},
			expected: article,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
		},
		"success article already deleted": {
			msg: &DeleteArticleMsg{
				Metadata:  &weave.Metadata{Schema: 1},
				ArticleKey: notExistingArticleID,
			},
			expected: nil,
			wantCheckErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
			},
			wantDeliverErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"PrimaryKey": nil,
				"BlogKey":     nil,
				"Owner":      nil,
				"Title":      nil,
				"Content":    nil,
				"CreatedAt":  nil,
				"DeleteAt":   nil,
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
			err := articleBucket.Save(kv, article)
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
				if err := articleBucket.Has(kv, tc.msg.(*DeleteArticleMsg).ArticleKey); err != nil && !errors.ErrNotFound.Is(err) {
					t.Fatal("article still exists")
				}
			}
		})
	}
}
