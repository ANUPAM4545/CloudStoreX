# Consistency Verification

CloudStoreX runs continuous background verification to ensure metadata catalog integrity matches the actual physical provider storage.

## Types of Checks
- Object Presence (Missing from provider)
- Object Zombie (Missing from metadata)
- Object Size/Checksum mismatch

Discrepancies emit `ConsistencyCheckFailed` events which trigger auto-repairs.
