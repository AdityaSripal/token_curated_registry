package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank"
	types "github.com/AdityaSripal/token_curated_registry/types"
	db "github.com/AdityaSripal/token_curated_registry/db"
)

func NewCandidacyHandler(accountKeeper bank.Keeper, accountMapper sdk.AccountMapper, registryMapper db.RegistryMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		declareMsg := msg.(types.DeclareCandidacyMsg)
		account := accountMapper.GetAccount(ctx, declareMsg.Owner)
		_, _, err := accountKeeper.SubtractCoins(ctx, account, []Coin{declareMsg.Bond})
		if err != nil {
			return err.Result()
		}

		err2 := registryMapper.AddListing(ctx, declareMsg.Identifier, declareMsg.Bond.Amount)
		if err2 != nil {
			return err2.Result()
		}
		return sdk.Result{}
	}
}

func NewChallengeHandler(accountKeeper bank.Keeper, accountMapper sdk.AccountMapper, registryMapper db.RegistryMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		challengeMsg := msg.(types.ChallengeMsg)
		account := accountMapper.GetAccount(ctx, challengeMsg.Owner)
		_, _, err := accountKeeper.SubtractCoins(ctx, account, []Coin{declareMsg.Bond})
		if err != nil {
			return err.Result()
		}

		store := ctx.KVStore(registryMapper.ListingKey)
		key, _ := registryMapper.Cdc.MarshalBinary(challengeMsg.Identifier)
		bz := store.Get(key)

		if bz == nil {
			return sdk.NewError(2, 108, "Listing with given identifier does not exist").Result()
		}
		listing := &types.Listing{}
		err2 := registryMapper.Cdc.UnmarshalBinary(bz, listing)
		if err2 != nil {
			panic(err2)
		}

		err3 := registryMapper.ChallengeListing(ctx, listing.Owner, challengeMsg.Owner, challengeMsg.Identifier, challengeMsg.Bond.Amount)
		if err3 != nil {
			return err3.Result()
		}
		return sdk.Result{}
	}
}

func NewCommitHandler(accountMapper sdk.AccountMapper, registryMapper db.RegistryMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		commitMsg := msg.(types.CommitMsg)
		err := registryMapper.CommitListing(ctx, commitMsg.Owner, commitMsg.Identifier, commitMsg.Commitment)
		if err != nil {
			return err.Result()
		}
		return sdk.Result{}
	}
}

func NewRevealHandler(accountKeeper bank.Keeper, accountMapper sdk.AccountMapper, registryMapper db.RegistryMapper) sdk.Handler {
	return func(ctx sdk.Context, msg sdk.Msg) sdk.Result {
		revealMsg := msg.(types.RevealMsg)
		account := accountMapper.GetAccount(ctx, revealMsg.Owner)
		_, _, err := accountKeeper.SubtractCoins(ctx, account, []Coin{declareMsg.Bond})
		if err != nil {
			return err.Result()
		}

		err2 := registryMapper.RevealListing(ctx, revealMsg.Owner, revealMsg.Identifier, revealMsg.Vote, revealMsg.Nonce, revealMsg.Bond.Amount)
		if err2 != nil {
			return err2.Result()
		}
		return sdk.Result{}
	}
}

