
#!/bin/sh
# Using /bin/sh for maximum compatibility

# Set options
set -e

# Support more Agents? better different ports but this should work for now!
# 1. Export the variable
#export OLLAMA_NUM_PARALLEL=4
# 2. Stop any existing Ollama instance (if running)
#pkill -f "ollama serve"
# 3. Restart Ollama with the env var applied
#nohup ollama serve > ollama.log 2>&1 &
#ollama serve

# Chain configuration
CHAIN_ID="swechain"
KEYRING="test"
DENOM="token"
CHAIN_DIR="$HOME/.swechain"
BINARY="swechaind"
LOG_FILE="$HOME/swechaind.log"
START_TIMEOUT=90
FIRST_BLOCK_WAIT=10  # Wait time for first block in seconds
FAUCET_MNEMONIC="attitude motion repair drive edge chapter cave radar genius vault unique diesel scissors eagle matter pudding boring rose income cake target curve chunk such"

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check dependencies
check_dependencies() {
  if ! command -v jq > /dev/null 2>&1; then
    echo -e "${RED}❌ Missing dependency: jq${NC}"
    echo "Please install with: sudo apt-get install -y jq"
    exit 1
  fi
  
  if ! command -v lsof > /dev/null 2>&1; then
    echo -e "${RED}❌ Missing dependency: lsof${NC}"
    echo "Please install with: sudo apt-get install -y lsof"
    exit 1
  fi
}

# Clean up previous run
cleanup() {
  echo -e "${YELLOW}⚠ Cleaning up previous run...${NC}"
  pkill -f "$BINARY" || true
  rm -rf $CHAIN_DIR
  rm -f "$LOG_FILE"
  
  # Also clean up any previous project directory
  rm -rf ~/swechain_project
}

# Initialize chain directly without ignite
init_chain() {
  echo -e "${YELLOW}🧪 Initializing chain...${NC}"
  
  # Initialize the chain with the binary
  $BINARY init mynode --chain-id $CHAIN_ID
  
  # Modify genesis to use our token denom instead of stake
  sed -i "s/\"stake\"/\"$DENOM\"/g" $CHAIN_DIR/config/genesis.json
  
  # Set minimum gas price in app.toml
  sed -i 's/minimum-gas-prices = ""/minimum-gas-prices = "0.001token"/g' $CHAIN_DIR/config/app.toml
}

# Setup accounts
setup_accounts() {
  echo -e "${YELLOW}🔑 Setting up accounts...${NC}"
  
  # Add validator account
  $BINARY keys add validator --keyring-backend $KEYRING
  
  # Add faucet account from mnemonic
  echo "$FAUCET_MNEMONIC" | $BINARY keys add faucet --keyring-backend $KEYRING --recover
  
  # Add agent accounts
  $BINARY keys add agent1 --keyring-backend $KEYRING
  $BINARY keys add agent2 --keyring-backend $KEYRING
  $BINARY keys add agent3 --keyring-backend $KEYRING
  $BINARY keys add agent4 --keyring-backend $KEYRING
  $BINARY keys add agent5 --keyring-backend $KEYRING
  $BINARY keys add agent6 --keyring-backend $KEYRING
  
  # Get addresses
  VALIDATOR_ADDR=$($BINARY keys show validator --keyring-backend $KEYRING -a)
  FAUCET_ADDR=$($BINARY keys show faucet --keyring-backend $KEYRING -a)
  AGENT1_ADDR=$($BINARY keys show agent1 --keyring-backend $KEYRING -a)
  AGENT2_ADDR=$($BINARY keys show agent2 --keyring-backend $KEYRING -a)
  AGENT3_ADDR=$($BINARY keys show agent3 --keyring-backend $KEYRING -a)
  AGENT4_ADDR=$($BINARY keys show agent4 --keyring-backend $KEYRING -a)
  AGENT5_ADDR=$($BINARY keys show agent5 --keyring-backend $KEYRING -a)
  AGENT6_ADDR=$($BINARY keys show agent6 --keyring-backend $KEYRING -a)
  
  # Add genesis accounts with sufficient funds
  $BINARY genesis add-genesis-account $VALIDATOR_ADDR 1000000000000$DENOM
  $BINARY genesis add-genesis-account $FAUCET_ADDR 1000000000000$DENOM
  $BINARY genesis add-genesis-account $AGENT1_ADDR 100000000$DENOM
 $BINARY genesis add-genesis-account $AGENT2_ADDR 100000000$DENOM
$BINARY genesis add-genesis-account $AGENT3_ADDR 100000000$DENOM
 $BINARY genesis add-genesis-account $AGENT4_ADDR 100000000$DENOM
  $BINARY genesis add-genesis-account $AGENT5_ADDR 100000000$DENOM
  $BINARY genesis add-genesis-account $AGENT6_ADDR 100000000$DENOM
  
  # Create validator gentx with sufficient delegation
  $BINARY genesis gentx validator 1000000000$DENOM --chain-id $CHAIN_ID --keyring-backend $KEYRING
  
  # Collect gentxs
  $BINARY genesis collect-gentxs
  
  # Validate genesis
  $BINARY genesis validate-genesis
}

