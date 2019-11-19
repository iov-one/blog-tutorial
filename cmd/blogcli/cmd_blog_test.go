package main

import (
	"bytes"
	"github.com/iov-one/weave/weavetest"
	"testing"
	"time"

	"github.com/iov-one/blog-tutorial/x/blog"
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/weavetest/assert"
)

func TestCreateBlogUser(t *testing.T) {
	var output bytes.Buffer
	args := []string{
		"-username", "test-username",
		"-bio", "test bio",
	}
	if err := cmdCreateUser(nil, &output, args); err != nil {
		t.Fatalf("cannot create a new user transaction: %s", err)
	}

	tx, _, err := readTx(&output)
	if err != nil {
		t.Fatalf("cannot unmarshal created transaction: %s", err)
	}

	txmsg, err := tx.GetMsg()
	if err != nil {
		t.Fatalf("cannot get transaction message: %s", err)
	}
	msg := txmsg.(*blog.CreateUserMsg)

	assert.Equal(t, "test-username", msg.Username)
	assert.Equal(t, "test bio", msg.Bio)
}

func TestCreateBlog(t *testing.T) {
	var output bytes.Buffer
	args := []string{
		"-title", "test title",
		"-desc", "test desc",
	}
	if err := cmdCreateBlog(nil, &output, args); err != nil {
		t.Fatalf("cannot create a new blog transaction: %s", err)
	}

	tx, _, err := readTx(&output)
	if err != nil {
		t.Fatalf("cannot unmarshal created transaction: %s", err)
	}

	txmsg, err := tx.GetMsg()
	if err != nil {
		t.Fatalf("cannot get transaction message: %s", err)
	}
	msg := txmsg.(*blog.CreateBlogMsg)

	assert.Equal(t, "test title", msg.Title)
	assert.Equal(t, "test desc", msg.Description)
}

func TestChangeBlogOwner(t *testing.T) {
	var output bytes.Buffer
	args := []string{
		"-blog_key", "122333",
		"-new_owner", "E28AE9A6EB94FC88B73EB7CBD6B87BF93EB9BEF0",
	}
	if err := cmdChangeBlogOwner(nil, &output, args); err != nil {
		t.Fatalf("cannot create a new change blog owner transaction: %s", err)
	}

	tx, _, err := readTx(&output)
	if err != nil {
		t.Fatalf("cannot unmarshal created transaction: %s", err)
	}

	txmsg, err := tx.GetMsg()
	if err != nil {
		t.Fatalf("cannot get transaction message: %s", err)
	}
	msg := txmsg.(*blog.ChangeBlogOwnerMsg)

	assert.Equal(t, weavetest.SequenceID(122333), []byte(msg.BlogKey))
	assert.Equal(t, fromHex(t, "E28AE9A6EB94FC88B73EB7CBD6B87BF93EB9BEF0"), []byte(msg.NewOwner))
}

func TestCreateArticle(t *testing.T) {
	var output bytes.Buffer
	currentTime := time.Now().UTC()
	currentTimeStr := currentTime.Format(flagTimeFormat)
	args := []string{
		"-blog_key", "122333",
		"-title", "test title",
		"-content", "test content",
		"-delete_at", currentTimeStr,
	}
	if err := cmdCreateArticle(nil, &output, args); err != nil {
		t.Fatalf("cannot create a new article transaction: %s", err)
	}

	tx, _, err := readTx(&output)
	if err != nil {
		t.Fatalf("cannot unmarshal created transaction: %s", err)
	}

	txmsg, err := tx.GetMsg()
	if err != nil {
		t.Fatalf("cannot get transaction message: %s", err)
	}
	msg := txmsg.(*blog.CreateArticleMsg)

	assert.Equal(t, weavetest.SequenceID(122333), msg.BlogKey)
	assert.Equal(t, "test title", msg.Title)
	assert.Equal(t, "test content", msg.Content)

	testT, _ := time.Parse(flagTimeFormat, currentTimeStr)
	assert.Equal(t, weave.AsUnixTime(testT), msg.DeleteAt)
}

func TestDeleteArticle(t *testing.T) {
	var output bytes.Buffer
	args := []string{
		"-article_key", "122333",
	}
	if err := cmdDeleteArticle(nil, &output, args); err != nil {
		t.Fatalf("cannot create a delete article: %s", err)
	}

	tx, _, err := readTx(&output)
	if err != nil {
		t.Fatalf("cannot unmarshal created transaction: %s", err)
	}

	txmsg, err := tx.GetMsg()
	if err != nil {
		t.Fatalf("cannot get transaction message: %s", err)
	}
	msg := txmsg.(*blog.DeleteArticleMsg)

	assert.Equal(t, weavetest.SequenceID(122333), msg.ArticleKey)
}

func TestCancelDeleteArticleTask(t *testing.T) {
	var output bytes.Buffer
	args := []string{
		"-task_id", "122333",
	}
	if err := cmdCancelDeleteArticleTask(nil, &output, args); err != nil {
		t.Fatalf("cannot create a new cancel delete article task transaction: %s", err)
	}

	tx, _, err := readTx(&output)
	if err != nil {
		t.Fatalf("cannot unmarshal created transaction: %s", err)
	}

	txmsg, err := tx.GetMsg()
	if err != nil {
		t.Fatalf("cannot get transaction message: %s", err)
	}
	msg := txmsg.(*blog.CancelDeleteArticleTaskMsg)

	assert.Equal(t, weavetest.SequenceID(122333), msg.TaskID)
}
