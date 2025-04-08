#!/bin/bash

# Configuration
COMPOSE_FILE="build/postgres/docker-compose.yml"
CONTAINER_NAME="starpivot_postgres"
TIMEOUT=30

# Colors
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check docker-compose availability
if ! command -v docker-compose &> /dev/null; then
  echo -e "${RED}Error: docker-compose not found${NC}"
  exit 1
fi

function start_service() {
  echo -e "${YELLOW}Starting PostgreSQL container...${NC}"
  docker-compose -f $COMPOSE_FILE up -d
  
  echo -e "${YELLOW}Waiting for health check (max ${TIMEOUT}s)...${NC}"
  for ((i=1; i<=$TIMEOUT; i++)); do
    status=$(docker inspect --format='{{.State.Health.Status}}' $CONTAINER_NAME 2>/dev/null)
    if [ "$status" == "healthy" ]; then
      echo -e "${GREEN}PostgreSQL is ready${NC}"
      return 0
    fi
    sleep 1
  done
  
  echo -e "${RED}Timeout: PostgreSQL did not become healthy${NC}"
  return 1
}

function stop_service() {
  echo -e "${YELLOW}Stopping PostgreSQL container...${NC}"
  docker-compose -f $COMPOSE_FILE down -v
}

case "$1" in
  run)
    start_service
    ;;
  stop)
    stop_service
    ;;
  *)
    echo "Usage: $0 {run|stop}"
    exit 1
    ;;
esac

exit $?