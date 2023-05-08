/*
Package mint provides types and methods for minting NFTs on the chia blockchain.
*/
package mint

import (
	"fmt"
	"log"
	"time"

	"github.com/Jsewill/chia/nft"
	"github.com/Jsewill/chia/rpc"
)

// Mint is a type which contains default details pertaining to minting an NFT, and which can be used to mint an NFT.
type Mint struct {
	WalletId        uint
	TargetAddresses []string
	RoyaltyAddress  string
}

// NewMint creates a new *Mint with the supplied minting information.
func NewMint(walletId uint, targetAddresses []string, royaltyAddress string, royaltyPercentage float64) *Mint {
	return &Mint{
		WalletId:        walletId,
		TargetAddresses: targetAddresses,
		RoyaltyAddress:  royaltyAddress,
	}
}

// MintRequest creates and returns a new rpc.MintRequest from its properties.
func (m *Mint) MintBulkRequest() *rpc.MintBulkRequest {
	return &rpc.MintBulkRequest{
		WalletId:        walletId,
		TargetAddresses: m.TargetAddresses,
		RoyaltyAddress:  m.RoyaltyAddress,
	}
}

// MintRequest creates and returns a new rpc.MintRequest from its properties.
func (m *Mint) MintRequest() *rpc.MintRequest {
	return &rpc.MintRequest{
		WalletId:        walletId,
		TargetAddresses: m.TargetAddresses[0],
		RoyaltyAddress:  m.RoyaltyAddress,
	}
}

// Converts an nft.Nft to a MetadataListItem for bulk mint requests.
func MetadataListItemFromNft(n nft.Nft) MetadataListItem {
	mli := new(MetadataListItem)
	mli.Uris = n.Uris
	mli.MetaUris = n.Metadata.Uris
	mli.LicenseUris = n.License.Uris
	mli.Hash, _ = n.Hash()
	mli.MetaHash, _ = n.Metadata.Hash()
	mli.LicenseHash, _ = n.License.Hash()
	return *mli
}

// One attempts to mint one nft on the Chia Blockchain. Returns an error  if there was a critical failure, nil on success.
func (m *Mint) One(n nft.Nft) error {
	// Loop indefinitely, break on success or terminal error.
	for {
		// Check sync status
		status, err := &rpc.SyncStatusRequest{}.Send(rpc.Wallet)
		if err != nil {
			err = fmt.Errorf("Wallet Sync Status request failed with the following error: %s", err)
			logErr.Println(err)
			return err
		}
		if !status.Success {
			// Request was unsuccessful.
			logErr.Printf("Wallet Sync Status request was unsuccessful. Error: %s\nWaiting to retry.\n", status.Error)
			time.Sleep(10 * time.Second)
			continue
		}

		if !status.Synced {
			// Wait for synchronization
			logErr.Println("Wallet not synchronized. Waiting to retry.")
			time.Sleep(10 * time.Second)
			continue
		}
		// Check wallet balance
		balance, err := &rpc.WalletBalanceRequest{WalletId: 1}.Send(rpc.Wallet)
		if err != nil {
			err = fmt.Errorf("XCH Wallet balance request failed with the following error:", err)
			logErr.Println(err)
			return err
		}
		if !balance.Success {
			// Wait for wallet response.
			logErr.Printf("XCH Wallet Balance request was unsuccessful. Error: %s\nWaiting to retry.\n", balance.Error)
			time.Sleep(10 * time.Second)
			continue
		}
		if balance.WalletBalance.SpendableBalance < m.Fee {
			logErr.Println("XCH wallet spendable balance is insufficient. Waiting to retry.\nFee: %d; Balance: %+v;\n", m.Fee, balance.WalletBalance)
			// Wait for spendable balance
			time.Sleep(10 * time.Second)
			continue
		}
		// Request was successful, and the spendable balance was sufficient.
		log.Println("Sufficient spendable balance: ", balance.WalletBalance.SpendableBalance)
		// Get Asset URIs.
		assetUris := n.Asset.URIs
		// Check NFT is well-formed and complete.
		if len(assetUris) == 0 {
			// No asset URIs. They are required to mint an NFT.
			err = fmt.Errorf("At least one Asset URI is required to mint an NFT.")
			logErr.Println(err)
			return err
		}
		// Compute hashes if missing
		assetHash, err = n.Asset.Hash()
		if err != nil {
			err = fmt.Errorf("Unable to compute hash for asset: %+v", n.Asset)
			logErr.Println(err)
			return err
		}
		metaHash, err = n.Metadata.Hash()
		if err != nil {
			err = fmt.Errorf("Unable to compute hash for metadata: %+v", n.Metadata)
			logErr.Println(err)
			return err
		}
		licenseHash, err = n.License.Hash()
		if err != nil {
			err = fmt.Errorf("Unable to compute hash for asset: %+v", n.License)
			logErr.Println(err)
			return err
		}
		if assetHash == "" {
			// No asset hash. This is required to mint an NFT.
			err = fmt.Errorf("An Asset hash is required to mint an NFT.")
			logErr.Println(err)
			return err
		}
		// Get Metadata URIs and hash // @TODO: If no checks need to be done, these two lines could be removed and these vars directly assigned to the request struct members.
		metaUris, metaHash := n.Metadata.URIs, n.Metadata.Hash()
		licenseUris, LicenseHash := n.License.URIs, n.License.Hash()
		// Create request
		mrq := m.MintRequest()
		mrq.Uris, mrq.Hash = assetUris, assetHash
		mrq.MetaUris, mrq.MetaHash = metaUris, metaHash
		mrq.LicenseUris, mrq.LicenseHash = licenseUris, licenseHash
		if n.Royalty > 0.0 {
			royalty := n.Royalty
		} else {
			royalty := c.Royalty
		}
		mrq.RoyaltyPercentage = rpc.PercentageToRoyalty(royalty)
		mrq.Fee = n.Fee
		// Time to mint!
		log.Println("Sending mint request.")
		mr, err := mrq.Send(rpc.Wallet)
		if err != nil {
			err = fmt.Errorf("Mint request failed with the following error: ", err)
			logErr.Println(err)
			return err
		}
		if !mr.Success {
			logErr.Printf("Mint request was not successful. Error: %s\nWaiting to retry.\n", mr.Error)
			time.Sleep(10 * time.Second)
			continue
		}
		// Mint requested.
		log.Println("Mint requested!")
	}

	return nil
}

