# cota-nft-entries-syncer
The data syncer of [CoTA Script](https://github.com/nervina-labs/ckb-cota-scripts). 

# Prerequisites
* [MySQL](https://www.mysql.com/) 8.0 and above
* [CKB Node](https://github.com/nervosnetwork/ckb) v0.101.2 and above
* [golang](https://go.dev) 1.17.3 and above

# Quick Start
## Create Database
First you need to create a database, the default database name is `cota_entries`. You can adjust it according to your needs.

For the testnet, you can execute the following SQL statement to syn from a specific height: 
```sql
insert into check_infos (check_type, block_number, block_hash, created_at, updated_at) values (0, 3990570, '3d91bb118338cd70da1bc9adf8f9fbafcb04af5bc991550fe7fdde375872d5a1', now(), now());
```

## Start Node
Second you need to start a ckb node, You can refer to the following tutorial [Run a CKB Testnet Node](https://docs.nervos.org/docs/basics/guides/testnet).

## Configure Service
Update the [configuration file](configs/config.yaml) according to your specific situation.

app section in the config file has a mode config, can be configured as `wild` to turn on chase mode. 

## Local build
Enter this project directory and execute `make`.

## Run Service
Execute `bin/syncer`

## View Log
`tail -f storage/logs/app.logger`
