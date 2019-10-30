package storage_market

import addr "github.com/filecoin-project/specs/systems/filecoin_vm/actor/address"
import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"
import deal "github.com/filecoin-project/specs/systems/filecoin_markets/deal"
import actor "github.com/filecoin-project/specs/systems/filecoin_vm/actor"
import msg "github.com/filecoin-project/specs/systems/filecoin_vm/message"
import vmr "github.com/filecoin-project/specs/systems/filecoin_vm/runtime"
import ipld "github.com/filecoin-project/specs/libraries/ipld"
import util "github.com/filecoin-project/specs/util"

////////////////////////////////////////////////////////////////////////////////
// Boilerplate
////////////////////////////////////////////////////////////////////////////////
type InvocOutput = msg.InvocOutput
type Runtime = vmr.Runtime
type Bytes = util.Bytes
type State = StorageMarketActorState

func (a *StorageMarketActorCode_I) State(rt Runtime) (vmr.ActorStateHandle, State) {
	h := rt.AcquireState()
	stateCID := h.Take()
	stateBytes := rt.IpldGet(ipld.CID(stateCID))
	if stateBytes.Which() != vmr.Runtime_IpldGet_FunRet_Case_Bytes {
		rt.Abort("IPLD lookup error")
	}
	state := DeserializeState(stateBytes.As_Bytes())
	return h, state
}
func Release(rt Runtime, h vmr.ActorStateHandle, st State) {
	checkCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.Release(checkCID)
}
func UpdateRelease(rt Runtime, h vmr.ActorStateHandle, st State) {
	newCID := actor.ActorSubstateCID(rt.IpldPut(st.Impl()))
	h.UpdateRelease(newCID)
}
func (st *StorageMarketActorState_I) CID() ipld.CID {
	panic("TODO")
}
func DeserializeState(x Bytes) State {
	panic("TODO")
}

////////////////////////////////////////////////////////////////////////////////

func (st *StorageMarketActorState_I) _generateStorageDealID(rt Runtime, storageDeal deal.StorageDeal) deal.DealID {
	// TODO
	var dealID deal.DealID
	return dealID
}

// Call by PublishStorageDeals and GetLastDealExpirationFromDealIDs (consider remove this)
// This is the check before a StorageDeal appears onchain
// It checks the following:
//   - verify deal did not expire when it is signed
//   - verify deal hits the chain before StartEpoch
//   - verify client and provider address and signature are correct (TODO may not be needed)
//   - verify StorageDealCollateral match requirements for MinimumStorageDealCollateral
//   - verify client and provider has sufficient balance
func (st *StorageMarketActorState_I) _verifyStorageDeal(rt Runtime, d deal.StorageDeal) bool {
	// TODO verify client and provider signature
	// TODO verify minimum StoragePrice, ProviderCollateralPerEpoch, and ClientCollateralPerEpoch
	// TODO: verify deal did not expire when it is signed

	currEpoch := rt.CurrEpoch()
	p := d.Proposal()

	// deal has started before publish
	if p.StartEpoch() < currEpoch {
		return false
	}

	// TODO: verify client and provider address and signature are correct (may not be needed)

	// verify StorageDealCollateral match requirements for MinimumStorageDealCollateral
	if p.ProviderCollateralPerEpoch() < deal.MIN_PROVIDER_DEAL_COLLATERAL_PER_EPOCH ||
		p.ClientCollateralPerEpoch() < deal.MIN_CLIENT_DEAL_COLLATERAL_PER_EPOCH {
		return false
	}

	// verify client and provider has sufficient balance
	clientBalanceA := st.Balances()[p.Client()].Available()
	providerBalanceA := st.Balances()[p.Provider()].Available()

	if clientBalanceA < (p.ClientBalanceRequirement()) ||
		providerBalanceA < p.ProviderBalanceRequirement() {
		return false
	}

	return true
}

// TODO: consider returning a boolean
func (st *StorageMarketActorState_I) _lockBalance(rt Runtime, addr addr.Address, amount actor.TokenAmount) {
	if amount < 0 {
		rt.Abort("negative amount.")
	}

	currBalance, found := st.Balances()[addr]
	if !found {
		rt.Abort("addr not found.")
	}

	currBalance.Impl().Available_ -= amount
	currBalance.Impl().Locked_ += amount
}

func (st *StorageMarketActorState_I) _unlockBalance(rt Runtime, addr addr.Address, amount actor.TokenAmount) {
	if amount < 0 {
		rt.Abort("negative amount.")
	}

	currBalance, found := st.Balances()[addr]
	if !found {
		rt.Abort("addr not found.")
	}

	currBalance.Impl().Locked_ -= amount
	currBalance.Impl().Available_ += amount
}

func (st *StorageMarketActorState_I) _lockFundsForStorageDeal(rt Runtime, deal deal.StorageDeal) {
	p := deal.Proposal()

	st._lockBalance(rt, p.Client(), p.ClientBalanceRequirement())
	st._lockBalance(rt, p.Provider(), p.ProviderBalanceRequirement())

}

func (st *StorageMarketActorState_I) _processStorageDealPayment(rt Runtime, deal deal.StorageDeal, duration block.ChainEpoch) {
	p := deal.Proposal()

	// move funds from locked in client to available in provider
	st.Balances()[p.Client()].Impl().Locked_ -= actor.TokenAmount(uint64(p.StoragePricePerEpoch()) * uint64(duration))
	st.Balances()[p.Provider()].Impl().Available_ += actor.TokenAmount(uint64(p.StoragePricePerEpoch()) * uint64(duration))
}

