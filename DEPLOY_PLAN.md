# Sub2API 同步部署计划（保持 sqllm.dpdns.org 域名）

## 问题原因
上次部署后提示的访问地址是 IP+端口，导致用户用 IP 访问。实际上域名配置是正常的，只需要用域名访问即可。

## 服务器域名架构
```
用户 → https://sqllm.dpdns.org → Cloudflare → AWS:443 → Caddy → localhost:8080 → sub2api
```

Caddy 配置正确：
```caddyfile
sqllm.dpdns.org {
    reverse_proxy localhost:8080
}
```

## 部署步骤

### 1. 上传新版本到服务器
```bash
# 备份
ssh -i /tmp/qinshu.pem ubuntu@44.218.81.97 "cp /home/ubuntu/sub2api/sub2api /home/ubuntu/sub2api/sub2api.backup-$(date +%Y%m%d-%H%M%S)"

# 上传
scp -i /tmp/qinshu.pem /home/qinshu/sub2api/backend/sub2api ubuntu@44.218.81.97:/home/ubuntu/sub2api/sub2api.new

# 替换并重启
ssh -i /tmp/qinshu.pem ubuntu@44.218.81.97 "sudo systemctl stop sub2api && mv /home/ubuntu/sub2api/sub2api.new /home/ubuntu/sub2api/sub2api && chmod +x /home/ubuntu/sub2api/sub2api && sudo systemctl start sub2api"
```

### 2. 验证部署
```bash
# 检查健康状态
curl -s https://sqllm.dpdns.org/health
```

### 3. 访问地址
**请使用域名访问（不是 IP）**：
- 网页登录: https://sqllm.dpdns.org
- Claude API: https://sqllm.dpdns.org/v1/messages
- Gemini API: https://sqllm.dpdns.org/v1beta/

## 本地版本包含的功能
1. 上游 v0.1.45 新功能：Claude Code 客户端限制、回退分组
2. 本地自定义：model_rates 分组倍率、API 文档页面（含 sqllm.dpdns.org 示例）
