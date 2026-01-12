#!/bin/bash
# =============================================================================
# 停止统一代理服务
# =============================================================================

echo "停止 Caddy..."
pkill -f "caddy run" 2>/dev/null || true

echo "停止 LiteLLM..."
pkill -f "litellm" 2>/dev/null || true

echo "停止 Sub2API..."
pkill -f "sub2api" 2>/dev/null || true

echo ""
echo "清理 Windows 端口转发..."
powershell.exe -Command "netsh interface portproxy delete v4tov4 listenport=8000 listenaddress=0.0.0.0" 2>/dev/null || true

echo ""
echo "所有服务已停止"
