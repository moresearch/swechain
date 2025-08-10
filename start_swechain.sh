
#!/bin/sh
# Using /bin/sh for maximum compatibility

set -e

# ===========
#   CONFIG
# ===========
CHAIN_ID="swechain"
KEYRING="test"
DENOM="token"
CHAIN_DIR="$HOME/.swechain"
BINARY="swechaind"
LOG_FILE="$HOME/swechaind.log"
START_TIMEOUT=90
FIRST_BLOCK_WAIT=10
FAUCET_MNEMONIC="attitude motion repair drive edge chapter cave radar genius vault unique diesel scissors eagle matter pudding boring rose income cake target curve chunk such"

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# ===========
#   ARGUMENTS
# ===========
if [ "$1" = "balances" ]; then
    MODE="balances"
else
    if [ -z "$1" ]; then
        echo -e "${RED}❌ Usage:${NC} $0 <number_of_agents> | balances"
        echo -e "Example: ${YELLOW}$0 20${NC}"
        exit 1
    fi
    NUM_AGENTS="$1"
fi

# ===========
#   UTILS
# ===========
log_section() {
    echo -e "\n${YELLOW}========== $1 ==========${NC}"
}

check_dependencies() {
    log_section "Checking dependencies"
    for dep in jq lsof; do
        if ! command -v "$dep" >/dev/null 2>&1; then
            echo -e "${RED}❌ Missing dependency: $dep${NC}"
            exit 1
        fi
    done
    echo -e "${GREEN}✅ All dependencies are installed${NC}"
}

cleanup() {
    log_section "Cleaning up previous run"
    pkill -f "$BINARY" || true
    rm -rf "$CHAIN_DIR" "$LOG_FILE" ~/swechain_project
    echo -e "${GREEN}✅ Cleanup complete${NC}"
}

init_chain() {
    log_section "Initializing chain"
    $BINARY init mynode --chain-id "$CHAIN_ID"
    sed -i "s/\"stake\"/\"$DENOM\"/g" "$CHAIN_DIR/config/genesis.json"
    sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001token"/g' "$CHAIN_DIR/config/app.toml"
    echo -e "${GREEN}✅ Chain initialized${NC}"
}

setup_accounts() {
    log_section "Setting up accounts"

    # Validator & faucet
    $BINARY keys add validator --keyring-backend "$KEYRING"
    echo "$FAUCET_MNEMONIC" | $BINARY keys add faucet --keyring-backend "$KEYRING" --recover

    # Agents
    for i in $(seq 1 "$NUM_AGENTS"); do
        agent_name=$(printf "agent_%03d" "$i")
        $BINARY keys add "$agent_name" --keyring-backend "$KEYRING"
    done

    # Addresses
    VALIDATOR_ADDR=$($BINARY keys show validator --keyring-backend "$KEYRING" -a)
    FAUCET_ADDR=$($BINARY keys show faucet --keyring-backend "$KEYRING" -a)

    # Fund accounts
    $BINARY genesis add-genesis-account "$VALIDATOR_ADDR" 1000000000000$DENOM
    $BINARY genesis add-genesis-account "$FAUCET_ADDR" 1000000000000$DENOM
    for i in $(seq 1 "$NUM_AGENTS"); do
        agent_name=$(printf "agent_%03d" "$i")
        agent_addr=$($BINARY keys show "$agent_name" --keyring-backend "$KEYRING" -a)
        $BINARY genesis add-genesis-account "$agent_addr" 100000000$DENOM
    done

    # Validator gentx
    $BINARY genesis gentx validator 1000000000$DENOM --chain-id "$CHAIN_ID" --keyring-backend "$KEYRING"
    $BINARY genesis collect-gentxs
    $BINARY genesis validate-genesis
    echo -e "${GREEN}✅ Accounts created and funded${NC}"
}

start_blockchain() {
    log_section "Starting blockchain"

    for PORT in 26657 1317 9090; do
        if lsof -i :$PORT >/dev/null; then
            echo -e "${RED}❌ Port $PORT is in use${NC}"
            exit 1
        fi
    done

    $BINARY start > "$LOG_FILE" 2>&1 &
    CHAIN_PID=$!
    echo -e "Blockchain started with PID ${GREEN}$CHAIN_PID${NC}"
    wait_for_chain
}

