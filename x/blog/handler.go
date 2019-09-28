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
	newCommentCost int64 = 100
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
	r.Handle(&CreateCommentMsg{}, NewCreateCommentHandler(auth))
	r.Handle(&CreateLikeMsg{}, NewCreateLikeHandler(auth))
}

// RegisterCronRoutes registers routes that are not exposed to
// routers
func RegisterCronRoutes(
	r weave.Registry,
	auth x.Authenticator,
) {
	r.Handle(&DeleteArticleMsg{}, newCronDeleteArticleHandler(auth))
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
func (h CreateUserHandler) validate(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*CreateUserMsg, *User, error) {
	var msg CreateUserMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, nil, errors.Wrap(err, "load msg")
	}

	blockTime, err := weave.BlockTime(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "no block time in header")
	}
	now := weave.AsUnixTime(blockTime)

	user := &User{
		Metadata:     msg.Metadata,
		Username:     msg.Username,
		Bio:          msg.Bio,
		RegisteredAt: now,
	}

	return &msg, user, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateUserHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, _, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newUserCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h CreateUserHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	_, user, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
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
func (h CreateBlogHandler) validate(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*CreateBlogMsg, *Blog, error) {
	var msg CreateBlogMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, nil, errors.Wrap(err, "load msg")
	}

	blockTime, err := weave.BlockTime(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "no block time in header")
	}
	now := weave.AsUnixTime(blockTime)

	blog := &Blog{
		Metadata:    msg.Metadata,
		Owner:       x.MainSigner(ctx, h.auth).Address(),
		Title:       msg.Title,
		Description: msg.Description,
		CreatedAt:   now,
	}

	return &msg, blog, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateBlogHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, _, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newBlogCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h CreateBlogHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	_, blog, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
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
func (h CreateArticleHandler) validate(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*CreateArticleMsg, *Article, error) {
	var msg CreateArticleMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, nil, errors.Wrap(err, "load msg")
	}

	var blog Blog
	if err := h.bb.One(store, msg.BlogID, &blog); err != nil {
		return nil, nil, errors.Wrapf(err, "blog id with %s does not exist", msg.BlogID)
	}

	signer := x.MainSigner(ctx, h.auth).Address()
	if !blog.Owner.Equals(signer) {
		return nil, nil, errors.Wrapf(errors.ErrUnauthorized, "signer %s is unauthorized to post article to the blog with ID %s", signer, blog.ID)
	}

	blockTime, err := weave.BlockTime(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "no block time in header")
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

	return &msg, article, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateArticleHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, _, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newArticleCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h CreateArticleHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	_, article, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
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
func (h DeleteArticleHandler) validate(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*DeleteArticleMsg, *Article, error) {
	var msg DeleteArticleMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, nil, errors.Wrap(err, "load msg")
	}

	var article Article
	if err := h.b.One(store, msg.ArticleID, &article); err != nil {
		return nil, nil, errors.Wrapf(err, "cannot retrieve article with ID %s", msg.ArticleID)
	}

	signer := x.MainSigner(ctx, h.auth).Address()
	if !article.Owner.Equals(signer) {
		return nil, nil, errors.Wrapf(errors.ErrUnauthorized, "signer %s is unauthorized to delete article with ID %s", signer, article.ID)
	}

	return &msg, &article, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h DeleteArticleHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, _, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	// Deleting is free of charge
	return &weave.CheckResult{}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h DeleteArticleHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	_, article, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	if err := h.b.Delete(store, article.ID); err != nil {
		return nil, errors.Wrapf(err, "cannot delete article with ID %s", article.ID)
	}

	return &weave.DeliverResult{}, nil
}

// ------------------- CronDeleteArticleHandler -------------------

// CronDeleteArticleHandler will handle scheduled DeleteArticleMsg
type CronDeleteArticleHandler struct {
	auth x.Authenticator
	b    *ArticleBucket
}

var _ weave.Handler = CronDeleteArticleHandler{}

