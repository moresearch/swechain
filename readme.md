# swechain
**swechain** is a blockchain built using Cosmos SDK and Tendermint and created with [Ignite CLI](https://ignite.com/cli).

## Prerequisites

Golang


`curl -sSL https://get.ignite.com/cli\! | sudo bash`


GEX is a real time in-terminal explorer for Cosmos SDK blockchains. 
`go install github.com/cosmos/gex@latest` 


## Get started


```
ignite chain build --release --release.targets="linux:amd64,windows:amd64"
```


```
swechaind start
```


Run  to view transcations(tx) in real-time.

```
gex explorer 
```


### Configure

Your blockchain in development can be configured with `config.yml` to edit balances or add accounts for example. 


## Usage Example:
### Setup chain Id
swechaind config set client chain-id swechain

### Create accounts
swechaind keys add alice
swechaind keys add bob

### Check Initial Balances
swechaind query bank balances alice --output json
swechaind query bank balances bob --output json

### Create an auction (Alice)
swechaind tx issuemarket create-auction "BUG-123" "Fix critical security vulnerability" "open" "" --from alice --yes --output json

### Place bids (Bob)
BOB=$(swechaind keys show bob -a)
swechaind tx issuemarket create-bid  "0" "0" "$BOB" "5000" "Will fix in 7 days" --from bob --yes --output json
swechaind tx issuemarket create-bid  "0" "0" "$BOB" "4000" "Will fix in 6 days" --from bob --yes --output json

### List all bids and filter for auction 1
swechaind query issuemarket list-bid --output json | jq '.bid | .[] | select(.auctionId == "1")'

### Close the auction (Alice - automatically selects lowest bidder)
swechaind tx issuemarket update-auction 0 "BUG-123" "Fix critical security vulnerability" "closed" "" --from alice --yes --output json

### View the closed auction
swechaind query issuemarket get-auction 0 --output json


### Make a Transaction 
swechaind tx bank send alice $BOB 4000stake --from alice --yes

### Check Final Balances 
swechaind query bank balances alice --output json
swechaind query bank balances bob --output json
