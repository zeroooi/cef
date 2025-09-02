#!/bin/bash

# CEF 进程清理脚本
# 解决 macOS 上 CEF 应用的进程冲突问题

echo "=== CEF 进程清理脚本 ==="

# 查找并终止相关进程
echo "正在查找 CEF 相关进程..."

# 终止我们的应用进程
pkill -f "cef-app" 2>/dev/null || true
pkill -f "main.go" 2>/dev/null || true

# 终止可能的 CEF helper 进程
pkill -f "CEF Helper" 2>/dev/null || true
pkill -f "Chromium Helper" 2>/dev/null || true

# 终止可能的 Energy 进程
pkill -f "energy" 2>/dev/null || true

# 清理缓存和用户数据目录
echo "正在清理缓存和临时文件..."
rm -rf ./cache 2>/dev/null || true
rm -rf ./userdata 2>/dev/null || true
rm -rf ./cef_debug.log 2>/dev/null || true

# 等待进程完全退出
sleep 2

echo "清理完成！现在可以安全启动 CEF 应用了。"
echo ""
echo "使用以下命令启动应用："
echo "go run main.go env=dev"