wait_for_chain() {
    log_section "Waiting for chain readiness"
    start_time=$(date +%s)
    timeout=$((start_time + START_TIMEOUT))

    while true; do
        if [ "$(date +%s)" -ge "$timeout" ]; then
            echo -e "${RED}❌ Timeout waiting for chain${NC}"
            show_logs
            exit 1
        fi

        if ! ps -p $CHAIN_PID >/dev/null; then
            echo -e "${RED}❌ Blockchain process died${NC}"
            show_logs
            exit 1
        fi

        if curl -s http://localhost:26657/status | jq -e '.result.sync_info.catching_up == false' >/dev/null 2>&1; then
            echo -e "${GREEN}✅ Chain is ready${NC}"
            break
        fi
        sleep 1
        echo -n "."
    done

    echo -e "\n${YELLOW}⏳ Waiting $FIRST_BLOCK_WAIT seconds for first block...${NC}"
    sleep "$FIRST_BLOCK_WAIT"

    LATEST_HEIGHT=$(curl -s http://localhost:26657/status | jq -r '.result.sync_info.latest_block_height')
    if [ "$LATEST_HEIGHT" -gt 0 ]; then
        echo -e "${GREEN}✅ Block ${LATEST_HEIGHT} produced${NC}"
    else
        echo -e "${YELLOW}⚠ No blocks yet, continuing...${NC}"
    fi
}

show_logs() {
    log_section "Last 20 log lines"
    tail -n 20 "$LOG_FILE"
}

safe_query() {
    local cmd="$1"
    local retry=0
    while [ $retry -lt 3 ]; do
        if eval "$cmd" >/dev/null 2>&1; then
            eval "$cmd"
            return
        fi
        retry=$((retry+1))
        echo -e "${YELLOW}⚠ Retry $retry/3 in 5s...${NC}"
        sleep 5
    done
    echo -e "${RED}❌ Query failed after retries${NC}"
}

list_accounts() {
    log_section "Account addresses"
    echo "Validator: $($BINARY keys show validator --keyring-backend $KEYRING -a)"
    echo "Faucet: $($BINARY keys show faucet --keyring-backend $KEYRING -a)"
    for i in $(seq 1 "$NUM_AGENTS"); do
        agent_name=$(printf "agent_%03d" "$i")
        echo "$agent_name: $($BINARY keys show "$agent_name" --keyring-backend $KEYRING -a)"
    done
}

check_balances() {
    log_section "Checking balances"
    echo "Validator:"
    safe_query "$BINARY query bank balances $($BINARY keys show validator --keyring-backend $KEYRING -a) --output json | jq"
    echo ""
    echo "Faucet:"
    safe_query "$BINARY query bank balances $($BINARY keys show faucet --keyring-backend $KEYRING -a) --output json | jq"
    for i in $(seq 1 "$NUM_AGENTS"); do
        agent_name=$(printf "agent_%03d" "$i")
        echo ""
        echo "$agent_name:"
        safe_query "$BINARY query bank balances $($BINARY keys show "$agent_name" --keyring-backend $KEYRING -a) --output json | jq"
    done
}

get_all_balances() {
    log_section "All Account Balances"
    accounts=$($BINARY keys list --keyring-backend "$KEYRING" --output json | jq -r '.[] | .name')
    for account in $accounts; do
        echo -e "${GREEN}Account: $account${NC}"
        safe_query "$BINARY query bank balances $($BINARY keys show $account -a --keyring-backend $KEYRING) --output json | jq '.balances'"
        echo -e "${YELLOW}------------------------${NC}"
    done
}

main() {
    echo -e "${GREEN}🚀 Starting setup at: $(date)${NC}"
    check_dependencies
    cleanup
    init_chain
    setup_accounts
    start_blockchain
    list_accounts
    check_balances
    echo -e "\n${GREEN}✅ Setup completed successfully!${NC}"
    echo -e "📋 Logs: ${YELLOW}tail -f $LOG_FILE${NC}"
}

if [ "$MODE" = "balances" ]; then
    get_all_balances
else
    main
fi
