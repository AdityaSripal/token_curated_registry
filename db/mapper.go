package db

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/tendermint/go-amino"
	"github.com/AdityaSripal/token_curated_registry/types"
	crypto "github.com/tendermint/go-crypto"
	"reflect"
)

type RegistryMapper struct {
	ListingKey sdk.StoreKey

	CommitKey sdk.StoreKey

	RevealKey sdk.StoreKey

	VoteKey sdk.StoreKey

	Cdc *amino.Codec
}

func NewRegistryMapper(listingkey sdk.StoreKey, commitKey sdk.StoreKey, revealKey sdk.StoreKey, voteKey sdk.StoreKey, _cdc *amino.Codec) RegistryMapper {
	return RegistryMapper{
		ListingKey: listingkey,
		CommitKey: commitKey,
		RevealKey: revealKey,
		VoteKey: voteKey,
		Cdc: _cdc,
	}
}

// Will get Listing using unique identifier. Do not need to specify status
func (rm RegistryMapper) GetListing(ctx sdk.Context, identifier string) types.Listing {
	store := ctx.KVStore(rm.ListingKey)
	key, _ := rm.Cdc.MarshalBinary(identifier)
	val := store.Get(key)
	if val == nil {
		return types.Listing{}
	}
	listing := &types.Listing{}
	err := rm.Cdc.UnmarshalBinary(val, listing)
	if err != nil {
		panic(err)
	}
	return *listing
}

func (rm RegistryMapper) AddListing(ctx sdk.Context, identifier string, minBond int64) sdk.Error {
	store := ctx.KVStore(rm.ListingKey)

	// Cannot add an already existing listing
	if (!reflect.DeepEqual(rm.GetListing(ctx, identifier), types.Listing{})) {
		return sdk.NewError(2, 102, "Listing already exists")
	}

	newListing := types.Listing{
		Identifier: identifier,
		Status: "Pending",
		Display: false,
		Bond: minBond,
	}
	// Add listing with Pending Status
	key, _ := rm.Cdc.MarshalBinary(identifier)
	val, _ := rm.Cdc.MarshalBinary(newListing)
	store.Set(key, val)
	return nil
}

func (rm RegistryMapper) ChallengeListing(ctx sdk.Context, owner sdk.Address, challenger sdk.Address, identifier string, minBond int64) sdk.Error {
	store := ctx.KVStore(rm.ListingKey)
	voteStore := ctx.KVStore(rm.VoteKey)
	listing := rm.GetListing(ctx, identifier)

	// Cannot challenge a nonexistant listing
	if (!reflect.DeepEqual(rm.GetListing(ctx, identifier), types.Listing{})) {
		return sdk.NewError(2, 103, "Listing does not exist")
	}
	
	if listing.Status == "Challenged" {
		return sdk.NewError(2, 104, "Listing already Challenged")
	}

	listing.Status = "Challenged"

	key, _ := rm.Cdc.MarshalBinary(identifier)
	newVal, _ := rm.Cdc.MarshalBinary(listing)
	store.Set(key, newVal)

	ballot := types.Ballot{
		Identifier: identifier,
		Owner: owner,
		Challenger: challenger,
		Approve: 0,
		Deny: 0,
		Bond: minBond,
	}
	ballotVal, _ := rm.Cdc.MarshalBinary(ballot)
	voteStore.Set(key, ballotVal)
	return nil
}

func (rm RegistryMapper) CommitListing(ctx sdk.Context, owner sdk.Address, identifier string, commitment []byte) sdk.Error {
	commitStore := ctx.KVStore(rm.CommitKey)

	listing := rm.GetListing(ctx, identifier)
	if (!reflect.DeepEqual(rm.GetListing(ctx, identifier), types.Listing{})) {
		return sdk.NewError(2, 103, "Listing does not exist")
	}

	if listing.Status != "Challenged" {
		return sdk.NewError(2, 105, "Listing is not challenged")
	}

	voter := types.Voter{
		Owner: owner,
		Identifier: identifier,
	}
	key, _ := rm.Cdc.MarshalBinary(voter)
	commitStore.Set(key, commitment)
	return nil
}

func (rm RegistryMapper) RevealListing(ctx sdk.Context, owner sdk.Address, identifier string, vote bool, nonce []byte, power int64) sdk.Error {
	commitStore := ctx.KVStore(rm.CommitKey)
	revealStore := ctx.KVStore(rm.RevealKey)
	voteStore := ctx.KVStore(rm.VoteKey)

	listing := rm.GetListing(ctx, identifier)
	if (!reflect.DeepEqual(rm.GetListing(ctx, identifier), types.Listing{})) {
		return sdk.NewError(2, 103, "Listing does not exist")
	}

	if listing.Status != "Challenged" {
		return sdk.NewError(2, 105, "Listing is not challenged")
	}
	
	voter := types.Voter{
		Owner: owner,
		Identifier: identifier,
	}
	key, _ := rm.Cdc.MarshalBinary(voter)

	commitment := commitStore.Get(key)
	
	reveal := types.Vote{
		Choice: vote,
		Nonce: nonce,
		Power: power,
	}
	val, _ := rm.Cdc.MarshalBinary(reveal)
	if (!reflect.DeepEqual(crypto.Sha256(val), commitment)) {
		return sdk.NewError(2, 106, "Vote does not match commitment")
	}

	revealStore.Set(key, val)

	listingKey, _ := rm.Cdc.MarshalBinary(identifier)
	bz := voteStore.Get(listingKey)
	if bz == nil {
		return sdk.NewError(2, 107, "Ballot does not exist")
	}
	ballot := &types.Ballot{}
	err := rm.Cdc.UnmarshalBinary(bz, ballot)
	if err != nil {
		panic(err)
	}
	if vote {
		ballot.Approve += 1
	} else {
		ballot.Deny += 1
	}
	newBallot, _ := rm.Cdc.MarshalBinary(*ballot)
	voteStore.Set(listingKey, newBallot)
	return nil
}