package blog

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/migration"
)

func init() {
	migration.MustRegister(1, &CreateUserMsg{}, migration.NoModification)
	migration.MustRegister(1, &CreateBlogMsg{}, migration.NoModification)
	migration.MustRegister(1, &ChangeBlogOwnerMsg{}, migration.NoModification)
	migration.MustRegister(1, &CreateArticleMsg{}, migration.NoModification)
	migration.MustRegister(1, &DeleteArticleMsg{}, migration.NoModification)
	migration.MustRegister(1, &CreateCommentMsg{}, migration.NoModification)
	migration.MustRegister(1, &CreateLikeMsg{}, migration.NoModification)
}

var _ weave.Msg = (*CreateUserMsg)(nil)

// Path returns the routing path for this message.
func (CreateUserMsg) Path() string {
	return "blog/create_user"
}

// Validate ensures the CreateUserMsg is valid
func (m CreateUserMsg) Validate() error {
	var errs error

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
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

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())

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

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "BlogID", isGenID(m.BlogID, false))

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

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "BlogID", isGenID(m.BlogID, false))

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

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "ArticleID", isGenID(m.ArticleID, false))

	return errs
}

var _ weave.Msg = (*CreateCommentMsg)(nil)

// Path returns the routing path for this message.
func (CreateCommentMsg) Path() string {
	return "blog/create_comment"
}

// Validate ensures the CreateCommentMsg is valid
func (m CreateCommentMsg) Validate() error {
	var errs error

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "ArticleID", isGenID(m.ArticleID, false))

	if !validArticleContent(m.Content) {
		errs = errors.AppendField(errs, "Content", errors.ErrModel)
	}

	return errs
}

var _ weave.Msg = (*CreateLikeMsg)(nil)

// Path returns the routing path for this message.
func (CreateLikeMsg) Path() string {
	return "blog/create_like"
}

// Validate ensures the CreateLikeMsg is valid
func (m CreateLikeMsg) Validate() error {
	var errs error

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "ArticleID", isGenID(m.ArticleID, false))

	return errs
}
