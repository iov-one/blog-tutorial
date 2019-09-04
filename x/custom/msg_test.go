package custom

import (
	"testing"
	"time"

	"github.com/iov-one/weave"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/weavetest"
	"github.com/iov-one/weave/weavetest/assert"
)

func TestValidateCreateTimedStateMsg(t *testing.T) {
	now := weave.AsUnixTime(time.Now())
	future := now.Add(time.Hour)

	cases := map[string]struct {
		msg      weave.Msg
		wantErrs map[string]*errors.Error
	}{
		"success": {
			msg: &CreateTimedStateMsg{
				Metadata:       &weave.Metadata{Schema: 1},
				InnerStateEnum: InnerStateEnum_CaseOne,
				Str:            "cstm:str",
				Byte:           []byte{0, 1},
				DeleteAt:       future,
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            nil,
				"Byte":           nil,
				"DeleteAt":       nil,
			},
		},
		"success, no delete at": {
			msg: &CreateTimedStateMsg{
				Metadata:       &weave.Metadata{Schema: 1},
				InnerStateEnum: InnerStateEnum_CaseOne,
				Str:            "cstm:str",
				Byte:           []byte{0, 1},
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            nil,
				"Byte":           nil,
				"DeleteAt":       nil,
			},
		},
		"missing metadata": {
			msg: &CreateTimedStateMsg{
				InnerStateEnum: InnerStateEnum_CaseTwo,
				Str:            "cstm:str",
				Byte:           []byte{0, 1},
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":       errors.ErrMetadata,
				"InnerStateEnum": nil,
				"Str":            nil,
				"Byte":           nil,
				"DeleteAt":       nil,
			},
		},
		"missing inner state enum": {
			msg: &CreateTimedStateMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Str:      "cstm:str",
				Byte:     []byte{0, 1},
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": errors.ErrState,
				"Str":            nil,
				"Byte":           nil,
				"DeleteAt":       nil,
			},
		},
		"missing str": {
			msg: &CreateTimedStateMsg{
				Metadata:       &weave.Metadata{Schema: 1},
				InnerStateEnum: InnerStateEnum_CaseTwo,
				Byte:           []byte{0, 1},
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            errors.ErrEmpty,
				"Byte":           nil,
				"DeleteAt":       nil,
			},
		},
		"str does not have 'cstm' prefix": {
			msg: &CreateTimedStateMsg{
				Metadata:       &weave.Metadata{Schema: 1},
				InnerStateEnum: InnerStateEnum_CaseTwo,
				Str:            "str",
				Byte:           []byte{0, 1},
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            errors.ErrInput,
				"Byte":           nil,
				"DeleteAt":       nil,
			},
		},
		"missing byte": {
			msg: &CreateTimedStateMsg{
				Metadata:       &weave.Metadata{Schema: 1},
				InnerStateEnum: InnerStateEnum_CaseOne,
				Str:            "cstm:str",
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":       nil,
				"InnerStateEnum": nil,
				"Str":            nil,
				"Byte":           errors.ErrEmpty,
				"DeleteAt":       nil,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			err := tc.msg.Validate()
			for field, wantErr := range tc.wantErrs {
				assert.FieldError(t, err, field, wantErr)
			}
		})
	}
}

func TestValidateDeleteStateMsg(t *testing.T) {
	cases := map[string]struct {
		msg      weave.Msg
		wantErrs map[string]*errors.Error
	}{
		"success": {
			msg: &DeleteTimedStateMsg{
				Metadata:     &weave.Metadata{Schema: 1},
				TimedStateID: weavetest.SequenceID(1),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"TimedStateID": nil,
			},
		},
		"failure missing metadata": {
			msg: &DeleteTimedStateMsg{
				TimedStateID: weavetest.SequenceID(1),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     errors.ErrMetadata,
				"TimedStateID": nil,
			},
		},
		"failure invalid id": {
			msg: &DeleteTimedStateMsg{
				Metadata:     &weave.Metadata{Schema: 1},
				TimedStateID: []byte{7, 99, 0},
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"TimedStateID": errors.ErrInput,
			},
		},
		"failure missing id": {
			msg: &DeleteTimedStateMsg{
				Metadata: &weave.Metadata{Schema: 1},
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":     nil,
				"TimedStateID": errors.ErrEmpty,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			err := tc.msg.Validate()
			for field, wantErr := range tc.wantErrs {
				assert.FieldError(t, err, field, wantErr)
			}
		})
	}
}

func TestValidateCreateStateMsg(t *testing.T) {
	cases := map[string]struct {
		msg      weave.Msg
		wantErrs map[string]*errors.Error
	}{
		"success": {
			msg: &CreateStateMsg{
				Metadata:   &weave.Metadata{Schema: 1},
				InnerState: &InnerState{St1: 1, St2: 2},
				Address:    weavetest.NewCondition().Address(),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"InnerState": nil,
				"Address":    nil,
			},
		},
		"missing metadata": {
			msg: &CreateStateMsg{
				InnerState: &InnerState{St1: 1, St2: 2},
				Address:    weavetest.NewCondition().Address(),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   errors.ErrMetadata,
				"InnerState": nil,
				"Address":    nil,
			},
		},
		"missing inner state": {
			msg: &CreateStateMsg{
				Metadata: &weave.Metadata{Schema: 1},
				Address:  weavetest.NewCondition().Address(),
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"InnerState": errors.ErrEmpty,
				"Address":    nil,
			},
		},
		"bad address": {
			msg: &CreateStateMsg{
				Metadata:   &weave.Metadata{Schema: 1},
				InnerState: &InnerState{St1: 1, St2: 2},
				Address:    []byte{0, 1},
			},
			wantErrs: map[string]*errors.Error{
				"Metadata":   nil,
				"InnerState": nil,
				"Address":    errors.ErrInput,
			},
		},
	}
	for testName, tc := range cases {
		t.Run(testName, func(t *testing.T) {
			err := tc.msg.Validate()
			for field, wantErr := range tc.wantErrs {
				assert.FieldError(t, err, field, wantErr)
			}
		})
	}
}
