package main

import (
	"github.com/iov-one/weave/coin"
	"github.com/iov-one/weave/errors"
	"github.com/iov-one/weave/x/msgfee"
)

func msgfeeConf(nodeUrl string, msgPath string) (*coin.Coin, error) {
	store := tendermintStore(nodeUrl)
	b := msgfee.NewMsgFeeBucket()
	var fee msgfee.MsgFee
	switch err := b.One(store, []byte(msgPath), &fee); {
	case err == nil:
		return &fee.Fee, nil
	case errors.ErrNotFound.Is(err):
		return nil, nil
	default:
		return nil, errors.Wrap(err, "cannot get fee")
	}
}

