#!/bin/bash
# =============================================================================
# Sub2API + LiteLLM 统一代理启动脚本
# =============================================================================

set -e

CADDY_CONFIG="/home/qinshu/sub2api/deploy/Caddyfile.unified"
LITELLM_CONFIG="/home/qinshu/litellm/litellm_config.yaml"
LITELLM_VENV="/home/qinshu/litellm/venv"
SUB2API_DIR="/home/qinshu/sub2api/backend"
SUDO_PASS="1"

# 颜色输出
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m'

info() { echo -e "${GREEN}[INFO]${NC} $1"; }
warn() { echo -e "${YELLOW}[WARN]${NC} $1"; }
error() { echo -e "${RED}[ERROR]${NC} $1"; exit 1; }

# 1. 检查并启动 PostgreSQL
check_postgres() {
    info "检查 PostgreSQL..."
    if ! pg_isready -h localhost -p 5432 > /dev/null 2>&1; then
        warn "PostgreSQL 未运行，正在启动..."
        echo "$SUDO_PASS" | sudo -S service postgresql start
        sleep 2
    fi
    info "PostgreSQL 已运行"
}

# 2. 检查并启动 Redis
check_redis() {
    info "检查 Redis..."
    if ! redis-cli -h localhost -p 6379 ping > /dev/null 2>&1; then
        warn "Redis 未运行，正在启动..."
        echo "$SUDO_PASS" | sudo -S service redis-server start
        sleep 1
    fi
    info "Redis 已运行"
}

# 3. 配置 DNS (WSL2 修复)
fix_dns() {
    info "检查 DNS..."
    if ! host api.anthropic.com > /dev/null 2>&1; then
        warn "DNS 异常，正在修复..."
        echo "$SUDO_PASS" | sudo -S bash -c 'echo -e "nameserver 8.8.8.8\nnameserver 114.114.114.114" > /etc/resolv.conf'
    fi
    info "DNS 正常"
}

# 4. 启动 Sub2API
start_sub2api() {
    info "检查 Sub2API..."
    if lsof -i :8080 > /dev/null 2>&1; then
        info "Sub2API 已在运行 (端口 8080)"
        return
    fi

    info "启动 Sub2API..."
    cd "$SUB2API_DIR"
    nohup ./sub2api > /tmp/sub2api.log 2>&1 &
    sleep 2

    if lsof -i :8080 > /dev/null 2>&1; then
        info "Sub2API 启动成功 (端口 8080)"
    else
        error "Sub2API 启动失败，查看日志: /tmp/sub2api.log"
    fi
}

# 5. 启动 LiteLLM
start_litellm() {
    info "检查 LiteLLM..."
    if lsof -i :4000 > /dev/null 2>&1; then
        info "LiteLLM 已在运行 (端口 4000)"
        return
    fi

    info "启动 LiteLLM..."
    cd /home/qinshu/litellm
    source "$LITELLM_VENV/bin/activate"
    export NO_PROXY="127.0.0.1,localhost"
    nohup litellm --config "$LITELLM_CONFIG" --port 4000 > /tmp/litellm.log 2>&1 &
    sleep 3

    if lsof -i :4000 > /dev/null 2>&1; then
        info "LiteLLM 启动成功 (端口 4000)"
    else
        warn "LiteLLM 启动失败，查看日志: /tmp/litellm.log"
    fi
}

# 6. 启动 Caddy 反向代理
start_caddy() {
    info "检查 Caddy..."
    if lsof -i :8000 > /dev/null 2>&1; then
        info "Caddy 已在运行 (端口 8000)"
        return
    fi

    info "启动 Caddy 反向代理..."
    nohup caddy run --config "$CADDY_CONFIG" --adapter caddyfile > /tmp/caddy.log 2>&1 &
    sleep 2

    if lsof -i :8000 > /dev/null 2>&1; then
        info "Caddy 启动成功 (端口 8000)"
    else
        error "Caddy 启动失败，查看日志: /tmp/caddy.log"
    fi
}

# 7. 配置 Windows 端口转发
setup_windows_port_forward() {
    info "配置 Windows 端口转发..."

    WSL_IP=$(hostname -I | awk '{print $1}')

    # 清理旧规则并添加新规则
    powershell.exe -Command "netsh interface portproxy delete v4tov4 listenport=8000 listenaddress=0.0.0.0" 2>/dev/null || true
    powershell.exe -Command "netsh interface portproxy add v4tov4 listenport=8000 listenaddress=0.0.0.0 connectport=8000 connectaddress=$WSL_IP" 2>/dev/null

    info "Windows 端口转发配置完成 (WSL IP: $WSL_IP)"
}

# 打印访问信息
print_info() {
    WIN_IP=$(powershell.exe -Command "(Get-NetIPAddress -AddressFamily IPv4 | Where-Object { \$_.InterfaceAlias -notlike '*Loopback*' -and \$_.InterfaceAlias -notlike '*vEthernet*' -and \$_.IPAddress -notlike '172.*' -and \$_.IPAddress -notlike '169.*' } | Select-Object -First 1).IPAddress" 2>/dev/null | tr -d '\r')

    echo ""
    echo "=========================================="
    echo "  统一代理启动成功!"
    echo "=========================================="
    echo ""
    echo "统一入口 (端口 8000):"
    echo "  - 本地: http://localhost:8000"
    [ -n "$WIN_IP" ] && echo "  - 局域网: http://${WIN_IP}:8000"
    echo ""
    echo "API 端点:"
    echo "  Claude:    /v1/messages"
    echo "  OpenAI:    /v1/chat/completions"
    echo "  Gemini:    /v1beta/"
    echo "  Responses: /v1/responses"
    echo ""
    echo "管理界面: http://localhost:8000/"
    echo ""
    echo "=========================================="
}

# 主流程
main() {
    echo "=========================================="
    echo "  Sub2API + LiteLLM 统一代理启动"
    echo "=========================================="

    check_postgres
    check_redis
    fix_dns
    start_sub2api
    start_litellm
    start_caddy
    setup_windows_port_forward
    print_info
}

main "$@"
