package main

import (
	"flag"
	"fmt"
	"io"
	"time"

	"github.com/iov-one/weave"
	customd "github.com/iov-one/blog-tutorial/cmd/customd/app"
	"github.com/iov-one/blog-tutorial/x/custom"
)

func cmdCreateState(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), `
Create a transaction for creating a custom state.
		`)
		fl.PrintDefaults()
	}
	var (
		innerState1 = fl.Int64("innerstate", 0, "inner state 1")
		innerState2 = fl.Int64("innerstate", 0, "inner state 2")
		addressFl   = flAddress(fl, "address", "", "Address representation")
	)
	fl.Parse(args)

	innerState := custom.InnerState{
		St1: *innerState1,
		St2: *innerState2,
	}

	msg := custom.CreateStateMsg{
		Metadata:   &weave.Metadata{Schema: 1},
		InnerState: &innerState,
		Address:    *addressFl,
	}

	if err := msg.Validate(); err != nil {
		return fmt.Errorf("given data produce an invalid invalid message: %s", err)
	}

	tx := &customd.Tx{
		Sum: &customd.Tx_CustomCreateStateMsg{
			CustomCreateStateMsg: &msg,
		},
	}
	_, err := writeTx(output, tx)
	return err
}

func cmdCreateTimedState(input io.Reader, output io.Writer, args []string) error {
	fl := flag.NewFlagSet("", flag.ExitOnError)
	fl.Usage = func() {
		fmt.Fprintln(flag.CommandLine.Output(), `
Create a transaction for creating a timed custom state.
		`)
		fl.PrintDefaults()
	}
	var (
		innerStateEnum = fl.Int64("innerstateenum", 3, "Invalid = 0, CaseOne = 1, CaseTwo = 2")
		str            = fl.String("string", "", "string must start with 'cstm' to be valid")
		bytes          = fl.String("bytes", "", "Byte representation")
		deleteAt       = fl.Int64("deleteat", 0, "Delete at represents the unix time of deletion of custom state")
	)
	fl.Parse(args)

	var ise custom.InnerStateEnum
	switch *innerStateEnum {
	case 0:
		ise = custom.InnerStateEnum_Invalid
	case 1:
		ise = custom.InnerStateEnum_CaseOne
	case 2:
		ise = custom.InnerStateEnum_CaseTwo
	default:
		return fmt.Errorf("unknown inner state enumeration = %d, must be 0, 1 or 2", innerStateEnum)
	}

	da := weave.AsUnixTime(time.Unix(*deleteAt, 0))

	msg := custom.CreateTimedStateMsg{
		Metadata:       &weave.Metadata{Schema: 1},
		Str:            *str,
		Byte:           []byte(*bytes),
		InnerStateEnum: ise,
		DeleteAt:       da,
	}

	if err := msg.Validate(); err != nil {
		return fmt.Errorf("given data produce an invalid invalid message: %s", err)
	}

	tx := &customd.Tx{
		Sum: &customd.Tx_CustomCreateTimedStateMsg{
			CustomCreateTimedStateMsg: &msg,
		},
	}
	_, err := writeTx(output, tx)
	return err
}
