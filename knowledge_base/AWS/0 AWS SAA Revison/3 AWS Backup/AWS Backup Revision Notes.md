
## AWS Backup
- AWS Backup is a **fully managed**, **centralized and automated solution** for backing up AWS Resources like:
    - EBS
    - EFS
    - EC2
    - RDS
    - DynamoDB
    - Storage Gateway
    - Aurora
- Uses **default backup mechanism** for each of these AWS services. For customized backup mechanism, for example for EBS Volume, can use AWS lambda and Event bridge

## How to create backup
- **Create Vaults**
    - Vaults are containers to organize backups - grouping backups for same app together
    - if we have an application named "my-app", we can store EBS, EFS, EC2 and RDS backups around this application in a same vault
- Create **Backup plans**
    - **Schedule** when to take backups
    - **Frequency** of backups
    - **Retention period** - can be set to always or specicy fixed period