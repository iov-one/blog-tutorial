package main

import (
	"bytes"
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

	assert.Equal(t, "test-username", msg.GetUsername())
	assert.Equal(t, "test bio", msg.GetBio())
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

	assert.Equal(t, "test title", msg.GetTitle())
	assert.Equal(t, "test desc", msg.GetDescription())
}

func TestChangeBlogOwner(t *testing.T) {
	var output bytes.Buffer
	args := []string{
		"-blog_id", "1",
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

	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 1}, []byte(msg.GetBlogID()))
	assert.Equal(t, fromHex(t, "E28AE9A6EB94FC88B73EB7CBD6B87BF93EB9BEF0"), []byte(msg.GetNewOwner()))
}

func TestCreateArticle(t *testing.T) {
	var output bytes.Buffer
	currentTime := time.Now().UTC()
	currentTimeStr := currentTime.Format(flagTimeFormat)
	args := []string{
		"-blog_id", "1",
		"-title", "test title",
		"-content", "test content",
		"-delete_at", currentTimeStr,
	}
	if err := cmdCreateArticle(nil, &output, args); err != nil {
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
	msg := txmsg.(*blog.CreateArticleMsg)

	assert.Equal(t, []byte{0, 0, 0, 0, 0, 0, 0, 1}, []byte(msg.GetBlogID()))
	assert.Equal(t, "test title", msg.GetTitle())
	assert.Equal(t, "test content", msg.GetContent())
	assert.Equal(t, weave.AsUnixTime(currentTime), msg.GetDeleteAt())
}
