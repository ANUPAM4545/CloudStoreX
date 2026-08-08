# Backups & Retention

CloudStoreX supports scheduled Full and Incremental backups, routing data to cold storage classes like AWS S3 Glacier or Google Cloud Storage Coldline.

## Process

1. Policies define the schedule and targets.
2. Background workers execute the extraction.
3. Retention engines automatically prune expired backups based on configured retention lifecycles.
