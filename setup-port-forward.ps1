# Sub2API WSL2 端口转发配置脚本
# 以管理员身份运行此 PowerShell 脚本

# 获取 WSL2 IP 地址
$wslIp = (wsl hostname -I).Trim().Split()[0]
$port = 8080

Write-Host "WSL IP: $wslIp" -ForegroundColor Green
Write-Host "端口: $port" -ForegroundColor Green

# 删除旧的端口转发规则（如果存在）
netsh interface portproxy delete v4tov4 listenport=$port listenaddress=0.0.0.0 2>$null

# 添加新的端口转发规则
Write-Host "配置端口转发..." -ForegroundColor Yellow
netsh interface portproxy add v4tov4 listenport=$port listenaddress=0.0.0.0 connectport=$port connectaddress=$wslIp

# 配置防火墙规则
Write-Host "配置防火墙..." -ForegroundColor Yellow
netsh advfirewall firewall delete rule name="WSL2 Sub2API" 2>$null
netsh advfirewall firewall add rule name="WSL2 Sub2API" dir=in action=allow protocol=TCP localport=$port

# 显示配置结果
Write-Host "`n端口转发配置完成！" -ForegroundColor Green
netsh interface portproxy show all

Write-Host "`n你现在可以通过以下方式访问:" -ForegroundColor Cyan
Write-Host "- 本地访问: http://localhost:$port" -ForegroundColor White
Write-Host "- 局域网访问: http://你的Windows主机IP:$port" -ForegroundColor White
Write-Host "`n获取 Windows IP 命令: ipconfig" -ForegroundColor Gray
