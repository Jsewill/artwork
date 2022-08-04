/*
Package mint provides types and methods for minting NFTs on the chia blockchain.
*/
package mint

import (
	"fmt"
	"time"

	"github.com/Jsewill/chia/nft"
	"github.com/Jsewill/chia/rpc"
)

// Mint is a type which contains details pertaining to minting an NFT, and which can be used to mint an NFT.
type Mint struct {
	WalletId       uint
	MintWalletId   uint
	TargetAddress  string
	Royalty        uint
	RoyaltyAddress string
	Fee            uint
}

// NewMint creates a new *Mint with the supplied minting information. Fee and Royalty are used as defaults, which an individual NFT can override.
func NewMint(walletId uint, mintWalletId uint, targetAddress string, royalty uint, royaltyAddress string, fee uint) *Mint {
	return &Mint{
		WalletId:       walletId,
		MintWalletId:   mintWalletId,
		TargetAddress:  targetAddress,
		Royalty:        royalty,
		RoyaltyAddress: royaltyAddress,
		Fee:            fee,
	}
}

// MintRequest creates and returns a new rpc.MintRequest from its properties.
func (m *Mint) ToRequest() *rpc.MintRequest {
	return &rpc.MintRequest{
		WalletId:          m.MintWalletId,
		TargetAddress:     m.TargetAddress,
		RoyaltyPercentage: m.Royalty,
		RoyaltyAddress:    m.RoyaltyAddress,
		Fee:               m.Fee,
	}
}

// One attempts to mint one nft on the Chia Blockchain. Returns an error  if there was a critical failure, nil on success.
func (m *Mint) One(n nft.Nft) error {
	// Loop indefinitely, break on success or terminal error.
	for {
		// Check sync status
		status, err := &rpc.SyncStatusRequest{}.Send(rpc.Wallet)
		if err != nil {
			return fmt.Errorf("Wallet Sync Status request failed with the following error:", err)
		}
		if !status.Success {
			// Request was unsuccessful.
			fmt.Printf("Wallet Sync Status request was unsuccessful. Error: %s\nWaiting to retry.\n", status.Error)
			time.Sleep(10 * time.Second)
			continue
		}

		if !status.Synced {
			// Wait for synchronization
			fmt.Println("Wallet not synchronized. Waiting to retry.")
			time.Sleep(10 * time.Second)
			continue
		}
		// Check wallet balance
		balance, err := &rpc.WalletBalanceRequest{WalletId: m.WalletId}.Send(rpc.Wallet)
		if err != nil {
			return fmt.Errorf("XCH Wallet balance request failed with the following error:", err)
		}
		if !balance.Success {
			// Wait for wallet response.
			fmt.Printf("XCH Wallet Balance request was unsuccessful. Error: %s\nWaiting to retry.\n", balance.Error)
			time.Sleep(10 * time.Second)
			continue
		}
		if balance.WalletBalance.SpendableBalance < m.Fee {
			// We have enough to pay fees. Report and break out of the switch.
			fmt.Println("XCH wallet spendable balance is insufficient. Waiting to retry.\nFee: %d; Balance: %+v;\n", m.Fee, balance.WalletBalance)
			// Wait for spendable balance
			time.Sleep(10 * time.Second)
			continue
		}
		// Request was successful, and the spendable balance was sufficient.
		fmt.Println("Sufficient spendable balance: ", balance.WalletBalance.SpendableBalance)

		// Get Asset URIs and hash.
		assetUris, assetHash := n.Asset.URIs, n.Asset.Hash
		// Check NFT is well-formed and complete.
		if len(assetUris) == 0 {
			// No asset URIs. They are required to mint an NFT.
			return fmt.Errorf("At least one Asset URI is required to mint an NFT.")
		}
		if assetHash == "" {
			// No asset hash. This is require to work.
			return fmt.Errorf("An Asset hash is required to mint an NFT.")
		}
		// Get Metadata URIs and hash // @TODO: If no checks need to be done, these two lines could be removed and these vars directly assigned to the request struct members.
		metaUris, metaHash := n.Metadata.URIs, n.Metadata.Hash
		licenseUris, LicenseHash := n.License.URIs, n.License.Hash
		// Create request
		mrq := m.ToRequest()
		mrq.Uris, mrq.Hash = assetUris, assetHash
		mrq.MetaUris, mrq.MetaHash = metaUris, metaHash
		mrq.LicenseUris, mrq.LicenseHash = licenseUris, licenseHash
		// Time to mint!
		mr, err := mrq.Send(rpc.Wallet)
		if err != nil {
			return fmt.Errorf("Mint request failed with the following error: ", err)
		}
		if !mr.Success {
			fmt.Printf("Mint request was not successful. Error: %s\nWaiting to retry.\n", mr.Error)
			time.Sleep(10 * time.Second)
			continue
		}
		// Mint requested.
		fmt.Printf("Mint requested!\n")
	}

	return nil
}

// Many attempts to mint at least one nft on the Chia Blockchain. Returns an error if there was a critical failure, nil on success.
func (m *Mint) Many(c *nft.Collection) error {
	// For now, we'll just loop over collection items and use Mint.One()

	// Wait to mint another to avoid misses.
	time.Sleep(48 * time.Second)
	return fmt.Errorf("The Mint.Many function has not yet been written.")
}