package main

import (
	"fmt"

	"github.com/STNS/libstns-go/libstns"
	"github.com/k0kubun/pp"
)

func main() {
	stns, err := libstns.NewSTNS("https://stns.lolipop.io/v1/", nil)
	if err != nil {
		panic(err)
	}

	user, err := stns.GetUserByName("pyama")
	if err != nil {
		panic(err)
	}
	pp.Println(user)

	signature, err := stns.Sign([]byte("secret message"))
	if err != nil {
		panic(err)
	}

	// it is ok
	fmt.Println(stns.VerifyWithUser("pyama", []byte("secret message"), signature))

	// verify error
	fmt.Println(stns.VerifyWithUser("pyama", []byte("invalid message"), signature))
}
