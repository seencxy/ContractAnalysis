#!/bin/bash

# ============================================
# Futures Analysis - Docker Deployment Script
# ============================================

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Script directory
SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_ROOT="$(cd "${SCRIPT_DIR}/../.." && pwd)"

# Default values
COMPOSE_FILE="${PROJECT_ROOT}/docker-compose.yml"
ENV_FILE="${PROJECT_ROOT}/.env.docker"
ACTION="up"

# Functions
log_info() {
    echo -e "${BLUE}[INFO]${NC} $1"
}

log_success() {
    echo -e "${GREEN}[SUCCESS]${NC} $1"
}

log_warning() {
    echo -e "${YELLOW}[WARNING]${NC} $1"
}

log_error() {
    echo -e "${RED}[ERROR]${NC} $1"
}

print_banner() {
    echo -e "${BLUE}"
    echo "╔══════════════════════════════════════════════════╗"
    echo "║     Binance Futures Analysis - Docker Deploy     ║"
    echo "╚══════════════════════════════════════════════════╝"
    echo -e "${NC}"
}

print_usage() {
    echo "Usage: $0 [COMMAND] [OPTIONS]"
    echo ""
    echo "Commands:"
    echo "  up          Start all services (default)"
    echo "  down        Stop all services"
    echo "  restart     Restart all services"
    echo "  build       Build/rebuild containers"
    echo "  logs        Show logs (use -f to follow)"
    echo "  status      Show service status"
    echo "  clean       Stop and remove all data (DESTRUCTIVE)"
    echo "  init        Initialize environment file"
    echo ""
    echo "Options:"
    echo "  -f, --follow    Follow log output (for 'logs' command)"
    echo "  -h, --help      Show this help message"
    echo ""
    echo "Examples:"
    echo "  $0 up           # Start all services"
    echo "  $0 logs -f      # Follow logs"
    echo "  $0 down         # Stop services"
}

check_docker() {
    if ! command -v docker &> /dev/null; then
        log_error "Docker is not installed. Please install Docker first."
        exit 1
    fi

    if ! command -v docker-compose &> /dev/null && ! docker compose version &> /dev/null; then
        log_error "Docker Compose is not installed. Please install Docker Compose first."
        exit 1
    fi
}

get_compose_cmd() {
    if docker compose version &> /dev/null 2>&1; then
        echo "docker compose"
    else
        echo "docker-compose"
    fi
}

init_env() {
    log_info "Initializing environment file..."
    
    if [ -f "${ENV_FILE}" ]; then
        log_warning "Environment file already exists: ${ENV_FILE}"
        read -p "Do you want to overwrite it? (y/N): " confirm
        if [ "$confirm" != "y" ] && [ "$confirm" != "Y" ]; then
            log_info "Skipping environment file initialization"
            return
        fi
    fi

    cp "${PROJECT_ROOT}/.env.docker.example" "${ENV_FILE}"
    log_success "Environment file created: ${ENV_FILE}"
    log_warning "Please edit ${ENV_FILE} and update the passwords before deployment!"
}

do_up() {
    log_info "Starting services..."
    
    # Check environment file
    if [ ! -f "${ENV_FILE}" ]; then
        log_warning "Environment file not found. Creating from example..."
        init_env
    fi

    # Create logs directory
    mkdir -p "${PROJECT_ROOT}/logs"

    # Build and start containers
    COMPOSE_CMD=$(get_compose_cmd)
    cd "${PROJECT_ROOT}"
    $COMPOSE_CMD --env-file "${ENV_FILE}" -f "${COMPOSE_FILE}" up -d --build
    
    log_success "Services started successfully!"
    echo ""
    log_info "Checking service status..."
    sleep 5
    do_status
}

do_down() {
    log_info "Stopping services..."
    COMPOSE_CMD=$(get_compose_cmd)
    cd "${PROJECT_ROOT}"
    $COMPOSE_CMD -f "${COMPOSE_FILE}" down
    log_success "Services stopped"
}

do_restart() {
    log_info "Restarting services..."
    do_down
    do_up
}

