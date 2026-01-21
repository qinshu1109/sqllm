#!/bin/bash
# Start LiteLLM with NO_PROXY to bypass SOCKS proxy for local connections
cd /home/qinshu/litellm
source venv/bin/activate

# Exclude local addresses from proxy
export NO_PROXY="127.0.0.1,localhost"
export no_proxy="127.0.0.1,localhost"

# Start LiteLLM
litellm --config litellm_config.yaml --port 4000
