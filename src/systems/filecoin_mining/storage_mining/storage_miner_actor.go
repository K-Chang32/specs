package storage_mining

import base_mining "github.com/filecoin-project/specs/systems/filecoin_mining"
import sector "github.com/filecoin-project/specs/systems/filecoin_mining/sector"
import sealing "github.com/filecoin-project/specs/systems/filecoin_mining/sector"

func (sm *StorageMinerActor_I) SubmitPoSt(postProof base_mining.PoStProof, sectorStateSets sector.SectorStateSets) {
	panic("TODO")
}

// TODO: currently deals must be posted on chain via sma.PublishStorageDeals
// as an optimization, in the future CommitSector could contain a list of
// deals that are not published yet.
func (sm *StorageMinerActor_I) CommitSector(onChainInfo sealing.OnChainSealVerifyInfo) {
	// TODO: sm.verifySeal(sectorID SectorID, comm sealing.OnChainSealVerifyInfo, proof SealProof)

	// TODO: you don't need to store the proof and the randomseed in the state tree
	// just verify it and drop it, just SealCommitments{CommD, CommR, DealIDs}
	// TODO: @porcuquine verifies
	currSector, found := sm.Sectors()[onChainInfo.SectorNumber()]
	if !found {
		sm.Sectors()[onChainInfo.SectorNumber()] = []sealing.OnChainSealVerifyInfo{onChainInfo}
	} else {
		newSector := append(currSector, onChainInfo)
		sm.Sectors()[onChainInfo.SectorNumber()] = newSector
	}

	// TODO write state change
}