do_build() {
    log_info "Building containers..."
    COMPOSE_CMD=$(get_compose_cmd)
    cd "${PROJECT_ROOT}"
    $COMPOSE_CMD --env-file "${ENV_FILE}" -f "${COMPOSE_FILE}" build --no-cache
    log_success "Build complete"
}

do_logs() {
    COMPOSE_CMD=$(get_compose_cmd)
    cd "${PROJECT_ROOT}"
    
    if [ "$FOLLOW_LOGS" = true ]; then
        $COMPOSE_CMD -f "${COMPOSE_FILE}" logs -f
    else
        $COMPOSE_CMD -f "${COMPOSE_FILE}" logs --tail=100
    fi
}

do_status() {
    COMPOSE_CMD=$(get_compose_cmd)
    cd "${PROJECT_ROOT}"
    
    echo ""
    echo "Container Status:"
    echo "─────────────────────────────────────────────────"
    $COMPOSE_CMD -f "${COMPOSE_FILE}" ps
    echo ""
    
    # Check service health
    echo "Service Health:"
    echo "─────────────────────────────────────────────────"
    
    # MySQL
    if docker ps --filter "name=futures-mysql" --filter "status=running" | grep -q futures-mysql; then
        echo -e "MySQL:  ${GREEN}●${NC} Running"
    else
        echo -e "MySQL:  ${RED}●${NC} Not Running"
    fi
    
    # Redis
    if docker ps --filter "name=futures-redis" --filter "status=running" | grep -q futures-redis; then
        echo -e "Redis:  ${GREEN}●${NC} Running"
    else
        echo -e "Redis:  ${RED}●${NC} Not Running"
    fi
    
    # App
    if docker ps --filter "name=futures-app" --filter "status=running" | grep -q futures-app; then
        echo -e "App:    ${GREEN}●${NC} Running"
    else
        echo -e "App:    ${RED}●${NC} Not Running"
    fi
    
    echo ""
    echo "Access Points:"
    echo "─────────────────────────────────────────────────"
    echo "  API:       http://localhost:${APP_PORT:-8080}"
    echo "  Metrics:   http://localhost:${METRICS_PORT:-9090}/metrics"
    echo "  Health:    http://localhost:${APP_PORT:-8080}/health"
    echo ""
}

do_clean() {
    log_warning "This will stop all containers and DELETE ALL DATA!"
    read -p "Are you sure? (type 'yes' to confirm): " confirm
    
    if [ "$confirm" != "yes" ]; then
        log_info "Cancelled"
        return
    fi
    
    log_info "Stopping and cleaning up..."
    COMPOSE_CMD=$(get_compose_cmd)
    cd "${PROJECT_ROOT}"
    $COMPOSE_CMD -f "${COMPOSE_FILE}" down -v --remove-orphans
    
    log_success "Cleanup complete"
}

# Main execution
print_banner
check_docker

FOLLOW_LOGS=false

# Parse arguments
while [[ $# -gt 0 ]]; do
    case $1 in
        up|start)
            ACTION="up"
            shift
            ;;
        down|stop)
            ACTION="down"
            shift
            ;;
        restart)
            ACTION="restart"
            shift
            ;;
        build)
            ACTION="build"
            shift
            ;;
        logs)
            ACTION="logs"
            shift
            ;;
        status|ps)
            ACTION="status"
            shift
            ;;
        clean)
            ACTION="clean"
            shift
            ;;
        init)
            ACTION="init"
            shift
            ;;
        -f|--follow)
            FOLLOW_LOGS=true
            shift
            ;;
        -h|--help)
            print_usage
            exit 0
            ;;
        *)
            log_error "Unknown option: $1"
            print_usage
            exit 1
            ;;
    esac
done

# Execute action
case $ACTION in
    up)
        do_up
        ;;
    down)
        do_down
        ;;
    restart)
        do_restart
        ;;
    build)
        do_build
        ;;
    logs)
        do_logs
        ;;
    status)
        do_status
        ;;
    clean)
        do_clean
        ;;
    init)
        init_env
        ;;
esac
