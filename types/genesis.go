package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
)

type GenesisState struct {
	Accounts []*GenesisAccount `json:"accounts"`
}

// GenesisAccount doesn't need pubkey or sequence
type GenesisAccount struct {
	Address sdk.Address `json:"address"`
	Coins   sdk.Coins   `json:"coins"`
}

func NewGenesisAccount(aa *auth.BaseAccount) *GenesisAccount {
	return &GenesisAccount{
		Address: aa.Address,
		Coins:   aa.Coins.Sort(),
	}
}

// convert GenesisAccount to AppAccount
func (ga *GenesisAccount) ToAccount() (acc *auth.BaseAccount, err error) {
	baseAcc := auth.BaseAccount{
		Address: ga.Address,
		Coins:   ga.Coins.Sort(),
	}
	return &baseAcc, nil
}