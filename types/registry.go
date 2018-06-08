package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
)

type Listing struct {
	Identifier string
	Owner sdk.Address
	Status string
	Display bool
	Bond int64
}

func EncodeListing(cdc *amino.Codec, listing Listing) []byte {
	b, err := cdc.MarshalBinary(listing)
	if err != nil {
		panic(err)
	}
	return b
}

// Create new Voter for address on each Listing
type Voter struct {
	Owner sdk.Address
	Identifier string 
}

// Vote revealed during reveal phase
type Vote struct {
	Choice bool
	Nonce []byte
	Power int64
}

type Ballot struct {
	Identifier string
	Owner sdk.Address
	Challenger sdk.Address
	Approve uint64
	Deny uint64
	Bond int64
}