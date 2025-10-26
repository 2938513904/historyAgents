# AI多智能体系统 - 部署指南

一个基于Web的多智能体聊天室应用，支持创建多个AI智能体角色，并在聊天室中进行自动讨论。

## 项目概述

本项目是一个现代化的多智能体聊天室系统，允许用户创建具有不同角色和性格的AI智能体，并在主题聊天室中进行自动化讨论。系统通过调用阿里云DashScope API实现真实的AI对话功能，并通过WebSocket提供实时消息推送。

## 功能特性

### 🎯 核心功能
- **智能体管理**: 创建和管理多个AI智能体角色，每个智能体都有独特的名称、角色定位和性格描述
- **聊天室系统**: 创建主题聊天室，选择参与的智能体，启动自动化讨论
- **实时通信**: 使用WebSocket实现实时消息推送，观察智能体的讨论过程

### 🎨 界面特性
- **现代化UI**: 采用渐变背景和毛玻璃效果，界面美观现代
- **响应式设计**: 支持桌面端和移动端访问
- **实时更新**: 聊天消息实时显示，支持讨论状态管理
- **用户友好**: 直观的操作界面，易于使用

## 技术架构

### 后端技术
- **语言**: Go 1.21+
- **Web框架**: Gin Web Framework
- **实时通信**: Gorilla WebSocket
- **数据存储**: MySQL数据库
- **API集成**: 阿里云DashScope API
- **配置管理**: YAML配置文件支持

### 前端技术
- **HTML5 + CSS3**: 现代化界面设计
- **原生JavaScript**: 无需额外框架依赖
- **WebSocket**: 实时通信支持

## 系统架构

- **后端**: Go + Gin + GORM + MySQL
- **前端**: HTML + CSS + JavaScript
- **数据库**: MySQL 8.0
- **部署**: Docker + Docker Compose

## 目录结构

```
AIagents2.0/
├── backend/                 # 后端代码
│   ├── internal/           # 内部包
│   │   ├── config/        # 配置管理
│   │   ├── database/      # 数据库连接
│   │   ├── model/         # 数据模型
│   │   └── router/        # 路由处理
│   ├── main.go            # 主程序入口
│   ├── config.yaml        # 配置文件
│   ├── go.mod             # Go模块文件
│   └── go.sum             # 依赖校验文件
├── frontend/               # 前端代码
│   ├── index.html         # 主页面
│   ├── style.css          # 样式文件
│   └── app.js             # JavaScript逻辑
├── nginx/                  # Nginx配置
│   └── nginx.conf         # Nginx配置文件
├── mysql/                  # MySQL初始化
│   └── init/              # 初始化脚本
│       └── 01-init.sql    # 数据库初始化
├── Dockerfile             # Docker构建文件
├── docker-compose.yml     # Docker Compose配置
├── .env.example           # 环境变量示例
└── README.md              # 本文档
```

## 配置说明

### 1. YAML配置文件 (backend/config.yaml)

系统支持通过YAML文件进行配置，主要配置项包括：

```yaml
# 服务器配置
server:
  port: 8080          # 后端服务端口
  host: "0.0.0.0"     # 监听地址

# 数据库配置
database:
  type: "mysql"       # 数据库类型
  host: "localhost"   # 数据库主机
  port: 3306          # 数据库端口
  username: "root"    # 数据库用户名
  password: "root"    # 数据库密码
  dbname: "ai_agents" # 数据库名称
  charset: "utf8mb4"  # 字符集
  timezone: "Local"   # 时区
  max_idle_conns: 10  # 最大空闲连接数
  max_open_conns: 100 # 最大打开连接数

# API配置
api:
  api_key: "your-api-key"     # DashScope API Key
  base_url: "https://..."     # API基础URL
  model: "qwen-plus"          # 默认模型
  endpoint: "https://..."     # API端点

# 日志配置
log:
  level: "info"               # 日志级别
  file: "logs/app.log"        # 日志文件路径
  console: true               # 是否输出到控制台
```

### 2. 环境变量配置 (.env)

Docker Compose支持通过环境变量配置：

