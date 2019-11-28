package storage_mining

import block "github.com/filecoin-project/specs/systems/filecoin_blockchain/struct/block"

// TODO: placeholder epoch value -- this will be set later
const MAX_PROVE_COMMIT_SECTOR_PERIOD = block.ChainEpoch(3)    // placeholder
const MAX_SURPRISE_POST_RESPONSE_PERIOD = block.ChainEpoch(4) // placeholder
const POST_CHALLENGE_TIME = block.ChainEpoch(1)               // placeholder
const PROVING_PERIOD = block.ChainEpoch(2)                    // placeholder
// how many times per PP should a miner get challenged on expectation
const SURPRISE_CHALLENGE_FREQUENCY = 2 // placeholder
// how far into their PP does a miner get their first challenge
const SUPRISE_NO_CHALLENGE_PERIOD = PROVING_PERIOD / SURPRISE_CHALLENGE_FREQUENCY
