package main

import (
	"fmt"

	"github.com/Jsewill/chia/nft"
)

type Minter interface {
	One(nft.Nft) error
	Many([]*nft.Nft) error
}

type Generator interface {
	Generate() (nft.Collection, error)
}

func main() {
	// Remove as this gets filled out
	fmt.Println("Artwork")
}
