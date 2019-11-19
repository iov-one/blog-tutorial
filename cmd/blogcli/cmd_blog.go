package main

import (
	"flag"
	"fmt"
	"io"

	app "github.com/iov-one/blog-tutorial/cmd/blog/app"
	"github.com/iov-one/blog-tutorial/x/blog"
	"github.com/iov-one/weave"
)

func cmdCreateUser(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), `
Create a blog user.
		`)
		fl.PrintDefaults()
	}
	var (
		usernameFl = fl.String("username", "", "The username. For example 'alice'")
		bioFl      = fl.String("bio", "", "Bio of the user")
	)
	fl.Parse(args)

	msg := blog.CreateUserMsg{
		Metadata: &weave.Metadata{Schema: 1},
		Username: *usernameFl,
		Bio:      *bioFl,
	}
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("given data produce an invalid message: %s", err)
	}

	tx := &app.Tx{
		Sum: &app.Tx_BlogCreateUserMsg{
			BlogCreateUserMsg: &msg,
		},
	}
	_, err := writeTx(output, tx)
	return err
}

func cmdCreateBlog(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), `
Create a blog.
		`)
		fl.PrintDefaults()
	}
	var (
		titleFl = fl.String("title", "", "Title of the blog")
		descFl  = fl.String("desc", "", "Description of the blog")
	)
	fl.Parse(args)

	msg := blog.CreateBlogMsg{
		Metadata:    &weave.Metadata{Schema: 1},
		Title:       *titleFl,
		Description: *descFl,
	}
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("given data produce an invalid message: %s", err)
	}

	tx := &app.Tx{
		Sum: &app.Tx_BlogCreateBlogMsg{
			BlogCreateBlogMsg: &msg,
		},
	}
	_, err := writeTx(output, tx)
	return err
}

func cmdChangeBlogOwner(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), `
Change owner of the blog.
		`)
		fl.PrintDefaults()
	}
	var (
		blogKeyFl  = flSeq(fl, "blog_key", "", "Identifier of the blog")
		newOwnerFl = flAddress(fl, "new_owner", "", "Address of the new owner")
	)
	fl.Parse(args)

	msg := blog.ChangeBlogOwnerMsg{
		Metadata: &weave.Metadata{Schema: 1},
		BlogKey:  *blogKeyFl,
		NewOwner: *newOwnerFl,
	}
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("given data produce an invalid message: %s", err)
	}

	tx := &app.Tx{
		Sum: &app.Tx_BlogChangeBlogOwnerMsg{
			BlogChangeBlogOwnerMsg: &msg,
		},
	}
	_, err := writeTx(output, tx)
	return err
}

func cmdCreateArticle(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), `
Post an article under a blog.
		`)
		fl.PrintDefaults()
	}
	var (
		blogKeyFl  = flSeq(fl, "blog_key", "", "Identifier of the blog that article will be posted at")
		titleFl    = fl.String("title", "", "Title of the article")
		contentFl  = fl.String("content", "", "Content of the article")
		deleteAtFl = flTime(fl, "delete_at", nil, "Deletion time of the article, format: 2006-01-02 15:04")
	)
	fl.Parse(args)

	msg := blog.CreateArticleMsg{
		Metadata: &weave.Metadata{Schema: 1},
		BlogKey:  *blogKeyFl,
		Title:    *titleFl,
		Content:  *contentFl,
		DeleteAt: deleteAtFl.UnixTime(),
	}
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("given data produce an invalid message: %s", err)
	}

	tx := &app.Tx{
		Sum: &app.Tx_BlogCreateArticleMsg{
			BlogCreateArticleMsg: &msg,
		},
	}
	_, err := writeTx(output, tx)
	return err
}

func cmdDeleteArticle(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), `
Delete an article.
		`)
		fl.PrintDefaults()
	}
	var (
		articleKeyFl = flSeq(fl, "article_key", "", "Identifer of the article")
	)
	fl.Parse(args)

	msg := blog.DeleteArticleMsg{
		Metadata:   &weave.Metadata{Schema: 1},
		ArticleKey: *articleKeyFl,
	}
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("given data produce an invalid message: %s", err)
	}

	tx := &app.Tx{
		Sum: &app.Tx_BlogDeleteArticleMsg{
			BlogDeleteArticleMsg: &msg,
		},
	}
	_, err := writeTx(output, tx)
	return err
}

func cmdCancelDeleteArticleTask(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), `
Cancel a delete article task.
		`)
		fl.PrintDefaults()
	}
	var (
		taskIDFl = flSeq(fl, "task_id", "", "Identifer of the delete task")
	)
	fl.Parse(args)

	msg := blog.CancelDeleteArticleTaskMsg{
		Metadata: &weave.Metadata{Schema: 1},
		TaskID:   *taskIDFl,
	}
	if err := msg.Validate(); err != nil {
		return fmt.Errorf("given data produce an invalid message: %s", err)
	}

	tx := &app.Tx{
		Sum: &app.Tx_BlogCancelDeleteArticleTaskMsg{
			BlogCancelDeleteArticleTaskMsg: &msg,
		},
	}
	_, err := writeTx(output, tx)
	return err
}
