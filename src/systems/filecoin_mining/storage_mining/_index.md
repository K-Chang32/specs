---
menuTitle: Storage Miner
title: Storage Miner
entries:
  - mining_cycle
  - storage_miner_actor
  - mining_scheduler
---

{{<label storage_mining_subsystem>}}

TODO:

- rename "Storage Mining Worker" ?

Filecoin Storage Mining Subsystem

{{< readfile file="storage_mining_subsystem.id" code="true" lang="go" >}}

# Sector in StorageMiner State Machine (new one)

{{< diagram src="diagrams/sector_state.dot.svg" title="Sector State (new one)" >}}

{{< diagram src="diagrams/sector_state_legend.dot.svg" title="Sector State Legend (new one)" >}}

# Sector in StorageMiner State Machine (both)

{{< diagram src="diagrams/sector_fsm.dot.svg" title="Sector State Machine (both)" >}}
