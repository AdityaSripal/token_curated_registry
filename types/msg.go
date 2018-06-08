package types

import (
	"encoding/json"
	sdk "github.com/cosmos/cosmos-sdk/types"
	amino "github.com/tendermint/go-amino"
)

const (
	tokenName = "RegistryCoin"
)

type DeclareCandidacyMsg struct {
	Owner sdk.Address
	Identifier string
	Bond sdk.Coin
}

func NewDeclareCandidacyMsg(owner sdk.Address, identifier string, bond sdk.Coin) DeclareCandidacyMsg  {
	return DeclareCandidacyMsg{
		Owner: owner,
		Identifier: identifier,
		Bond: bond,
	}
}

func (msg DeclareCandidacyMsg) Type() string {
	return "DeclareCandidacy"
}

func (msg DeclareCandidacyMsg) ValidateBasic() sdk.Error {
	if (msg.Bond.Amount <= 0 || msg.Bond.Denom != tokenName) {
		return sdk.NewError(2, 101, "Must submit a bond in RegistryCoins")
	}
	return nil
}

func (msg DeclareCandidacyMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg DeclareCandidacyMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Owner}
}

type ChallengeMsg struct {
	Owner sdk.Address
	Identifier string
	Bond sdk.Coin
}

func NewChallengeMsg(owner sdk.Address, identifier string, bond sdk.Coin) ChallengeMsg {
	return ChallengeMsg{
		Owner: owner,
		Identifier: identifier,
		Bond: bond,
	}
}

func (msg ChallengeMsg) Type() string {
	return "Challenge"
}

func (msg ChallengeMsg) ValidateBasic() sdk.Error {
	if (msg.Bond.Amount <= 0 || msg.Bond.Denom != tokenName) {
		return sdk.NewError(2, 101, "Must submit a bond in RegistryCoins")
	}
	return nil
}

func (msg ChallengeMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg ChallengeMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Owner}
}

type CommitMsg struct {
	Owner sdk.Address
	Identifier string
	Commitment []byte
}

func NewCommitMsg(owner sdk.Address, identifier string, commitment []byte) CommitMsg {
	return CommitMsg{
		Owner: owner,
		Identifier: identifier,
		Commitment: commitment,
	}
}

func (msg CommitMsg) Type() string {
	return "Commit"
}

func (msg CommitMsg) ValidateBasic() sdk.Error {
	return nil
}

func (msg CommitMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg CommitMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Owner}
}

type RevealMsg struct {
	Owner sdk.Address
	Identifier string
	Vote bool
	Nonce []byte
	Bond sdk.Coin
}

func (msg RevealMsg) Type() string {
	return "Reveal"
}

func (msg RevealMsg) ValidateBasic() sdk.Error {
	if (msg.Bond.Amount <= 0 || msg.Bond.Denom != tokenName) {
		return sdk.NewError(2, 101, "Must submit a bond in RegistryCoins")
	}
	return nil
}

func (msg RevealMsg) GetSignBytes() []byte {
	b, err := json.Marshal(msg)
	if err != nil {
		panic(err)
	}
	return b
}

func (msg RevealMsg) GetSigners() []sdk.Address {
	return []sdk.Address{msg.Owner}
}

func RegisterAmino(cdc *amino.Codec) {
	cdc.RegisterConcrete(DeclareCandidacyMsg{}, "types/DeclareCandidacyMsg", nil)
	cdc.RegisterConcrete(ChallengeMsg{}, "types/ChallengeMsg", nil)
	cdc.RegisterConcrete(CommitMsg{}, "types/CommintMsg", nil)
	cdc.RegisterConcrete(RevealMsg{}, "types/RevealMsg", nil)
	cdc.RegisterConcrete(Listing{}, "types/Listing", nil)
	cdc.RegisterConcrete(Voter{}, "types/Voter", nil)
	cdc.RegisterConcrete(Vote{}, "types/Vote", nil)
	cdc.RegisterConcrete(Ballot{}, "types/Ballot", nil)
}