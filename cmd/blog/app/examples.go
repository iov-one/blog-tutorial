package blog

import (
	"encoding/hex"
	"strings"

	"github.com/iov-one/weave"
	"github.com/iov-one/blog-tutorial/x/custom"
	"github.com/iov-one/weave/coin"
	"github.com/iov-one/weave/commands"
	"github.com/iov-one/weave/crypto"
	"github.com/iov-one/weave/x/cash"
	"github.com/iov-one/weave/x/sigs"
)

// we fix the private keys here for deterministic output with the same encoding
// these are not secure at all, but the only point is to check the format,
// which is easier when everything is reproduceable.
var (
	source     = makePrivKey("1234567890")
	dst        = makePrivKey("F00BA411").PublicKey().Address()
	randomAddr = makePrivKey("00CAFE00F00D").PublicKey().Address()
)

// makePrivKey repeats the string as long as needed to get 64 digits, then
// parses it as hex. It uses this repeated string as a "random" seed
// for the private key.
//
// nothing random about it, but at least it gives us variety
func makePrivKey(seed string) *crypto.PrivateKey {
	rep := 64/len(seed) + 1
	in := strings.Repeat(seed, rep)[:64]
	bin, err := hex.DecodeString(in)
	if err != nil {
		panic(err)
	}
	return crypto.PrivKeyEd25519FromSeed(bin)
}

// Examples generates some example structs to dump out with testgen
func Examples() []commands.Example {
	ticker := "CSTM"
	wallet := &cash.Set{
		Metadata: &weave.Metadata{Schema: 1},
		Coins: []*coin.Coin{
			{Whole: 150, Fractional: 567000, Ticker: ticker},
		},
	}

	pub := source.PublicKey()
	addr := pub.Address()
	user := &sigs.UserData{
		Metadata: &weave.Metadata{Schema: 1},
		Pubkey:   pub,
		Sequence: 17,
	}

	amt := coin.NewCoin(250, 0, ticker)
	msg := &cash.SendMsg{
		Metadata:    &weave.Metadata{Schema: 1},
		Amount:      &amt,
		Destination: dst,
		Source:      addr,
		Memo:        "Test payment",
	}

	unsigned := Tx{
		Sum: &Tx_CashSendMsg{msg},
	}
	tx := unsigned
	sig, err := sigs.SignTx(source, &tx, "test-123", 17)
	if err != nil {
		panic(err)
	}
	tx.Signatures = []*sigs.StdSignature{sig}

	createTimedStateMsg := &custom.CreateTimedStateMsg{
		Metadata:       &weave.Metadata{Schema: 1},
		InnerStateEnum: custom.InnerStateEnum_CaseOne,
		Str:            "cstm_str",
		Byte:           []byte{0, 1},
	}

	createStateMsg := &custom.CreateStateMsg{
		Metadata:   &weave.Metadata{Schema: 1},
		InnerState: &custom.InnerState{St1: 1, St2: 1},
		Address:    randomAddr,
	}

	return []commands.Example{
		{Filename: "wallet", Obj: wallet},
		{Filename: "priv_key", Obj: source},
		{Filename: "pub_key", Obj: pub},
		{Filename: "user", Obj: user},
		{Filename: "send_msg", Obj: msg},
		{Filename: "unsigned_tx", Obj: &unsigned},
		{Filename: "signed_tx", Obj: &tx},
		{Filename: "custom_create_timed_state_msg", Obj: createTimedStateMsg},
		{Filename: "custom_create_state_msg", Obj: createStateMsg},
	}
}
