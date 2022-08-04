package main

import (
	"fmt"

	"github.com/Jsewill/chia/nft"
)

func main() {
	fmt.Println("Artwork")
}

type Minter interface {
	One(nft.Nft) error
	Many([]*nft.Nft) error
}