```bash
# MySQL数据库配置
MYSQL_ROOT_PASSWORD=your_root_password
MYSQL_DATABASE=ai_agents
MYSQL_USER=aiuser
MYSQL_PASSWORD=your_mysql_password
MYSQL_PORT=3306

# 服务端口配置
BACKEND_PORT=8080
FRONTEND_PORT=80

# 时区配置
TZ=Asia/Shanghai
```

## 部署方式

### 方式一：Docker Compose 部署（推荐）

#### 1. 环境准备

确保Linux服务器已安装：
- Docker (版本 20.10+)
- Docker Compose (版本 2.0+)

```bash
# 安装Docker (Ubuntu/Debian)
curl -fsSL https://get.docker.com -o get-docker.sh
sudo sh get-docker.sh

# 安装Docker Compose
sudo curl -L "https://github.com/docker/compose/releases/latest/download/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose
sudo chmod +x /usr/local/bin/docker-compose

# 验证安装
docker --version
docker-compose --version
```

#### 2. 项目部署

```bash
# 1. 克隆或上传项目到服务器
git clone <your-repo-url> ai-agents
cd ai-agents

# 或者直接上传项目文件到服务器
scp -r ./AIagents2.0 user@your-server:/opt/ai-agents

# 2. 进入项目目录
cd /opt/ai-agents

# 3. 配置环境变量
cp .env.example .env
nano .env  # 编辑配置文件

# 4. 修改后端配置文件
nano backend/config.yaml

# 重要：修改数据库连接配置
# 将 host 改为 "mysql"（Docker服务名）
# 修改用户名密码等配置

# 5. 构建并启动服务
docker-compose up -d

# 6. 查看服务状态
docker-compose ps

# 7. 查看日志
docker-compose logs -f
```

#### 3. 配置文件修改示例

修改 `backend/config.yaml` 中的数据库配置：

```yaml
database:
  host: "mysql"              # 使用Docker服务名
  port: 3306
  username: "aiuser"         # 与.env中的MYSQL_USER一致
  password: "aipassword"     # 与.env中的MYSQL_PASSWORD一致
  dbname: "ai_agents"        # 与.env中的MYSQL_DATABASE一致
```

#### 4. 端口配置

在 `.env` 文件中配置暴露的端口：

```bash
# 前端访问端口
FRONTEND_PORT=80

# 后端API端口
BACKEND_PORT=8080

# MySQL端口（如需外部访问）
MYSQL_PORT=3306
```

#### 5. 访问应用

- 前端界面：`http://your-server-ip:80`
- 后端API：`http://your-server-ip:8080`
- 健康检查：`http://your-server-ip:8080/health`

### 方式二：手动部署

#### 1. 安装依赖

```bash
# 安装Go (版本 1.21+)
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# 安装MySQL
sudo apt update
sudo apt install mysql-server

# 安装Nginx
sudo apt install nginx
```

#### 2. 配置数据库
```bash
# 创建MySQL数据库
mysql -u root -p
CREATE DATABASE ai_agents CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
EXIT;
```

3. **配置后端**
```bash
cd backend
# 修改config.yaml文件中的数据库连接信息
nano config.yaml
```

4. **安装依赖并运行**
```bash
# 安装Go依赖
go mod tidy

# 运行后端服务
go run main.go
```

5. **访问应用**
打开浏览器访问 `http://localhost:8080`

## 常用运维命令

### Docker Compose 命令

```bash
# 启动所有服务
docker-compose up -d

# 停止所有服务
docker-compose down

# 重启服务
docker-compose restart

# 查看服务状态
docker-compose ps

# 查看日志
docker-compose logs -f [service_name]

# 进入容器
docker-compose exec backend sh
docker-compose exec mysql mysql -u root -p

# 更新服务
docker-compose pull
docker-compose up -d --force-recreate

# 清理资源
docker-compose down -v  # 删除卷
docker system prune -a  # 清理未使用的镜像
```

### 服务管理命令

```bash
# 查看服务状态
sudo systemctl status ai-agents-backend

# 重启服务
sudo systemctl restart ai-agents-backend

# 查看日志
sudo journalctl -u ai-agents-backend -f

# 重载Nginx配置
sudo nginx -t && sudo systemctl reload nginx
```