// newCronDeleteArticleHandler creates a article message handler
func newCronDeleteArticleHandler(auth x.Authenticator) weave.Handler {
	return CronDeleteArticleHandler{
		auth: auth,
		b:    NewArticleBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CronDeleteArticleHandler) validate(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*DeleteArticleMsg, error) {
	var msg DeleteArticleMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, errors.Wrap(err, "load msg")
	}

	return &msg, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CronDeleteArticleHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	// Deleting is free of charge
	return &weave.CheckResult{}, nil
}

// Deliver stages a scheduled deletion if all preconditions are met
func (h CronDeleteArticleHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	msg, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	if err := h.b.Delete(store, msg.ArticleID); err != nil {
		return nil, errors.Wrapf(err, "cannot delete article with ID %s", msg.ArticleID)
	}

	return &weave.DeliverResult{}, nil
}

// ------------------- CreateCommentHandler -------------------

// CreateCommentHandler will handle CreateCommentMsg
type CreateCommentHandler struct {
	auth x.Authenticator
	cb   *CommentBucket
	ab   *ArticleBucket
}

var _ weave.Handler = CreateCommentHandler{}

// NewCreateCommentHandler creates a comment message handler
func NewCreateCommentHandler(auth x.Authenticator) weave.Handler {
	return CreateCommentHandler{
		auth: auth,
		cb:   NewCommentBucket(),
		ab:   NewArticleBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CreateCommentHandler) validate(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*CreateCommentMsg, *Comment, error) {
	var msg CreateCommentMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, nil, errors.Wrap(err, "load msg")
	}

	// Check if article exists
	if err := h.ab.Has(store, msg.ArticleID); err != nil {
		return nil, nil, errors.Wrapf(err, "article with id %s does not exist", msg.ArticleID)
	}

	signer := x.MainSigner(ctx, h.auth).Address()

	blockTime, err := weave.BlockTime(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "no block time in header")
	}
	now := weave.AsUnixTime(blockTime)

	comment := &Comment{
		Metadata:  msg.Metadata,
		ArticleID: msg.ArticleID,
		Owner:     signer,
		Content:   msg.Content,
		CreatedAt: now,
	}

	return &msg, comment, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateCommentHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, _, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newCommentCost}, nil
}

// Deliver creates a comment and saves if all preconditions are met
func (h CreateCommentHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	_, comment, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	err = h.cb.Put(store, comment)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store comment")
	}

	// Returns generated user ID as response
	return &weave.DeliverResult{Data: comment.ID}, nil
}

// ------------------- CreateLikeHandler -------------------

// CreateLikeHander will handle CreateLikeMsg
type CreateLikeHandler struct {
	auth x.Authenticator
	ab   *ArticleBucket
	lb   *LikeBucket
}

var _ weave.Handler = CreateLikeHandler{}

// NewCreateLikeHandler creates a like message handler
func NewCreateLikeHandler(auth x.Authenticator) weave.Handler {
	return CreateLikeHandler{
		auth: auth,
		ab:   NewArticleBucket(),
		lb:   NewLikeBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CreateLikeHandler) validate(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*CreateLikeMsg, *Like, error) {
	var msg CreateLikeMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, nil, errors.Wrap(err, "load msg")
	}

	// Check if article exists
	if err := h.ab.Has(store, msg.ArticleID); err != nil {
		return nil, nil, errors.Wrapf(err, "article with id %s does not exist", msg.ArticleID)
	}

	signer := x.MainSigner(ctx, h.auth).Address()

	blockTime, err := weave.BlockTime(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "no block time in header")
	}
	now := weave.AsUnixTime(blockTime)

	like := &Like{
		Metadata:  msg.Metadata,
		ArticleID: msg.ArticleID,
		Owner:     signer,
		CreatedAt: now,
	}

	return &msg, like, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateLikeHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, _, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: newCommentCost}, nil
}

// Deliver creates a like and saves if all preconditions are met
func (h CreateLikeHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	_, like, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	err = h.lb.Put(store, like)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store like")
	}

	// Returns generated like ID as response
	return &weave.DeliverResult{Data: like.ID}, nil
}
