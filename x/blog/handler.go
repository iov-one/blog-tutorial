package blog

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/x"
)

const (
	packageName               = "blog"
	newUserCost         int64 = 1
	newBlogCost         int64 = 10
	changeBlogOwnerCost int64 = 5

	newArticleCost  int64 = 1
	articleCostUnit int64 = 1000 // first 1000 chars are free then pay 1 per mille
	newCommentCost  int64 = 1
)

// RegisterQuery registers buckets for querying.
func RegisterQuery(qr weave.QueryRouter) {
	NewUserBucket().Register("users", qr)
	NewBlogBucket().Register("blogs", qr)
	NewArticleBucket().Register("articles", qr)
}

// RegisterRoutes registers handlers for message processing.
func RegisterRoutes(r weave.Registry, auth x.Authenticator, scheduler weave.Scheduler) {
	//r = migration.SchemaMigratingRegistry(packageName, r)
	r.Handle(&CreateUserMsg{}, NewCreateUserHandler(auth))
	r.Handle(&CreateBlogMsg{}, NewCreateBlogHandler(auth))
	r.Handle(&ChangeBlogOwnerMsg{}, NewChangeBlogOwnerHandler(auth))
	r.Handle(&CreateArticleMsg{}, NewCreateArticleHandler(auth, scheduler))
	r.Handle(&DeleteArticleMsg{}, NewDeleteArticleHandler(auth))
	r.Handle(&CancelDeleteArticleTaskMsg{}, NewCancelDeleteArticleTaskHandler(auth, scheduler))
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
		Metadata:     &weave.Metadata{Schema: 1},
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

	err = h.b.Save(store, user)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store user")
	}

	// Returns generated user PrimaryKey as response
	return &weave.DeliverResult{Data: user.PrimaryKey}, nil
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
		Metadata:    &weave.Metadata{Schema: 1},
		Owner:       x.AnySigner(ctx, h.auth).Address(),
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

	err = h.b.Save(store, blog)
	if err != nil {
		return nil, errors.Wrap(err, "cannot store blog")
	}

	// Returns generated blog PrimaryKey as response
	return &weave.DeliverResult{Data: blog.PrimaryKey}, nil
}

// ------------------- ChangeBlogOwnerHandler -------------------

// ChangeBlogOwnerHandler will handle ChangeBlogOWnerMsg
type ChangeBlogOwnerHandler struct {
	auth x.Authenticator
	b    *BlogBucket
}

var _ weave.Handler = ChangeBlogOwnerHandler{}

// NewChangeBlogOwnerHandler creates a blog message handler
func NewChangeBlogOwnerHandler(auth x.Authenticator) weave.Handler {
	return ChangeBlogOwnerHandler{
		auth: auth,
		b:    NewBlogBucket(),
	}
}

// validate does all common pre-processing between Check and Deliver
func (h ChangeBlogOwnerHandler) validate(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*ChangeBlogOwnerMsg, *Blog, error) {
	var msg ChangeBlogOwnerMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, nil, errors.Wrap(err, "load msg")
	}

	var blog Blog
	if err := h.b.ByID(store, msg.BlogKey, &blog); err != nil {
		return nil, nil, errors.Wrapf(err, "cannot retrieve blog with id %s from database", msg.BlogKey)
	}

	if !h.auth.HasAddress(ctx, blog.Owner) {
		return nil, nil, errors.Wrap(errors.ErrUnauthorized, "only the blog owner can blog owner")
	}

	newBlog := &Blog{
		Metadata:    blog.Metadata,
		Owner:       msg.NewOwner,
		Title:       blog.Title,
		Description: blog.Description,
		CreatedAt:   blog.CreatedAt,
	}

	return &msg, newBlog, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h ChangeBlogOwnerHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, _, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	return &weave.CheckResult{GasAllocated: changeBlogOwnerCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h ChangeBlogOwnerHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	_, blog, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	err = h.b.Save(store, blog)
	if err != nil {
		return nil, errors.Wrap(err, "cannot update blog")
	}

	// Returns generated blog PrimaryKey as response
	return &weave.DeliverResult{Data: blog.PrimaryKey}, nil
}

// ------------------- CreateArticleHandler -------------------

// CreateArticleHandler will handle CreateArticleMsg
type CreateArticleHandler struct {
	auth      x.Authenticator
	ab        *ArticleBucket
	bb        *BlogBucket
	scheduler weave.Scheduler
}

var _ weave.Handler = CreateArticleHandler{}