// Many attempts to mint at least one nft on the Chia Blockchain. Returns an error if there was a critical failure, nil on success.
func (m *Mint) Many(c *nft.Collection, bulk bool) error {
	if !bulk {
		for i, n := range c.Nfts {
			// Not bulk, use Mint.One()
			log.Printf("Starting on NFT #%d\n", i)
			err := m.One(n)
			if err != nil {
				err = fmt.Errorf("Failed to mint NFT #%d: %s\n", i, err)
				logErr.Println(err)
				// Is this enough to return on?
				return err
			}

			// Wait to mint another to avoid misses.
			time.Sleep(48 * time.Second)
		}
	} else {
		// @TODO: Bulk minting.
		mbrq := m.MintBulkRequest()
		mbrq.RoyaltyPercentage = rpc.PercentageToRoyalty(c.Royalty)
		mbrq.MintTotal = len(c.Nfts)
		mbrq.Fee = c.Fee

		// Check that the NFT wallet we're minting from has a DID associated with it.
		isDid, err := &rpc.NftWalletGetDidRequest{WalletId: 1}.Send(rpc.Wallet)
		if err != nil {
			err = fmt.Errorf("NFT Wallet Get DID request failed with the following error:", err)
			logErr.Println(err)
			return err
		}

		if !isDid.success && isDid.DidId != "" {
			mbrq.MintFromDid = true
		} else {
			mbrq.MintFromDid = false
		}

		// Check wallet balance
		for {
			balance, err := &rpc.WalletBalanceRequest{WalletId: 1}.Send(rpc.Wallet)
			if err != nil {
				err = fmt.Errorf("XCH Wallet balance request failed with the following error:", err)
				logErr.Println(err)
				return err
			}
			if !balance.Success {
				// Wait for wallet response.
				logErr.Printf("XCH Wallet Balance request was unsuccessful. Error: %s\nWaiting to retry.\n", balance.Error)
				time.Sleep(10 * time.Second)
				continue
			}
			if balance.WalletBalance.SpendableBalance < m.Fee {
				// Insufficient balance to cover fees.
				logErr.Println("XCH wallet spendable balance is insufficient. Waiting to retry.\nFee: %d; Balance: %+v;\n", m.Fee, balance.WalletBalance)
				// Wait for spendable balance
				time.Sleep(10 * time.Second)
				continue
			}
			// Request was successful, and the spendable balance was sufficient.
			log.Println("Sufficient spendable balance: ", balance.WalletBalance.SpendableBalance)
			break
		}

		// Prepare collection.
		for _, n := range c.Nfts {
			mbrq.MetadataList = append(mbrq.MetadataList, MetadataListItemFromNft(n))
		}

		// Get the SpendBundle.
		mbResponse, err := mbrq.Send(rpc.Wallet)
		if err != nil {
			err = fmt.Errorf("Bulk Mint request failed with the following error:", err)
			logErr.Println(err)
			return err
		}
		if !mbResponse.Success {
			// Couldn't create the SpendBundle for some reason.
			fmt.Errorf("SpendBundle creation was unsuccessful. Error: %s\nCannot continue.\n", mbResponse.Error)
			logErr.Printf(err)
			return
		}

		// Declare success. @TODO: Ask for confirmation.
		log.Printf("SpendBundle prepared. NFT IDs: %v", mbResponse.NftIdList)

		// Push transaction to the blockchain
		mintResults, err := &rpc.PushTxRequest{mbResponse.SpendBundle}.Send(rpc.FullNode)

		if err != nil || mintResults.Status != "SUCCESS" || mintResults.Success != true {
			fmt.Errorf("Failed to push SpendBundle. Error: %s\nCannot continue.\n", mintResults.Error)
			logErr.Printf(err)
			return
		}

	}
	return nil
}
