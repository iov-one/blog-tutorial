package blog

import (
	"encoding/binary"

	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/orm"
)

type UserBucket struct {
	orm.SerialModelBucket
}

// NewUserBucket returns a new user bucket
func NewUserBucket() *UserBucket {
	return &UserBucket{
		orm.NewSerialModelBucket("user", &User{}),
	}
}

type BlogBucket struct {
	orm.SerialModelBucket
}

// NewBlogBucket returns a new blog bucket
func NewBlogBucket() *BlogBucket {
	return &BlogBucket{
		orm.NewSerialModelBucket("blog", &Blog{},
			orm.WithIndexSerial("user", blogUserIDIndexer, false)),
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
	orm.SerialModelBucket
}

// NewArticleBucket returns a new article bucket
func NewArticleBucket() *ArticleBucket {
	return &ArticleBucket{
		orm.NewSerialModelBucket("article", &Article{},
			orm.WithIndexSerial("blog", articleBlogIDIndexer, false),
			orm.WithIndexSerial("timedBlog", blogTimedIndexer, false)),
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
	return article.BlogKey, nil

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

// BuildBlogTimedIndex produces 8 bytes BlogKey || big-endian createdAt
// This allows lexographical searches over the time ranges (or earliest or latest)
// of all articles within one blog
func BuildBlogTimedIndex(article *Article) ([]byte, error) {
	res := make([]byte, 16)
	copy(res, article.BlogKey)
	// this would violate lexographical ordering as negatives would be highest
	if article.CreatedAt < 0 {
		return nil, errors.Wrap(errors.ErrState, "cannot index negative creation times")
	}
	binary.BigEndian.PutUint64(res[8:], uint64(article.CreatedAt))
	return res, nil
}

type DeleteArticleTaskBucket struct {
	orm.SerialModelBucket
}

// NewDeleteArticleTaskBucket returns a new delete article task bucket
func NewDeleteArticleTaskBucket() *DeleteArticleTaskBucket {
	return &DeleteArticleTaskBucket{
		orm.NewSerialModelBucket("deleteart", &DeleteArticleTask{}),
	}
}
