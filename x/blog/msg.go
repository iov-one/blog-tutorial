package blog

import (
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/migration"
)

func init() {
	migration.MustRegister(1, &CreateUserMsg{}, migration.NoModification)
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

	if !validBio(m.Bio) {
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
