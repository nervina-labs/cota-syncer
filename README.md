# cota-syncer
The data syncer of [CoTA Script](https://talk.nervos.org/t/rfc-cota-a-compact-token-aggregator-standard-for-extremely-low-cost-nfts-and-fts/6338). 

# Prerequisites
* [MySQL](https://www.mysql.com/) 8.0 and above
* [CKB Node](https://github.com/nervosnetwork/ckb) v0.101.2 and above
* [golang](https://go.dev) 1.17.3 and above

# Quick Start
## Create Database
First you need to create a database, the default database name is `cota_entries`. You can adjust it according to your needs.

For the testnet, you can execute the following SQL statement to syn from a specific height: 
```sql
insert into check_infos (check_type, block_number, block_hash, created_at, updated_at) values (0, 4163980, 'ab6d9453628ee854062615acf05f899e8c84e4e61d417d0b13bbed128a862e23', now(), now());
insert into check_infos (check_type, block_number, block_hash, created_at, updated_at) values (1, 4163980, 'ab6d9453628ee854062615acf05f899e8c84e4e61d417d0b13bbed128a862e23', now(), now());
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
