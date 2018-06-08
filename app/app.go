package app

import (
	"encoding/json"
	bam "github.com/cosmos/cosmos-sdk/baseapp"
	sdk "github.com/cosmos/cosmos-sdk/types"
	abci "github.com/tendermint/abci/types"
	"github.com/tendermint/go-amino"
	crypto "github.com/tendermint/go-crypto"
	cmn "github.com/tendermint/tmlibs/common"
	dbm "github.com/tendermint/tmlibs/db"
	"github.com/tendermint/tmlibs/log"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/cosmos-sdk/x/bank"
	handle "github.com/AdityaSripal/token_curated_registry/auth"
	"github.com/AdityaSripal/token_curated_registry/db"
	//rlp "github.com/ethereum/go-ethereum/rlp"
)

const (
	appName = "Registry"
)

// Extended ABCI application
type ChildChain struct {
	*bam.BaseApp

	cdc *amino.Codec

	minDeposit uint64

	applyStage uint64

	commitStage uint64

	revealStage uint64

	dispensationPct float32

	quorum float32

	// keys to access the substores
	capKeyMainStore *sdk.KVStoreKey
	capKeyListings *sdk.KVStoreKey
	capKeyCommits *sdk.KVStoreKey
	capKeyReveals *sdk.KVStoreKey
	capKeyVotes *sdk.KVStoreKey

	registryMapper db.RegistryMapper

	// Manage addition and subtraction of account balances
	accountMapper sdk.AccountMapper
	accountKeeper sdk.AccountKeeper
}

func NewChildChain(logger log.Logger, db dmb.DB, mindeposit uint64, applystage uint64, commitstage uint64, revealstage uint64, dispensationpct float32, _quorum float32) *ChildChain {
	cdc := MakeCodec()
	var app = &ChildChain{
		BaseApp: bam.NewBaseApp(appName, cdc, logger, db),
		cdc: cdc,
		minDeposit: mindeposit,
		applyStage: applystage,
		commitStage: commitstage,
		revealStage: revealstage,
		dispensationPct: dispensationpct,
		quorum: _quorum,
		capKeyMainStore: sdk.NewKVStoreKey("main"),
		capKeyAccount: sdk.NewKVStoreKey("acc")
		capKeyListings: sdk.NewKVStoreKey("listings"),
		capKeyCommits: sdk.NewKVStoreKey("commits"),
		capKeyReveals: sdk.NewKVStoreKey("reveals")
		capKeyVotes: sdk.NewKVStoreKey("votes"),
	}

	app.registryMapper = db.NewRegistryMapper(app.capKeyListings, app.capKeyCommits, app.capKeyReveals, app.capKeyVotes, app.cdc)
	app.accountMapper = auth.NewAccountMapper(app.cdc, app.capKeyAccount, &auth.BaseAccount{})
	app.accountKeeper =  bank.NewKeeper(app.accountMapper)

	app.Router()
		.addRoute("DeclareCandidacy", handle.NewCandidacyHandler(app.accountKeeper, app.accountMapper, app.registryMapper))
		.addRoute("Challenge", handle.NewChallengeHandler(app.accountKeeper, app.accountMapper, app.registryMapper))
		.addRoute("Commit", handle.NewCommitHandler(app.accountMapper, app.registryMapper))
		.addRoute("Reveal", handle.NewRevealHandler(app.accountKeeper, app.accountMapper, app.registryMapper))

	app.SetTxDecoder(app.txDecoder)
	app.MountStoresIAVL(app.capKeyMainStore, app.capKeyAccount, app.capKeyListings, app.capKeyCommits, app.capKeyReveals, app.capKeyVotes)
	app.SetAnteHandler(handle.NewAnteHandler(app.accountMapper, app.minDeposit))
}

func (app *ChildChain) initChainer(ctx sdk.Context, req abci.RequestInitChain) abci.ResponseInitChain {
	stateJSON := req.AppStateBytes

	genesisState := new(types.GenesisState)
	err := app.cdc.UnmarshalJSON(stateJSON, genesisState)
	if err != nil {
		panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
		// return sdk.ErrGenesisParse("").TraceCause(err, "")
	}

	for _, gacc := range genesisState.Accounts {
		acc, err := gacc.ToAppAccount()
		if err != nil {
			panic(err) // TODO https://github.com/cosmos/cosmos-sdk/issues/468
			//	return sdk.ErrGenesisParse("").TraceCause(err, "")
		}
		app.accountMapper.SetAccount(ctx, acc)
	}
	return abci.ResponseInitChain{}
}

func (app *ChildChain) txDecoder(txBytes []byte) (sdk.Tx, sdk.Error) {
	var tx = sdk.StdTx
	err := json.Unmarshal(txBytes, &tx)
	if err != nil {
		return nil, sdk.ErrTxDecode("")
	}
	return tx, nil
}

func MakeCodec() *amino.Codec {
	cdc := amino.NewCodec()
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	types.RegisterAmino(cdc)
	return cdc
}