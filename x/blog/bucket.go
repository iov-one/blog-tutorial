package blog

import (
	"github.com/iov-one/blog-tutorial/morm"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/orm"
)

type UserBucket struct {
	morm.ModelBucket
}

// NewUserBucket returns a new user bucket
func NewUserBucket() *UserBucket {
	return &UserBucket{
		morm.NewModelBucket("user", &User{}),
	}
}

type BlogBucket struct {
	morm.ModelBucket
}

// NewBlogBucket returns a new blog bucket
func NewBlogBucket() *BlogBucket {
	return &BlogBucket{
		morm.NewModelBucket("blog", &Blog{},
			morm.WithIndex("user", blogUserIDIndexer, false)),
	}
}

// userIDIndexer enables querying blogs by user ids
func blogUserIDIndexer(obj orm.Object) ([]byte, error) {
	if obj == nil || obj.Value() == nil {
		return nil, nil
	}
	blog, ok := obj.Value().(*Blog)
	if !ok {
		return nil, errors.Wrapf(errors.ErrState, "expected blog, got %T", obj.Value())
	}
	return blog.Owner, nil
}

type ArticleBucket struct {
	morm.ModelBucket
}

// NewArticleBucket returns a new article bucket
func NewArticleBucket() *ArticleBucket {
	return &ArticleBucket{
		morm.NewModelBucket("article", &Article{},
			morm.WithIndex("blog", articleBlogIDIndexer, false)),
	}
}

// articleBlogIDIndexer enables querying articles by blog ids
func articleBlogIDIndexer(obj orm.Object) ([]byte, error) {
	if obj == nil || obj.Value() == nil {
		return nil, nil
	}
	article, ok := obj.Value().(*Article)
	if !ok {
		return nil, errors.Wrapf(errors.ErrState, "expected article, got %T", obj.Value())
	}
	return article.BlogID, nil
}

type CommentBucket struct {
	morm.ModelBucket
}

// NewCommentBucket returns a new comment bucket
func NewCommentBucket() *CommentBucket {
	return &CommentBucket{
		morm.NewModelBucket("comment", &Comment{},
			morm.WithIndex("article", commentArticleIDIndexer, false),
			morm.WithIndex("user", commentUserIDIndexer, false)),
	}
}

// commentArticleIDIndexer enables querying comment by article ID
func commentArticleIDIndexer(obj orm.Object) ([]byte, error) {
	if obj == nil || obj.Value() == nil {
		return nil, nil
	}
	comment, ok := obj.Value().(*Comment)
	if !ok {
		return nil, errors.Wrapf(errors.ErrState, "expected comment, got %T", obj.Value())
	}
	return comment.ArticleID, nil
}

// commentUserIDIndexer enables querying comment by user ID
func commentUserIDIndexer(obj orm.Object) ([]byte, error) {
	if obj == nil || obj.Value() == nil {
		return nil, nil
	}
	comment, ok := obj.Value().(*Comment)
	if !ok {
		return nil, errors.Wrapf(errors.ErrState, "expected comment, got %T", obj.Value())
	}
	return comment.Owner, nil
}

type LikeBucket struct {
	morm.ModelBucket
}

// NewLikeBucket returns a new like bucket
func NewLikeBucket() *LikeBucket {
	return &LikeBucket{
		morm.NewModelBucket("like", &Like{},
			morm.WithIndex("article", likeArticleIDIndexer, false)),
	}
}

// likeArticleIDIndexer enables querying comment by user ID
func likeArticleIDIndexer(obj orm.Object) ([]byte, error) {
	if obj == nil || obj.Value() == nil {
		return nil, nil
	}
	like, ok := obj.Value().(*Like)
	if !ok {
		return nil, errors.Wrapf(errors.ErrState, "expected like, got %T", obj.Value())
	}
	return like.ArticleID, nil
}
