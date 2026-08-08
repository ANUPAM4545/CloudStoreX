# Self-Healing

CloudStoreX automatically repairs broken, missing, or corrupted objects discovered during consistency checks.

## Process
1. Subscriber listens to `ConsistencyCheckFailed` reliability events.
2. The Self-Healing service queries replicas or cold backup vaults.
3. `SELF_HEAL_JOB` workers pull data from healthy targets to overwrite the broken provider target.