// NewCreateArticleHandler creates a article message handler
func NewCreateArticleHandler(auth x.Authenticator, scheduler weave.Scheduler) weave.Handler {
	return CreateArticleHandler{
		auth:      auth,
		ab:        NewArticleBucket(),
		bb:        NewBlogBucket(),
		scheduler: scheduler,
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CreateArticleHandler) validate(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*CreateArticleMsg, *Article, error) {
	var msg CreateArticleMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, nil, errors.Wrap(err, "load msg")
	}

	var blog Blog
	if err := h.bb.ByID(store, msg.BlogKey, &blog); err != nil {
		return nil, nil, errors.Wrapf(err, "blog id with %s does not exist", msg.BlogKey)
	}

	if !h.auth.HasAddress(ctx, blog.Owner) {
		return nil, nil, errors.Wrap(errors.ErrUnauthorized, "only the blog owner can post an article under a blog")
	}

	blockTime, err := weave.BlockTime(ctx)
	if err != nil {
		return nil, nil, errors.Wrap(err, "no block time in header")
	}

	if msg.DeleteAt != 0 && weave.InThePast(ctx, msg.DeleteAt.Time()) {
		return nil, nil, errors.Wrap(errors.ErrState, "delete at is in the past")
	}

	now := weave.AsUnixTime(blockTime)

	article := &Article{
		Metadata:  &weave.Metadata{Schema: 1},
		BlogKey:   msg.BlogKey,
		Owner:     blog.Owner,
		Title:     msg.Title,
		Content:   msg.Content,
		CreatedAt: now,
		DeleteAt:  msg.DeleteAt,
	}

	return &msg, article, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CreateArticleHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	msg, _, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	// Calculate gas cost
	gasCost := int64(len(msg.Content)) * newArticleCost / articleCostUnit

	return &weave.CheckResult{GasAllocated: gasCost}, nil
}

// Deliver creates an custom state and saves if all preconditions are met
func (h CreateArticleHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	msg, article, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	// schedule delete task
	if msg.DeleteAt != 0 {
		deleteArticleMsg := &DeleteArticleMsg{
			Metadata:   &weave.Metadata{Schema: 1},
			ArticleKey: article.PrimaryKey,
		}

		taskID, err := h.scheduler.Schedule(store, article.DeleteAt.Time(), nil, deleteArticleMsg)
		if err != nil {
			return nil, errors.Wrap(err, "cannot schedule deletion task")
		}

		article.DeleteTaskID = taskID
	}

	if err := h.ab.Save(store, article); err != nil {
		return nil, errors.Wrap(err, "cannot store article")
	}
	// Returns generated article PrimaryKey as response
	return &weave.DeliverResult{Data: article.PrimaryKey}, nil
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
	if err := h.b.ByID(store, msg.ArticleKey, &article); err != nil {
		return nil, nil, errors.Wrapf(err, "cannot retrieve article with PrimaryKey %s", msg.ArticleKey)
	}

	if !h.auth.HasAddress(ctx, article.Owner) {
		return nil, nil, errors.Wrapf(errors.ErrUnauthorized, "only the article owner can delete the article")
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

	if err := h.b.Delete(store, article.PrimaryKey); err != nil {
		return nil, errors.Wrapf(err, "cannot delete article with PrimaryKey %s", article.PrimaryKey)
	}

	return &weave.DeliverResult{}, nil
}

// ------------------- CancelDeleteArticleTaskHandler -------------------

// CancelDeleteArticleTaskHandler will handle CancelDeleteArticleTaskMsg
type CancelDeleteArticleTaskHandler struct {
	auth      x.Authenticator
	b         *ArticleBucket
	scheduler weave.Scheduler
}

var _ weave.Handler = CancelDeleteArticleTaskHandler{}

// NewCancelDeleteArticleTaskHandler creates a cancel delete article task msg handler
func NewCancelDeleteArticleTaskHandler(auth x.Authenticator, scheduler weave.Scheduler) weave.Handler {
	return CancelDeleteArticleTaskHandler{
		auth:      auth,
		b:         NewArticleBucket(),
		scheduler: scheduler,
	}
}

// validate does all common pre-processing between Check and Deliver
func (h CancelDeleteArticleTaskHandler) validate(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*CancelDeleteArticleTaskMsg, *Article, error) {
	var msg CancelDeleteArticleTaskMsg

	if err := weave.LoadMsg(tx, &msg); err != nil {
		return nil, nil, errors.Wrap(err, "load msg")
	}

	var article Article
	if err := h.b.ByID(store, msg.ArticleKey, &article); err != nil {
		return nil, nil, errors.Wrapf(err, "article with key %s not found", msg.ArticleKey)
	}

	if !h.auth.HasAddress(ctx, article.Owner) {
		return nil, nil, errors.Wrap(errors.ErrUnauthorized, "not authorized to execute this tx")
	}

	if article.DeleteTaskID == nil {
		return nil, nil, errors.Wrap(errors.ErrNotFound, "no scheduled delete task for article")
	}

	return &msg, &article, nil
}

// Check just verifies it is properly formed and returns
// the cost of executing it.
func (h CancelDeleteArticleTaskHandler) Check(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.CheckResult, error) {
	_, _, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	// Cancelling is free of charge
	return &weave.CheckResult{}, nil
}

// Deliver cancels delete task if conditions are met
func (h CancelDeleteArticleTaskHandler) Deliver(ctx weave.Context, store weave.KVStore, tx weave.Tx) (*weave.DeliverResult, error) {
	_, article, err := h.validate(ctx, store, tx)
	if err != nil {
		return nil, err
	}

	if err := h.scheduler.Delete(store, article.DeleteTaskID); err != nil {
		return nil, errors.Wrapf(err, "cannot deschedule with task id %s", article.DeleteTaskID)
	}

	article.DeleteTaskID = nil
	if err := h.b.Save(store, article); err != nil {
		return nil, errors.Wrapf(err, "cannot update article %s ", article.PrimaryKey)
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

	if err := h.b.Delete(store, msg.ArticleKey); err != nil {
		return nil, errors.Wrapf(err, "cannot delete article with PrimaryKey %s", msg.ArticleKey)
	}

	return &weave.DeliverResult{}, nil
}
