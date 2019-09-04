package custom

import (
	"strings"

	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/migration"
)

func init() {
	migration.MustRegister(1, &CreateStateMsg{}, migration.NoModification)
	migration.MustRegister(1, &CreateTimedStateMsg{}, migration.NoModification)
	migration.MustRegister(1, &DeleteTimedStateMsg{}, migration.NoModification)
}

var _ weave.Msg = (*CreateTimedStateMsg)(nil)

// Path returns the routing path for this message.
func (CreateTimedStateMsg) Path() string {
	return "custom/create_timed_state"
}

// Validate ensures the CreateTimedStateMsg is valid
func (m CreateTimedStateMsg) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "Str", stringValidation(m.Str))
	if m.Byte == nil {
		errs = errors.Append(errs, errors.Field("Byte", errors.ErrEmpty, "missing byte"))
	}
	if m.InnerStateEnum != InnerStateEnum_CaseOne && m.InnerStateEnum != InnerStateEnum_CaseTwo {
		errs = errors.AppendField(errs, "InnerStateEnum", errors.ErrState)
	}

	if m.DeleteAt == 0 {
		return errs
	} else if err := m.DeleteAt.Validate(); err != nil {
		errs = errors.AppendField(errs, "DeleteAt", m.DeleteAt.Validate())
	}
	return errs
}

var _ weave.Msg = (*DeleteTimedStateMsg)(nil)

// Path returns the routing path for this message.
func (DeleteTimedStateMsg) Path() string {
	return "custom/delete_timed_state"
}

// Validate ensures the DeleteTimedStateMsg is valid
func (m DeleteTimedStateMsg) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "TimedStateID", validID(m.TimedStateID))
	return errs
}

var _ weave.Msg = (*CreateStateMsg)(nil)

// Path returns the routing path for this message.
func (CreateStateMsg) Path() string {
	return "custom/create_state"
}

// Validate ensures the CreateStateMsg is valid
func (m CreateStateMsg) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "Address", m.Address.Validate())
	if m.InnerState == nil {
		errs = errors.AppendField(errs, "InnerState", errors.ErrEmpty)
	}
	return errs
}

// validID returns an error if this is not an 8-byte ID
// as expected for orm.IDGenBucket
func validID(id []byte) error {
	if len(id) == 0 {
		return errors.Wrap(errors.ErrEmpty, "id missing")
	}
	if len(id) != 8 {
		return errors.Wrap(errors.ErrInput, "id is invalid length (expect 8 bytes)")
	}
	return nil
}

func stringValidation(str string) error {
	if len(str) == 0 {
		return errors.Wrap(errors.ErrEmpty, "string missing")
	}
	if !strings.HasPrefix(str, "cstm") {
		return errors.Wrap(errors.ErrInput, "string does not have cstm prefix")
	}
	return nil
}