## 监控和日志

### 1. 应用日志

- 后端日志：`logs/app.log`
- Nginx日志：`/var/log/nginx/`
- MySQL日志：`/var/log/mysql/`

### 2. 健康检查

```bash
# 检查后端服务
curl http://localhost:8080/health

# 检查数据库连接
docker-compose exec mysql mysql -u aiuser -p ai_agents -e "SELECT 1"
```

### 3. 性能监控

```bash
# 查看容器资源使用
docker stats

# 查看系统资源
htop
df -h
free -h
```

## 故障排除

### 常见问题

1. **数据库连接失败**
   - 检查MySQL服务状态
   - 验证配置文件中的数据库连接信息
   - 确认防火墙设置

2. **端口冲突**
   - 修改 `.env` 文件中的端口配置
   - 检查端口占用：`netstat -tlnp | grep :8080`

3. **权限问题**
   - 确保文件权限正确：`chown -R www-data:www-data /opt/ai-agents`
   - 检查SELinux设置（如适用）

4. **内存不足**
   - 增加服务器内存
   - 调整MySQL配置参数

### 日志查看

```bash
# Docker容器日志
docker-compose logs backend
docker-compose logs mysql
docker-compose logs frontend

# 系统日志
sudo journalctl -u ai-agents-backend
sudo tail -f /var/log/nginx/error.log
```

## 安全建议

1. **数据库安全**
   - 使用强密码
   - 限制数据库访问IP
   - 定期备份数据

2. **网络安全**
   - 配置防火墙规则
   - 使用HTTPS（配置SSL证书）
   - 限制不必要的端口暴露

3. **应用安全**
   - 定期更新依赖包
   - 配置适当的日志级别
   - 实施访问控制

## 备份和恢复

### 数据备份

```bash
# 数据库备份
docker-compose exec mysql mysqldump -u root -p ai_agents > backup_$(date +%Y%m%d).sql

# 完整备份
tar -czf ai-agents-backup-$(date +%Y%m%d).tar.gz /opt/ai-agents
```

### 数据恢复

```bash
# 恢复数据库
docker-compose exec -T mysql mysql -u root -p ai_agents < backup_20240101.sql

# 恢复应用
tar -xzf ai-agents-backup-20240101.tar.gz -C /opt/
```

## 更新升级

### 应用更新

```bash
# 1. 备份当前版本
cp -r /opt/ai-agents /opt/ai-agents-backup-$(date +%Y%m%d)

# 2. 更新代码
git pull origin main

# 3. 重新构建和部署
docker-compose down
docker-compose build --no-cache
docker-compose up -d

# 4. 验证更新
docker-compose ps
curl http://localhost:8080/health
```

## 使用说明

### 创建智能体
1. 点击"创建智能体"按钮
2. 填写智能体名称、角色和性格描述
3. 点击"创建"完成智能体创建

### 创建聊天室
1. 点击"创建圆桌"按钮
2. 输入讨论主题
3. 选择参与讨论的智能体（至少2个）
4. 点击"开始讨论"启动自动化对话

### 观察讨论
- 智能体会自动进行讨论
- 消息会实时显示在聊天窗口中
- 可以随时停止讨论

## API文档

### 智能体相关接口

- `GET /agents` - 获取所有智能体
- `POST /agents` - 创建新智能体
- `DELETE /agents/:id` - 删除智能体

### 聊天室相关接口

- `GET /chatrooms` - 获取所有聊天室
- `POST /chatrooms` - 创建新聊天室
- `POST /chatrooms/:id/start` - 开始讨论
- `POST /chatrooms/:id/stop` - 停止讨论

### WebSocket接口

- `WS /ws` - WebSocket连接，用于实时消息推送

## 联系支持

如遇到部署问题，请提供以下信息：
- 操作系统版本
- Docker和Docker Compose版本
- 错误日志
- 配置文件内容（隐藏敏感信息）

---

**注意**: 请根据实际环境调整配置参数，确保系统安全和稳定运行。