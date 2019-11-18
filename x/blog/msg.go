package blog

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/migration"
	"github.com/iov-one/weave/orm"
)

func init() {
	migration.MustRegister(1, &CreateUserMsg{}, migration.NoModification)
	migration.MustRegister(1, &CreateBlogMsg{}, migration.NoModification)
	migration.MustRegister(1, &ChangeBlogOwnerMsg{}, migration.NoModification)
	migration.MustRegister(1, &CreateArticleMsg{}, migration.NoModification)
	migration.MustRegister(1, &DeleteArticleMsg{}, migration.NoModification)
}

var _ weave.Msg = (*CreateUserMsg)(nil)

// Path returns the routing path for this message.
func (CreateUserMsg) Path() string {
	return "blog/create_user"
}

// Validate ensures the CreateUserMsg is valid
func (m CreateUserMsg) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	if !validUsername(m.Username) {
		errs = errors.AppendField(errs, "Username", errors.ErrModel)
	}

	if m.Bio != "" && !validBio(m.Bio) {
		errs = errors.AppendField(errs, "Bio", errors.ErrModel)
	}

	return errs
}

var _ weave.Msg = (*CreateBlogMsg)(nil)

// Path returns the routing path for this message.
func (CreateBlogMsg) Path() string {
	return "blog/create_blog"
}

// Validate ensures the CreateBlogMsg is valid
func (m CreateBlogMsg) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	if !validBlogTitle(m.Title) {
		errs = errors.AppendField(errs, "Title", errors.ErrModel)
	}
	if !validBlogDescription(m.Description) {
		errs = errors.AppendField(errs, "Description", errors.ErrModel)
	}

	return errs
}

var _ weave.Msg = (*ChangeBlogOwnerMsg)(nil)

// Path returns the routing path for this message.
func (ChangeBlogOwnerMsg) Path() string {
	return "blog/change_blog_owner"
}

// Validate ensures the ChangeBlogOwner is valid
func (m ChangeBlogOwnerMsg) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "BlogKey", orm.ValidateSequence(m.BlogKey))

	if err := m.NewOwner.Validate(); err != nil {
		errs = errors.AppendField(errs, "NewOwner", errors.ErrInput)
	}

	return errs
}

var _ weave.Msg = (*CreateArticleMsg)(nil)

// Path returns the routing path for this message.
func (CreateArticleMsg) Path() string {
	return "blog/create_article"
}

// Validate ensures the CreateArticleMsg is valid
func (m CreateArticleMsg) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "BlogKey", orm.ValidateSequence(m.BlogKey))

	if !validBlogTitle(m.Title) {
		errs = errors.AppendField(errs, "Title", errors.ErrModel)
	}
	if !validBlogDescription(m.Content) {
		errs = errors.AppendField(errs, "Content", errors.ErrModel)
	}

	if m.DeleteAt != 0 {
		if err := m.DeleteAt.Validate(); err != nil {
			errs = errors.AppendField(errs, "DeleteAt", err)
		}
	}

	return errs
}

var _ weave.Msg = (*DeleteArticleMsg)(nil)

// Path returns the routing path for this message.
func (DeleteArticleMsg) Path() string {
	return "blog/delete_article"
}

// Validate ensures the DeleteArticle is valid
func (m DeleteArticleMsg) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "ArticleKey", orm.ValidateSequence(m.ArticleKey))

	return errs
}

var _ weave.Msg = (*CancelDeleteArticleTaskMsg)(nil)

// Path returns the routing path for this message.
func (CancelDeleteArticleTaskMsg) Path() string {
	return "blog/cancel_delete_article_task"
}

// Validate ensures the CancelDeleteArticleTaskMsg is valid
func (m CancelDeleteArticleTaskMsg) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "TaskID", orm.ValidateSequence(m.TaskID))

	return errs
}
