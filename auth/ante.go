package auth

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"reflect"
)

func NewAnteHandler(accountMapper sdk.AccountMapper, mindenom uint64) sdk.AnteHandler {
	return func(ctx sdk.Context, tx sdk.Tx) (_ sdk.Context, _ sdk.Result, abort bool) {
		sigs := tx.GetSignatures()
		if len(sigs) != 1 {
			return ctx,
				sdk.ErrUnauthorized("no signers").Result(),
				true
		}

		msg := tx.GetMsg()

		_, ok := tx.(sdk.StdTx)
		if !ok {
			return ctx, sdk.ErrInternal("tx must be in form of StdTx").Result(), true
		}

		signerAddr := msg.GetSigners()[0]

		if !reflect.DeepEqual(sigs[0].PubKey.Address().Bytes(), signerAddr.Bytes()) {
			return ctx, sdk.ErrInternal("Wrong signer address").Result(), true
		}

		if !sigs[0].PubKey.VerifyBytes(msg.GetSignBytes(), sigs[0].Signature) {
			return ctx, sdk.ErrInternal("Invalid Signature").Result(), true
		}

		if uint64(accountMapper.GetAccount(ctx, signerAddr).GetCoins().AmountOf("RegistryCoin")) < mindenom {
			return ctx, sdk.ErrInternal("Must bond at least minimum bond for candidacy").Result(), true
		}

		return ctx, sdk.Result{}, false
	}
}