package blog

import (
	"regexp"
	"time"

	"github.com/iov-one/blog-tutorial/morm"

	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/orm"
)

var _ morm.Model = (*User)(nil)

func (u *User) IsRegisteredAfterDate(date time.Time) bool {
	return u.RegisteredAt.Time().After(date)
}

// SetID is a minimal implementation, useful when the ID is a separate protobuf field
func (m *User) SetID(id []byte) error {
	m.ID = id
	return nil
}

// Copy produces a new copy to fulfill the Model interface
func (m *User) Copy() orm.CloneableData {
	return &User{
		Metadata:     m.Metadata.Copy(),
		ID:           copyBytes(m.ID),
		Username:     m.Username,
		Bio:          m.Bio,
		RegisteredAt: m.RegisteredAt,
	}
}

var validUsername = regexp.MustCompile(`^[a-zA-Z0-9_.-]{4,16}$`).MatchString
var validBio = regexp.MustCompile(`^[a-zA-Z0-9_ ]{4,200}$`).MatchString

// Validate validates user's fields
func (m *User) Validate() error {
	var errs error

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "ID", isGenID(m.ID, false))

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

var _ morm.Model = (*Blog)(nil)

// SetID is a minimal implementation, useful when the ID is a separate protobuf field
func (m *Blog) SetID(id []byte) error {
	m.ID = id
	return nil
}

// Copy produces a new copy to fulfill the Model interface
func (m *Blog) Copy() orm.CloneableData {
	return &Blog{
		Metadata:    m.Metadata.Copy(),
		ID:          copyBytes(m.ID),
		Owner:       m.Owner.Clone(),
		Title:       m.Title,
		Description: m.Description,
		CreatedAt:   m.CreatedAt,
	}
}

var validBlogTitle = regexp.MustCompile(`^[a-zA-Z0-9$@$!%*?&#'^;-_. +]{4,32}$`).MatchString
var validBlogDescription = regexp.MustCompile(`^[a-zA-Z0-9$@$!%*?&#'^;-_. +]{4,1000}$`).MatchString

// Validate validates blog's fields
func (m *Blog) Validate() error {
	var errs error

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "ID", isGenID(m.ID, false))
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

var _ morm.Model = (*Blog)(nil)

// SetID is a minimal implementation, useful when the ID is a separate protobuf field
func (m *Article) SetID(id []byte) error {
	m.ID = id
	return nil
}

// Copy produces a new copy to fulfill the Model interface
// TODO remove after weave 0.22.0 is released
func (m *Article) Copy() orm.CloneableData {
	return &Article{
		Metadata:     m.Metadata.Copy(),
		ID:           copyBytes(m.ID),
		BlogID:       copyBytes(m.BlogID),
		Owner:        m.Owner.Clone(),
		Title:        m.Title,
		Content:      m.Content,
		CreatedAt:    m.CreatedAt,
		DeleteAt:     m.DeleteAt,
	}
}

var validArticleTitle = regexp.MustCompile(`^[a-zA-Z0-9_ ]{4,32}$`).MatchString
var validArticleContent = regexp.MustCompile(`^[a-zA-Z0-9_ ]{4,1000}$`).MatchString

// Validate validates article's fields
func (m *Article) Validate() error {
	var errs error

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "ID", isGenID(m.ID, false))
	errs = errors.AppendField(errs, "BlogID", isGenID(m.BlogID, false))
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

var _ morm.Model = (*DeleteArticleTask)(nil)

// SetID is a minimal implementation, useful when the ID is a separate protobuf field
func (m *DeleteArticleTask) SetID(id []byte) error {
	m.ID = id
	return nil
}

// Copy produces a new copy to fulfill the Model interface
func (m *DeleteArticleTask) Copy() orm.CloneableData {
	return &DeleteArticleTask{
		Metadata:  m.Metadata.Copy(),
		ID:        copyBytes(m.ID),
		ArticleID: copyBytes(m.ArticleID),
		TaskOwner: m.TaskOwner.Clone(),
	}
}

// Validate validates user's fields
func (m *DeleteArticleTask) Validate() error {
	var errs error

	//errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "ID", isGenID(m.ID, false))
	errs = errors.AppendField(errs, "ArticleID", isGenID(m.ArticleID, false))
	errs = errors.AppendField(errs, "TaskOwner", m.TaskOwner.Validate())

	return errs
}

func copyBytes(in []byte) []byte {
	if in == nil {
		return nil
	}
	cpy := make([]byte, len(in))
	copy(cpy, in)
	return cpy
}

// isGenID ensures that the ID is 8 byte input.
// if allowEmpty is set, we also allow empty
// TODO change with validateSequence when weave 0.22.0 is released
func isGenID(id []byte, allowEmpty bool) error {
	if len(id) == 0 {
		if allowEmpty {
			return nil
		}
		return errors.Wrap(errors.ErrEmpty, "missing id")
	}
	if len(id) != 8 {
		return errors.Wrap(errors.ErrInput, "id must be 8 bytes")
	}
	return nil
}
