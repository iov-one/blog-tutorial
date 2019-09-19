package blog

import (
	"github.com/iov-one/blog-tutorial/morm"
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
		morm.NewModelBucket("article", &Article{}),
	}
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
