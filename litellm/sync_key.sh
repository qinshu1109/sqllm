#!/bin/bash
# LiteLLM API Key 同步脚本
# 用法: ./sync_key.sh <api_key>
# 功能: 更新 LiteLLM 配置中的 api_key 和 master_key，然后重启服务

set -e

API_KEY=$1
CONFIG_FILE="/home/qinshu/litellm/litellm_config.yaml"
LOG_FILE="/tmp/litellm.log"

if [ -z "$API_KEY" ]; then
    echo "错误: 请提供 API Key"
    echo "用法: $0 <api_key>"
    exit 1
fi

# 验证 Key 格式
if [[ ! "$API_KEY" =~ ^sk-[a-f0-9]{64}$ ]]; then
    echo "错误: API Key 格式无效，应为 sk- 开头的 64 位十六进制字符串"
    exit 1
fi

echo "开始同步 API Key 到 LiteLLM..."

# 备份配置文件
cp "$CONFIG_FILE" "${CONFIG_FILE}.bak"
echo "已备份配置文件到 ${CONFIG_FILE}.bak"

# 更新所有模型的 api_key
sed -i "s/api_key: sk-[a-f0-9]\{64\}/api_key: $API_KEY/g" "$CONFIG_FILE"
echo "已更新模型 api_key"

# 更新 master_key
sed -i "s/master_key: sk-[a-f0-9]\{64\}/master_key: $API_KEY/g" "$CONFIG_FILE"
echo "已更新 master_key"

# 停止现有 LiteLLM 进程（只杀 litellm 服务，不杀脚本自身）
echo "正在停止 LiteLLM..."
pkill -f "litellm --config" 2>/dev/null || true
sleep 2

# 重启 LiteLLM
echo "正在启动 LiteLLM..."
cd /home/qinshu/litellm
NO_PROXY="127.0.0.1,localhost" PYTHONPATH=/home/qinshu/litellm \
    nohup /home/qinshu/litellm/venv/bin/litellm --config litellm_config.yaml --host 0.0.0.0 --port 4000 > "$LOG_FILE" 2>&1 &

# 等待服务启动
sleep 5

# 检查服务是否启动成功
if pgrep -f "litellm" > /dev/null; then
    echo "LiteLLM 已成功启动"
    echo "新 API Key: ${API_KEY:0:20}..."
    exit 0
else
    echo "错误: LiteLLM 启动失败，请检查日志: $LOG_FILE"
    # 恢复备份
    cp "${CONFIG_FILE}.bak" "$CONFIG_FILE"
    echo "已恢复配置文件备份"
    exit 1
fi
