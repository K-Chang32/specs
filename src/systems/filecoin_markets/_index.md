---
menuIcon: ⚖️
menuTitle: "**Market**"
title: "Markets in Filecoin"
entries:
- order
- deal
- storage_market
- retrieval_market
---

{{< incTocMap "/docs/systems/filecoin_markets" 1 >}}

The Filecoin project is a protocol, a platform, and a marketplace. There are two major components to Filecoin markets, storage market and retrieval market. While both markets are expected to happen primarily off the blockchain, storage deals made in storage market will be published on chain and enforced by the protocol. Storage deal negotiation and order matching are expected to happen off chain in the first version of Filecoin. Retrieval deals are also negotiated off chain and executed with micropayments between transacting parties in payment channels.

Even though most of the market actions happen off the blockchain, there are on-chain invariant and structure that create economic structure for network success and allow for positive emergent behavior. Storage Mining in Filecoin can be compared to maintaining a storage cargo container (reference to `Sector`) with storage deals sitting in them. There must be at least one active deal in a `Sector` for the `Sector` to count towards a miner's power. The `Sector` is only considered Active after submitting it first PoSt successfully (reference to Sector FSM). A` Sector` no longer counts for power when it leaves the Active state either through fault or expiration. A `Sector` expires when all the deals in the `Sector` have expired and all `StorageDealCollateral` will only be returned then. In the first version of Filecoin, deals are immutable once added to the chain.
