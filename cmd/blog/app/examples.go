package blog

import (
	"encoding/hex"
	"strings"
	"time"

	"github.com/iov-one/blog-tutorial/x/blog"
	"github.com/iov-one/weave"
	"github.com/iov-one/weave/coin"
	"github.com/iov-one/weave/commands"
	"github.com/iov-one/weave/crypto"
	"github.com/iov-one/weave/x/cash"
	"github.com/iov-one/weave/x/sigs"

	"github.com/iov-one/weave/weavetest"
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
	ticker := "BLOG"
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

	createUserMsg := &blog.CreateUserMsg{
		Metadata: &weave.Metadata{Schema: 1},
		Username: "Crpto0X",
		Bio:      "Best hacker in the universe",
	}
	createBlogMsg := &blog.CreateBlogMsg{
		Metadata:    &weave.Metadata{Schema: 1},
		Title:       "insanely good title",
		Description: "best description in the existence",
	}

	blogID := weavetest.SequenceID(1)

	changeOwnerMsg := &blog.ChangeBlogOwnerMsg{
		Metadata: &weave.Metadata{Schema: 1},
		BlogID:   blogID,
		NewOwner: weavetest.NewCondition().Address(),
	}

	now := weave.AsUnixTime(time.Now())
	future := now.Add(time.Hour)

	createArticleMsg := &blog.CreateArticleMsg{
		Metadata: &weave.Metadata{Schema: 1},
		BlogID:   blogID,
		Title:    "insanely good title",
		Content:  "best content in the existence",
		DeleteAt: future,
	}
	articleID := weavetest.SequenceID(1)
	deleteArticleMsg := &blog.DeleteArticleMsg{
		Metadata:  &weave.Metadata{Schema: 1},
		ArticleID: articleID,
	}

	return []commands.Example{
		{Filename: "wallet", Obj: wallet},
		{Filename: "priv_key", Obj: source},
		{Filename: "pub_key", Obj: pub},
		{Filename: "user", Obj: user},
		{Filename: "send_msg", Obj: msg},
		{Filename: "unsigned_tx", Obj: &unsigned},
		{Filename: "signed_tx", Obj: &tx},
		{Filename: "blog_create_user_msg", Obj: createUserMsg},
		{Filename: "blog_create_blog_msg", Obj: createBlogMsg},
		{Filename: "blog_create_article_msg", Obj: createArticleMsg},
		{Filename: "blog_delete_article_msg", Obj: deleteArticleMsg},
		{Filename: "blog_change_blog_owner_msg", Obj: changeOwnerMsg},
	}
}
