package blog

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/x"
)

const (
	packageName          = "blog"
	newUserCost    int64 = 100
	newBlogCost    int64 = 100
	newArticleCost int64 = 100
)

// RegisterQuery registers buckets for querying.
func RegisterQuery(qr weave.QueryRouter) {
}

// RegisterRoutes registers handlers for message processing.
func RegisterRoutes(r weave.Registry, auth x.Authenticator) {
	//r = migration.SchemaMigratingRegistry(packageName, r)
	r.Handle(&CreateUserMsg{}, NewCreateUserHandler(auth))
	r.Handle(&CreateBlogMsg{}, NewCreateBlogHandler(auth))
	r.Handle(&CreateArticleMsg{}, NewCreateArticleHandler(auth))
	r.Handle(&DeleteArticleMsg{}, NewDeleteArticleHandler(auth))
}

// ------------------- CreateUserHandler -------------------

// CreateUserHandler will handle CreateUserMsg
type CreateUserHandler struct {
	auth x.Authenticator
	b    *UserBucket
}

var _ weave.Handler = CreateUserHandler{}

// NewCreateUserHandler creates a user message handler
func NewCreateUserHandler(auth x.Authenticator) weave.Handler {
	return CreateUserHandler{
		auth: auth,
		b:    NewUserBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CreateUserHandler) validate(ctx weave.Context, db weave.KVStore, tx weave.Tx) (*CreateUserMsg, error) {
	var msg CreateUserMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, errors.Wrap(err, "load msg")
	}

	return &msg, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateUserHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newUserCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h CreateUserHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	msg, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	blockTime, err := weave.BlockTime(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "no block time in header")
	}
	now := weave.AsUnixTime(blockTime)

	user := &User{
		Metadata:     msg.Metadata,
		Username:     msg.Username,
		Bio:          msg.Bio,
		RegisteredAt: now,
	}

	err = h.b.Put(store, user)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store user")
	}

	// Returns generated user ID as response
	return &weave.DeliverResult{Data: user.ID}, nil
}

// ------------------- CreateBlogHandler -------------------

// CreateBlogHandler will handle CreateBlogMsg
type CreateBlogHandler struct {
	auth x.Authenticator
	b    *BlogBucket
}

var _ weave.Handler = CreateBlogHandler{}

// NewCreateBlogHandler creates a blog message handler
func NewCreateBlogHandler(auth x.Authenticator) weave.Handler {
	return CreateBlogHandler{
		auth: auth,
		b:    NewBlogBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CreateBlogHandler) validate(ctx weave.Context, db weave.KVStore, tx weave.Tx) (*CreateBlogMsg, error) {
	var msg CreateBlogMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, errors.Wrap(err, "load msg")
	}

	return &msg, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateBlogHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newBlogCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h CreateBlogHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	msg, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	blockTime, err := weave.BlockTime(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "no block time in header")
	}
	now := weave.AsUnixTime(blockTime)

	blog := &Blog{
		Metadata:    msg.Metadata,
		Owner:       x.MainSigner(ctx, h.auth).Address(),
		Title:       msg.Title,
		Description: msg.Description,
		CreatedAt:   now,
	}

	err = h.b.Put(store, blog)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store blog")
	}

	// Returns generated blog ID as response
	return &weave.DeliverResult{Data: blog.ID}, nil
}

// ------------------- CreateArticleHandler -------------------

// CreateArticleHandler will handle CreateArticleMsg
type CreateArticleHandler struct {
	auth x.Authenticator
	ab   *ArticleBucket
	bb   *BlogBucket
}

var _ weave.Handler = CreateArticleHandler{}

// NewCreateArticleHandler creates a article message handler
func NewCreateArticleHandler(auth x.Authenticator) weave.Handler {
	return CreateArticleHandler{
		auth: auth,
		ab:   NewArticleBucket(),
		bb:   NewBlogBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CreateArticleHandler) validate(ctx weave.Context, db weave.KVStore, tx weave.Tx) (*CreateArticleMsg, error) {
	var msg CreateArticleMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, errors.Wrap(err, "load msg")
	}

	return &msg, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateArticleHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newArticleCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h CreateArticleHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	msg, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	var blog Blog
	h.bb.One(store, msg.BlogID, &blog)

	signer := x.MainSigner(ctx, h.auth).Address()
	if !blog.Owner.Equals(signer) {
		return nil, errors.Wrapf(errors.ErrUnauthorized, "signer %s is unauthorized to post article to the blog with ID %s", signer, blog.ID)
	}

	blockTime, err := weave.BlockTime(ctx)
	if err != nil {
		return nil, errors.Wrap(err, "no block time in header")
	}
	now := weave.AsUnixTime(blockTime)

	article := &Article{
		Metadata:     msg.Metadata,
		BlogID:       msg.BlogID,
		Owner:        signer,
		Title:        msg.Title,
		Content:      msg.Content,
		CommentCount: 0,
		LikeCount:    0,
		CreatedAt:    now,
		DeleteAt:     msg.DeleteAt,
	}

	err = h.ab.Put(store, article)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store article")
	}

	// Returns generated article ID as response
	return &weave.DeliverResult{Data: article.ID}, nil
}

// ------------------- DeleteArticleHandler -------------------

// DeleteArticleHandler will handle DeleteArticleMsg
type DeleteArticleHandler struct {
	auth x.Authenticator
	b    *ArticleBucket
}

var _ weave.Handler = DeleteArticleHandler{}

// NewDeleteArticleHandler creates a article message handler
func NewDeleteArticleHandler(auth x.Authenticator) weave.Handler {
	return DeleteArticleHandler{
		auth: auth,
		b:    NewArticleBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h DeleteArticleHandler) validate(ctx weave.Context, db weave.KVStore, tx weave.Tx) (*DeleteArticleMsg, error) {
	var msg DeleteArticleMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, errors.Wrap(err, "load msg")
	}

	return &msg, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h DeleteArticleHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	// Deleting is free of charge
	return &weave.CheckResult{}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h DeleteArticleHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	msg, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	var article Article
	if err := h.b.One(store, msg.ArticleID, &article); err != nil {
		return nil, errors.Wrapf(err, "cannot retrieve article with ID %s", msg.ArticleID)
	}

	signer := x.MainSigner(ctx, h.auth).Address()
	if !article.Owner.Equals(signer) {
		return nil, errors.Wrapf(errors.ErrUnauthorized, "signer %s is unauthorized to delete article with ID %s", signer, article.ID)
	}

	if err := h.b.Delete(store, article.ID); err != nil {
		return nil, errors.Wrapf(err, "cannot delete article with ID %s", article.ID)
	}

	return &weave.DeliverResult{}, nil
}