func (st *StorageMarketActorState_I) _settleExpiredStorageDeal(rt Runtime, deal deal.StorageDeal) {
	// TODO
}

func (st *StorageMarketActorState_I) _slashLockedFunds(rt Runtime, amount actor.TokenAmount) {
	// TODO
}

////////////////////////////////////////////////////////////////////////////////

func (a *StorageMarketActorCode_I) WithdrawBalance(rt Runtime, balance actor.TokenAmount) {
	h, st := a.State(rt)

	var msgSender addr.Address // TODO replace this from VM runtime

	if balance < 0 {
		rt.Abort("negative balance to withdraw.")
	}

	senderBalance, found := st.Balances()[msgSender]
	if !found {
		rt.Abort("sender address not found.")
	}

	if senderBalance.Available() < balance {
		rt.Abort("insufficient balance.")
	}

	senderBalance.Impl().Available_ = senderBalance.Available() - balance
	st.Balances()[msgSender] = senderBalance

	// TODO send funds to msgSender with `transferBalance` in VM runtime

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) AddBalance(rt Runtime) {
	h, st := a.State(rt)

	var msgSender addr.Address    // TODO replace this
	var balance actor.TokenAmount // TODO replace this

	// TODO subtract balance from msgSender
	// TODO add balance to StorageMarketActor
	if balance < 0 {
		rt.Abort("negative balance to add.")
	}

	senderBalance, found := st.Balances()[msgSender]
	if found {
		senderBalance.Impl().Available_ = senderBalance.Available() + balance
		st.Balances()[msgSender] = senderBalance
	} else {
		st.Balances()[msgSender] = &StorageParticipantBalance_I{
			Locked_:    0,
			Available_: balance,
		}
	}

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) PublishStorageDeals(rt Runtime, newStorageDeals []deal.StorageDeal) []PublishStorageDealResponse {
	h, st := a.State(rt)

	l := len(newStorageDeals)
	response := make([]PublishStorageDealResponse, l)

	// TODO: verify behavior here
	// some StorageDeal will pass and some will fail
	// if ealier StorageDeal consumes some balance such that
	// funds are no longer sufficient for later storage deals
	// all later storage deals will return error
	// TODO: confirm st here will be changing
	for i, newDeal := range newStorageDeals {
		if st._verifyStorageDeal(rt, newDeal) {
			st._lockFundsForStorageDeal(rt, newDeal)
			id := st._generateStorageDealID(rt, newDeal)
			st.Deals()[id] = newDeal
			response[i] = PublishStorageDealSuccess(id)
		} else {
			response[i] = PublishStorageDealError()
		}
	}

	UpdateRelease(rt, h, st)

	return response
}

func (a *StorageMarketActorCode_I) HandleCronAction(rt Runtime) {
	panic("TODO")
}

func (a *StorageMarketActorCode_I) SettleExpiredDeals(rt Runtime, storageDealIDs []deal.DealID) {
	// for dealID := range storageDealIDs {
	// Return the storage collateral
	// storageDeal := sma.Deals()[dealID]
	// storageCollateral := storageDeal.StorageCollateral()
	// provider := storageDeal.Provider()
	// assert(sma.Balances()[provider].Locked() >= storageCollateral)

	// // Move storageCollateral from locked to available
	// balance := sma.Balances()[provider]

	// sma.Balances()[provider] = &StorageParticipantBalance_I{
	// 	Locked_:    balance.Locked() - storageCollateral,
	// 	Available_: balance.Available() + storageCollateral,
	// }

	// // Delete reference to the deal
	// delete(sma.Deals_, dealID)
	// }
	panic("TODO")
}

func (a *StorageMarketActorCode_I) ProcessStorageDealsPayment(rt Runtime, dealIDs []deal.DealID, duration block.ChainEpoch) {
	h, st := a.State(rt)

	for _, dealID := range dealIDs {
		st._processStorageDealPayment(rt, st.Deals()[dealID], duration)
	}

	UpdateRelease(rt, h, st)
}

func (a *StorageMarketActorCode_I) SlashStorageDealsCollateral(rt Runtime, dealIDs []deal.DealID) {
	// for _, dealID := range storageDealIDs {
	// 	faultStorageDeal := sma.Deals()[dealID]
	// TODO remove locked funds and send slashed fund to TreasuryActor
	// TODO provider lose power for the FaultSet but not PledgeCollateral
	// }
	panic("TODO")
}

// Call by StorageMinerActor at CommitSector
func (a *StorageMarketActorCode_I) GetLastDealExpirationFromDealIDs(rt Runtime, dealIDs []deal.DealID) block.ChainEpoch {

	h, st := a.State(rt)

	var lastDealExpiration block.ChainEpoch
	for _, dealID := range dealIDs {
		deal, found := st.Deals()[dealID]
		if !found {
			rt.Abort("dealID not found.")
		}

		// this function is also called at sma.PublishStorageDeals()
		// TODO: decide if we want to refactor this
		// check balances, deal and proposal expiration, signatures
		isDealValid := st._verifyStorageDeal(rt, deal)
		if !isDealValid {
			rt.Abort("invalid deal.")
		}
		currExpiration := deal.Proposal().EndEpoch()
		if currExpiration > lastDealExpiration {
			lastDealExpiration = currExpiration
		}
	}

	Release(rt, h, st)

	return lastDealExpiration
}

func (a *StorageMarketActorCode_I) InvokeMethod(rt Runtime, method actor.MethodNum, params actor.MethodParams) InvocOutput {
	panic("TODO")
}