# Start the blockchain
start_blockchain() {
  echo -e "${YELLOW}⚠ Checking ports...${NC}"
  for PORT in 26657 1317 9090; do
    if lsof -i :$PORT > /dev/null; then
      echo -e "${RED}❌ Port $PORT is in use${NC}"
      exit 1
    fi
  done

  echo -e "${YELLOW}🚀 Starting blockchain...${NC}"
  $BINARY start > "$LOG_FILE" 2>&1 &
  
  CHAIN_PID=$!
  echo -e "Blockchain started with PID ${GREEN}$CHAIN_PID${NC}"

  wait_for_chain
}

# Wait for chain to be ready
wait_for_chain() {
  echo -e "${YELLOW}⏳ Waiting for chain to be ready (timeout: ${START_TIMEOUT}s)...${NC}"
  start_time=$(date +%s)
  timeout=$((start_time + START_TIMEOUT))

  while true; do
    # Check timeout
    current_time=$(date +%s)
    if [ $current_time -ge $timeout ]; then
      echo -e "${RED}❌ Timeout reached while waiting for chain to start${NC}"
      show_logs
      exit 1
    fi

    # Check if process is still running
    if ! ps -p $CHAIN_PID > /dev/null; then
      echo -e "${RED}❌ Blockchain process died${NC}"
      show_logs
      exit 1
    fi

    # Check chain status
    if curl -s http://localhost:26657/status > /dev/null 2>&1; then
      if curl -s http://localhost:26657/status | jq -e '.result.sync_info.catching_up == false' >/dev/null 2>&1; then
        echo -e "${GREEN}✅ Chain is ready!${NC}"
        break
      fi
    fi

    sleep 1
    echo -n "."
  done

  # Wait for the first block to be produced
  echo -e "${YELLOW}⏳ Waiting for first block to be produced (${FIRST_BLOCK_WAIT}s)...${NC}"
  sleep $FIRST_BLOCK_WAIT
  
  # Verify that we have at least one block
  LATEST_HEIGHT=$(curl -s http://localhost:26657/status | jq -r '.result.sync_info.latest_block_height')
  if [ "$LATEST_HEIGHT" -gt 0 ]; then
    echo -e "${GREEN}✅ Block ${LATEST_HEIGHT} produced!${NC}"
  else
    echo -e "${YELLOW}⚠ No blocks produced yet, but continuing...${NC}"
  fi
}

# Show recent logs if something fails
show_logs() {
  echo -e "\n${YELLOW}=== Last 20 lines of log ===${NC}"
  tail -n 20 "$LOG_FILE"
  echo -e "${YELLOW}===========================${NC}"
}

# Safe query function with retry
safe_query() {
  local MAX_RETRIES=3
  local RETRY_WAIT=5
  local cmd="$1"
  local retry=0
  
  while [ $retry -lt $MAX_RETRIES ]; do
    if eval "$cmd" > /dev/null 2>&1; then
      eval "$cmd"
      return 0
    else
      retry=$((retry + 1))
      echo -e "${YELLOW}⚠ Query failed, retrying in ${RETRY_WAIT}s (${retry}/${MAX_RETRIES})...${NC}"
      sleep $RETRY_WAIT
    fi
  done
  
  echo -e "${RED}❌ Query failed after ${MAX_RETRIES} retries${NC}"
  return 1
}

# Check account balances with retry
check_balances() {
  echo -e "${YELLOW}📤 Checking balances...${NC}"
  
  VALIDATOR_ADDR=$($BINARY keys show validator --keyring-backend $KEYRING -a)
  FAUCET_ADDR=$($BINARY keys show faucet --keyring-backend $KEYRING -a)
  AGENT1_ADDR=$($BINARY keys show agent1 --keyring-backend $KEYRING -a)
  AGENT2_ADDR=$($BINARY keys show agent2 --keyring-backend $KEYRING -a)
  AGENT3_ADDR=$($BINARY keys show agent3 --keyring-backend $KEYRING -a)
  AGENT4_ADDR=$($BINARY keys show agent4 --keyring-backend $KEYRING -a)
  AGENT5_ADDR=$($BINARY keys show agent5 --keyring-backend $KEYRING -a)
  AGENT6_ADDR=$($BINARY keys show agent6 --keyring-backend $KEYRING -a)
  
  echo -e "Validator (${VALIDATOR_ADDR}):"
  safe_query "$BINARY query bank balances $VALIDATOR_ADDR --output json | jq"
  
  echo -e "\nFaucet (${FAUCET_ADDR}):"
  safe_query "$BINARY query bank balances $FAUCET_ADDR --output json | jq"
  
  echo -e "\nAgent1 (${AGENT1_ADDR}):"
  safe_query "$BINARY query bank balances $AGENT1_ADDR --output json | jq"
  
  echo -e "\nAgent2 (${AGENT2_ADDR}):"
  safe_query "$BINARY query bank balances $AGENT2_ADDR --output json | jq"
  
  echo -e "\nAgent3 (${AGENT3_ADDR}):"
  safe_query "$BINARY query bank balances $AGENT3_ADDR --output json | jq"
  
  echo -e "\nAgent4 (${AGENT4_ADDR}):"
  safe_query "$BINARY query bank balances $AGENT4_ADDR --output json | jq"
  
  echo -e "\nAgent5 (${AGENT5_ADDR}):"
  safe_query "$BINARY query bank balances $AGENT5_ADDR --output json | jq"
  
  echo -e "\nAgent6 (${AGENT6_ADDR}):"
  safe_query "$BINARY query bank balances $AGENT6_ADDR --output json | jq"
}

# Display account addresses
list_accounts() {
  echo -e "${YELLOW}🔑 Account addresses:${NC}"
  echo -e "Validator: $($BINARY keys show validator --keyring-backend $KEYRING -a)"
  echo -e "Faucet: $($BINARY keys show faucet --keyring-backend $KEYRING -a)"
  echo -e "Agent1: $($BINARY keys show agent1 --keyring-backend $KEYRING -a)"
  echo -e "Agent2: $($BINARY keys show agent2 --keyring-backend $KEYRING -a)"
  echo -e "Agent3: $($BINARY keys show agent3 --keyring-backend $KEYRING -a)"
  echo -e "Agent4: $($BINARY keys show agent4 --keyring-backend $KEYRING -a)"
  echo -e "Agent5: $($BINARY keys show agent5 --keyring-backend $KEYRING -a)"
  echo -e "Agent6: $($BINARY keys show agent6 --keyring-backend $KEYRING -a)"
}

# Add a function to get all balances
get_all_balances() {
  echo -e "${YELLOW}===== All Account Balances =====${NC}"
  accounts=$(swechaind keys list --keyring-backend test --output json | jq -r '.[] | .name')

  for account in $accounts; do
    address=$(swechaind keys show $account -a --keyring-backend test)
    echo -e "${GREEN}Account: $account${NC}"
    safe_query "$BINARY query bank balances $address --output json | jq '.balances'"
    echo -e "${YELLOW}------------------------${NC}"
  done
}

# Main execution
main() {
  echo -e "${GREEN}Starting setup at: $(date)${NC}"
  
  check_dependencies
  cleanup
  init_chain
  setup_accounts
  start_blockchain
  list_accounts  # List accounts first, which will always succeed
  check_balances  # Then check balances with retry mechanism
  
  echo -e "\n${GREEN}✅ Setup completed successfully!${NC}"
  echo -e "🔗 Chain is running with PID ${GREEN}$CHAIN_PID${NC}"
  echo -e "📋 View logs: ${YELLOW}tail -f $LOG_FILE${NC}"
  echo -e "🛑 To stop: ${YELLOW}kill $CHAIN_PID${NC}"
  echo -e "\n${GREEN}💰 To check all balances at any time:${NC}"
  echo -e "${YELLOW}./$(basename "$0") balances${NC}"
  echo -e "\n${GREEN}ℹ️ Remember to use --keyring-backend test and --chain-id $CHAIN_ID for your commands${NC}"
}

# Check if being run with "balances" argument
if [ "$1" = "balances" ]; then
  get_all_balances
  exit 0
fi

main