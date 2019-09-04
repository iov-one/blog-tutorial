package custom

import (
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/migration"
	"github.com/iov-one/weave/orm"
)

func init() {
	// Migration needs to be registered for every message introduced in the codec.
	// This is the convention to message versioning.
	migration.MustRegister(1, &TimedState{}, migration.NoModification)
	migration.MustRegister(1, &State{}, migration.NoModification)
}

var _ orm.Model = (*TimedState)(nil)

// Validate ensures the TimedState fields are valid
func (m *TimedState) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	errs = errors.AppendField(errs, "Str", stringValidation(m.Str))
	if m.Byte == nil {
		errs = errors.AppendField(errs, "Byte", errors.ErrEmpty)
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

// Copy produces a new TimedState clone to fulfill the Model interface
func (m *TimedState) Copy() orm.CloneableData {
	return &TimedState{
		Metadata:       m.Metadata.Copy(),
		InnerStateEnum: m.InnerStateEnum,
		Str:            m.Str,
		Byte:           copyBytes(m.Byte),
		DeleteAt:       m.DeleteAt,
	}
}

var _ orm.Model = (*State)(nil)

// Validate ensures the State fields are valid
func (m *State) Validate() error {
	var errs error

	errs = errors.AppendField(errs, "Metadata", m.Metadata.Validate())
	if m.InnerState == nil {
		errs = errors.AppendField(errs, "InnerState", errors.ErrEmpty)
	}
	errs = errors.AppendField(errs, "Address", m.Address.Validate())
	if err := m.CreatedAt.Validate(); err != nil {
		errs = errors.AppendField(errs, "CreatedAt", m.CreatedAt.Validate())
	} else if m.CreatedAt == 0 {
		errs = errors.AppendField(errs, "CreatedAt", errors.ErrEmpty)
	}

	return errs
}

// Copy produces a new State clone to fulfill the Model interface
func (m *State) Copy() orm.CloneableData {
	return &State{
		Metadata:   m.Metadata.Copy(),
		InnerState: m.InnerState,
		Address:    copyBytes(m.Address),
		CreatedAt:  m.CreatedAt,
	}
}

func copyBytes(in []byte) []byte {
	if in == nil {
		return nil
	}
	cpy := make([]byte, len(in))
	copy(cpy, in)
	return cpy
}
