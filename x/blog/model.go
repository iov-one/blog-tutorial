package blog

import (
	"regexp"
	"time"

	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/orm"
)

var _ orm.SerialModel = (*User)(nil)

func (u *User) IsRegisteredAfterDate(date time.Time) bool {
	return u.RegisteredAt.Time().After(date)
}

// SetPrimaryKey is a minimal implementation, useful when the PrimaryKey is a separate protobuf field
func (m *User) SetPrimaryKey(pk []byte) error {
	m.PrimaryKey = pk
	return nil
}

var validUsername = regexp.MustCompile(`^[a-zA-Z0-9_.-]{4,16}$`).MatchString
var validBio = regexp.MustCompile(`^[a-zA-Z0-9_ ]{4,200}$`).MatchString

// Validate validates user's fields
func (m *User) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "PrimaryKey", orm.ValidateSequence(m.PrimaryKey))

	if !validUsername(m.Username) {
		errs = errors.AppendField(errs, "Username", errors.ErrModel)
	}

	if m.Bio != "" && !validBio(m.Bio) {
		errs = errors.AppendField(errs, "Bio", errors.ErrModel)
	}

	if err := m.RegisteredAt.Validate(); err != nil {
		errs = errors.AppendField(errs, "RegisteredAt", m.RegisteredAt.Validate())
	} else if m.RegisteredAt == 0 {
		errs = errors.AppendField(errs, "RegisteredAt", errors.ErrEmpty)
	}

	return errs
}

var _ orm.SerialModel = (*Blog)(nil)

// SetPrimaryKey is a minimal implementation, useful when the PrimaryKey is a separate protobuf field
func (m *Blog) SetPrimaryKey(pk []byte) error {
	m.PrimaryKey = pk
	return nil
}

var validBlogTitle = regexp.MustCompile(`^[a-zA-Z0-9$@$!%*?&#'^;-_. +]{4,32}$`).MatchString
var validBlogDescription = regexp.MustCompile(`^[a-zA-Z0-9$@$!%*?&#'^;-_. +]{4,1000}$`).MatchString

// Validate validates blog's fields
func (m *Blog) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "PrimaryKey", orm.ValidateSequence(m.PrimaryKey))
	errs = errors.AppendField(errs, "Owner", m.Owner.Validate())

	if !validBlogTitle(m.Title) {
		errs = errors.AppendField(errs, "Title", errors.ErrModel)
	}
	if !validBlogDescription(m.Description) {
		errs = errors.AppendField(errs, "Description", errors.ErrModel)
	}

	if err := m.CreatedAt.Validate(); err != nil {
		errs = errors.AppendField(errs, "CreatedAt", err)
	} else if m.CreatedAt == 0 {
		errs = errors.AppendField(errs, "CreatedAt", errors.ErrEmpty)
	}

	return errs
}

var _ orm.SerialModel = (*Blog)(nil)

// SetPrimaryKey is a minimal implementation, useful when the PrimaryKey is a separate protobuf field
func (m *Article) SetPrimaryKey(pk []byte) error {
	m.PrimaryKey = pk
	return nil
}

var validArticleTitle = regexp.MustCompile(`^[a-zA-Z0-9_ ]{4,32}$`).MatchString
var validArticleContent = regexp.MustCompile(`^[a-zA-Z0-9_ ]{4,1000}$`).MatchString

// Validate validates article's fields
func (m *Article) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "PrimaryKey", orm.ValidateSequence(m.PrimaryKey))
	errs = errors.AppendField(errs, "BlogKey", orm.ValidateSequence(m.BlogKey))
	errs = errors.AppendField(errs, "Owner", m.Owner.Validate())

	if !validBlogTitle(m.Title) {
		errs = errors.AppendField(errs, "Title", errors.ErrModel)
	}
	if !validBlogDescription(m.Content) {
		errs = errors.AppendField(errs, "Content", errors.ErrModel)
	}

	if err := m.CreatedAt.Validate(); err != nil {
		errs = errors.AppendField(errs, "CreatedAt", err)
	} else if m.CreatedAt == 0 {
		errs = errors.AppendField(errs, "CreatedAt", errors.ErrEmpty)
	}

	if m.DeleteAt != 0 {
		if err := m.DeleteAt.Validate(); err != nil {
			errs = errors.AppendField(errs, "DeleteAt", err)
		}
	}

	return errs
}

var _ orm.SerialModel = (*DeleteArticleTask)(nil)

// SetPrimaryKey is a minimal implementation, useful when the PrimaryKey is a separate protobuf field
func (m *DeleteArticleTask) SetPrimaryKey(pk []byte) error {
	m.PrimaryKey = pk
	return nil
}

// Validate validates user's fields
func (m *DeleteArticleTask) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "PrimaryKey", orm.ValidateSequence(m.PrimaryKey))
	errs = errors.AppendField(errs, "ArticleKey", orm.ValidateSequence(m.ArticleKey))
	errs = errors.AppendField(errs, "TaskOwner", m.TaskOwner.Validate())

	return errs
}
