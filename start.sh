#!/bin/bash

# 定义颜色代码
RED='\033[0;31m'
GREEN='\033[0;32m'
NC='\033[0m' # No Color

# 输出带颜色的消息函数
log() {
    echo -e "${GREEN}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

error() {
    echo -e "${RED}[$(date +'%Y-%m-%d %H:%M:%S')] $1${NC}"
}

# 检查前端目录是否存在
if [ ! -d "frontend" ]; then
    error "前端目录不存在"
    exit 1
fi

# 检查后端目录是否存在
if [ ! -d "backend" ]; then
    error "后端目录不存在"
    exit 1
fi

# 启动后端服务
log "正在启动后端服务..."
cd backend
go run main.go & 
BACKEND_PID=$!

# 返回项目根目录
cd ..

# 启动前端服务
log "正在启动前端服务..."
cd frontend
npm run dev &
FRONTEND_PID=$!

# 等待用户按Ctrl+C
log "服务已启动! 按Ctrl+C停止所有服务"
wait $BACKEND_PID $FRONTEND_PID