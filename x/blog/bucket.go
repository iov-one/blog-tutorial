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
		morm.NewModelBucket("blog", &Blog{}),
	}
}

type ArticleBucket struct {
	morm.ModelBucket
}

// NewArticleBucket returns a new article bucket
func NewArticleBucket() *ArticleBucket {
	return &ArticleBucket{
		morm.NewModelBucket("article", &Article{},
			morm.WithIndex("blog", blogIDIndexer, true)),
	}
}

// blogIDIndexer enables querying articles by blog ids
func blogIDIndexer(obj orm.Object) ([]byte, error) {
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
		morm.NewModelBucket("comment", &Comment{}),
	}
}

type LikeBucket struct {
	morm.ModelBucket
}

// NewLikeBucket returns a new like bucket
func NewLikeBucket() *LikeBucket {
	return &LikeBucket{
		morm.NewModelBucket("like", &Like{}),
	}
}
