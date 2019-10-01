package blog

import (
	"bytes"
	"encoding/binary"

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
			morm.WithIndex("blog", articleBlogIDIndexer, false),
			morm.WithIndex("timedBlog", blogTimedIndexer, false)),
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

// blogTimedIndexer indexes articles by
//   (blog id, createdAt)
// so give us easy lookup of the most recently posted articles on a given blog
// (we can also use this client side with range queries to select all trades on a given
// blog during any given timeframe)
func blogTimedIndexer(obj orm.Object) ([]byte, error) {
	if obj == nil || obj.Value() == nil {
		return nil, nil
	}
	article, ok := obj.Value().(*Article)
	if !ok {
		return nil, errors.Wrapf(errors.ErrState, "expected article, got %T", obj.Value())
	}

	return BuildBlogTimedIndex(article)
}

// BuildBlogTimedIndex produces 8 bytes BlogID || big-endian createdAt
// This allows lexographical searches over the time ranges (or earliest or latest)
// of all articles within one blog
func BuildBlogTimedIndex(article *Article) ([]byte, error) {
	res := make([]byte, 16)
	copy(res, article.BlogID)
	// this would violate lexographical ordering as negatives would be highest
	if article.CreatedAt < 0 {
		return nil, errors.Wrap(errors.ErrState, "cannot index negative creation times")
	}
	binary.BigEndian.PutUint64(res[8:], uint64(article.CreatedAt))
	return res, nil
}

type DeleteArticleTaskBucket struct {
	morm.ModelBucket
}

// NewDeleteArticleTaskBucket returns a new delete article task bucket
func NewDeleteArticleTaskBucket() *DeleteArticleTaskBucket {
	return &DeleteArticleTaskBucket{
		morm.NewModelBucket("deleteart", &DeleteArticleTask{}),
	}
}

type CommentBucket struct {
	morm.ModelBucket
}

// NewCommentBucket returns a new comment bucket
func NewCommentBucket() *CommentBucket {
	return &CommentBucket{
		morm.NewModelBucket("comment", &Comment{},
			morm.WithIndex("article", commentArticleIDIndexer, false),
			morm.WithIndex("user", commentUserIDIndexer, false),
			morm.WithIndex("articleuser", articleUserIndexer, false)),
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

// articleUserIndexer produces in SQL parlance, a compound index.
// Used for returning users all comments on an article
// (articleID, userID) -> index
func articleUserIndexer(obj orm.Object) ([]byte, error) {
	if obj == nil || obj.Value() == nil {
		return nil, nil
	}
	comment, ok := obj.Value().(*Comment)

	if !ok {
		return nil, errors.Wrapf(errors.ErrState, "expected comment, got %T", obj.Value())
	}

	return BuildArticleUserIndex(comment), nil
}

// BuildArticleUserIndex indexByteSize = 8(ArticleID) + 8(UserID)
func BuildArticleUserIndex(comment *Comment) []byte {
	return bytes.Join([][]byte{comment.ArticleID, comment.Owner}, nil)
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
