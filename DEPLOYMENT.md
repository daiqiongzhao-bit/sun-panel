# Sun-Panel 部署操作手册

## 版本信息
- 版本号: v1.0.0
- 默认端口: 3030
- 默认账号: admin
- 默认密码: 12345678

---

## 目录
1. [服务器环境要求](#1-服务器环境要求)
2. [环境检查脚本](#2-环境检查脚本)
3. [安装 Docker](#3-安装-docker)
4. [上传项目文件](#4-上传项目文件)
5. [构建 Docker 镜像](#5-构建-docker-镜像)
6. [启动容器](#6-启动容器)
7. [服务验证](#7-服务验证)
8. [常见问题排查](#8-常见问题排查)
9. [安全建议](#9-安全建议)

---

## 1. 服务器环境要求

### 最低配置
| 项目 | 要求 |
|------|------|
| 操作系统 | Ubuntu 20.04 / CentOS 7+ |
| CPU | 1核 |
| 内存 | 1GB |
| 磁盘空间 | 10GB 可用空间 |
| 网络 | 支持外网访问（用于拉取镜像和依赖） |

### 推荐配置
| 项目 | 要求 |
|------|------|
| 操作系统 | Ubuntu 22.04 / CentOS 8+ |
| CPU | 2核 |
| 内存 | 2GB |
| 磁盘空间 | 20GB 可用空间 |
| 网络 | 支持外网访问 |

---

## 2. 环境检查脚本

在部署前，请运行以下命令检查服务器环境：

```bash
#!/bin/bash

echo "=========================================="
echo "        Sun-Panel 环境检查脚本"
echo "=========================================="

# 检查操作系统
echo ""
echo "[1] 操作系统检查"
if [ -f /etc/os-release ]; then
    . /etc/os-release
    echo "操作系统: $PRETTY_NAME"
    echo "内核版本: $(uname -r)"
else
    echo "操作系统: 未知"
fi

# 检查 CPU
echo ""
echo "[2] CPU 检查"
cpu_cores=$(grep -c ^processor /proc/cpuinfo)
echo "CPU 核心数: $cpu_cores 核"

# 检查内存
echo ""
echo "[3] 内存检查"
total_mem=$(free -h | grep Mem | awk '{print $2}')
used_mem=$(free -h | grep Mem | awk '{print $3}')
available_mem=$(free -h | grep Mem | awk '{print $7}')
echo "总内存: $total_mem"
echo "已用内存: $used_mem"
echo "可用内存: $available_mem"

# 检查磁盘空间
echo ""
echo "[4] 磁盘空间检查"
root_available=$(df -h / | grep / | awk '{print $4}')
echo "根目录可用空间: $root_available"

# 检查 Docker 是否已安装
echo ""
echo "[5] Docker 检查"
if command -v docker &> /dev/null; then
    echo "Docker 版本: $(docker --version)"
    echo "Docker 状态: $(systemctl is-active docker)"
else
    echo "Docker: 未安装"
fi

# 检查 Docker Compose 是否已安装
echo ""
echo "[6] Docker Compose 检查"
if command -v docker-compose &> /dev/null; then
    echo "Docker Compose 版本: $(docker-compose --version)"
elif docker compose version &> /dev/null; then
    echo "Docker Compose 版本: $(docker compose version)"
else
    echo "Docker Compose: 未安装"
fi

# 检查端口 3030 是否被占用
echo ""
echo "[7] 端口检查"
if lsof -Pi :3030 -sTCP:LISTEN -t >/dev/null 2>&1; then
    echo "端口 3030: 已被占用"
else
    echo "端口 3030: 可用"
fi

# 检查网络连接
echo ""
echo "[8] 网络检查"
if curl -s --head --request GET https://www.docker.com/ | grep "200 OK" > /dev/null; then
    echo "外网连接: 正常"
else
    echo "外网连接: 异常，请检查网络设置"
fi

echo ""
echo "=========================================="
echo "         环境检查完成"
echo "=========================================="
```

### 使用方法

将上述脚本保存为 `check_env.sh`，然后执行：

```bash
chmod +x check_env.sh
./check_env.sh
```

---

## 3. 安装 Docker

### Ubuntu 系统

```bash
# 更新软件包
sudo apt update && sudo apt upgrade -y

# 安装必要依赖
sudo apt install -y apt-transport-https ca-certificates curl software-properties-common

# 添加 Docker GPG 密钥
curl -fsSL https://download.docker.com/linux/ubuntu/gpg | sudo gpg --dearmor -o /usr/share/keyrings/docker-archive-keyring.gpg

# 添加 Docker 仓库
echo "deb [arch=$(dpkg --print-architecture) signed-by=/usr/share/keyrings/docker-archive-keyring.gpg] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" | sudo tee /etc/apt/sources.list.d/docker.list > /dev/null

# 更新并安装 Docker
sudo apt update && sudo apt install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 启动 Docker 服务
sudo systemctl start docker
sudo systemctl enable docker

# 验证安装
docker --version
docker compose version
```

### CentOS 系统

```bash
# 更新软件包
sudo yum update -y

# 安装必要依赖
sudo yum install -y yum-utils device-mapper-persistent-data lvm2

# 添加 Docker 仓库
sudo yum-config-manager --add-repo https://download.docker.com/linux/centos/docker-ce.repo

# 安装 Docker
sudo yum install -y docker-ce docker-ce-cli containerd.io docker-compose-plugin

# 启动 Docker 服务
sudo systemctl start docker
sudo systemctl enable docker

# 验证安装
docker --version
docker compose version
```

### 配置 Docker 镜像加速（可选，国内服务器推荐）

```bash
# 创建配置文件
sudo mkdir -p /etc/docker
sudo tee /etc/docker/daemon.json <<-'EOF'
{
    "registry-mirrors": [
        "https://docker.m.daocloud.io",
        "https://hub-mirror.c.163.com",
        "https://mirror.baidubce.com"
    ]
}
EOF

# 重启 Docker 服务
sudo systemctl daemon-reload
sudo systemctl restart docker
```

---

## 4. 上传项目文件

### 方法一：使用 scp 上传

```bash
# 在本地执行
scp /path/to/sun-panel-v1.0.0.zip root@您的服务器IP:/root/
```

### 方法二：使用 wget 下载（如果已上传到文件服务器）

```bash
# 在服务器上执行
cd /root
wget https://your-file-server.com/sun-panel-v1.0.0.zip
```

### 解压文件

```bash
cd /root
unzip sun-panel-v1.0.0.zip
cd sun-panel-v1.3.0
```

---

## 5. 构建 Docker 镜像

### 方式一：使用 docker build 命令

```bash
cd /root/sun-panel-v1.3.0

# 构建镜像（需要较长时间，请耐心等待）
docker build -t sun-panel:v1.0.0 .

# 查看构建的镜像
docker images | grep sun-panel
```

### 方式二：使用 Docker Compose

```bash
cd /root/sun-panel-v1.3.0

# 构建并启动（会自动构建镜像）
docker-compose up -d --build
```

### 构建过程说明

构建过程包含以下步骤，每个步骤可能需要一些时间：

1. **拉取 Node.js 镜像** - 用于构建前端
2. **安装前端依赖** - 使用 pnpm 安装
3. **构建前端代码** - 生成静态文件
4. **拉取 Go 镜像** - 用于构建后端
5. **安装 Go 依赖** - 安装 go-bindata 等工具
6. **构建后端二进制** - 生成可执行文件
7. **创建运行镜像** - 合并前端和后端

### 构建失败处理

如果构建失败，请检查以下内容：

1. **网络问题**：确保服务器可以访问外网，或配置 Docker 镜像加速
2. **磁盘空间**：确保有足够的磁盘空间（至少 5GB）
3. **内存不足**：增加服务器内存或配置 Docker 资源限制

---

## 6. 启动容器

### 方式一：使用 docker run 命令

```bash
# 创建数据目录
mkdir -p /root/sun-panel-data/conf
mkdir -p /root/sun-panel-data/uploads
mkdir -p /root/sun-panel-data/database

# 启动容器
docker run -d \
  --name sun-panel \
  -p 3030:3030 \
  -v /root/sun-panel-data/conf:/app/conf \
  -v /root/sun-panel-data/uploads:/app/uploads \
  -v /root/sun-panel-data/database:/app/database \
  --restart=always \
  sun-panel:v1.0.0
```

### 方式二：使用 Docker Compose

```bash
# 创建数据目录
mkdir -p /root/sun-panel-data/conf
mkdir -p /root/sun-panel-data/uploads
mkdir -p /root/sun-panel-data/database

# 修改 docker-compose.yml 中的挂载路径（如果需要）
# 默认配置已经包含了挂载路径

# 启动容器
docker-compose up -d
```

### 参数说明

| 参数 | 说明 |
|------|------|
| `-d` | 后台运行容器 |
| `--name sun-panel` | 设置容器名称 |
| `-p 3030:3030` | 端口映射（主机端口:容器端口） |
| `-v /root/sun-panel-data/conf:/app/conf` | 挂载配置目录 |
| `-v /root/sun-panel-data/uploads:/app/uploads` | 挂载上传文件目录 |
| `-v /root/sun-panel-data/database:/app/database` | 挂载数据库目录 |
| `--restart=always` | 容器重启策略（开机自启） |

---

## 7. 服务验证

### 检查容器状态

```bash
# 查看容器状态
docker ps

# 查看容器日志
docker logs sun-panel

# 查看容器详细信息
docker inspect sun-panel
```

### 测试服务访问

```bash
# 使用 curl 测试本地访问
curl -I http://localhost:3030

# 预期输出（包含以下内容）
# HTTP/1.1 200 OK
```

### 浏览器访问

在浏览器中访问：

```
http://您的服务器IP:3030
```

### 登录验证

| 项目 | 值 |
|------|------|
| 地址 | http://您的服务器IP:3030 |
| 账号 | admin |
| 密码 | 12345678 |

---

## 8. 常见问题排查

### 问题 1：端口被占用

**现象**：启动容器时提示端口 3030 已被占用

**解决方案**：

```bash
# 查看占用端口的进程
lsof -i :3030

# 杀死占用端口的进程（替换 <PID> 为实际进程ID）
kill -9 <PID>

# 或者修改端口映射（例如使用 80 端口）
docker run -d \
  --name sun-panel \
  -p 80:3030 \
  -v /root/sun-panel-data/conf:/app/conf \
  -v /root/sun-panel-data/uploads:/app/uploads \
  -v /root/sun-panel-data/database:/app/database \
  --restart=always \
  sun-panel:v1.0.0
```

### 问题 2：容器启动失败

**现象**：容器启动后立即退出

**解决方案**：

```bash
# 查看容器日志
docker logs sun-panel

# 常见原因：
# 1. 数据目录权限问题
# 2. 配置文件错误
# 3. 数据库文件损坏

# 修复权限问题
chown -R 1000:1000 /root/sun-panel-data/

# 查看容器启动命令
docker inspect sun-panel | grep "Cmd"
```

### 问题 3：无法访问服务

**现象**：浏览器无法访问 http://服务器IP:3030

**解决方案**：

```bash
# 检查防火墙设置
sudo ufw status

# 如果启用了防火墙，允许 3030 端口
sudo ufw allow 3030/tcp

# 检查 SELinux（CentOS）
getenforce

# 如果 SELinux 启用，添加规则或临时禁用
sudo setenforce 0
```

### 问题 4：构建过程缓慢

**现象**：Docker 构建过程非常慢

**解决方案**：

```bash
# 配置 Docker 镜像加速（参见第 3 节）

# 或者使用国内源构建
# 修改 Dockerfile，取消注释以下行：
# RUN npm config set registry https://repo.huaweicloud.com/repository/npm/
# RUN sed -i "s@dl-cdn.alpinelinux.org@mirrors.aliyun.com@g" /etc/apk/repositories \
#     && go env -w GOPROXY=https://goproxy.cn,direct
```

### 问题 5：数据库文件损坏

**现象**：服务启动后无法正常登录，日志显示数据库错误

**解决方案**：

```bash
# 停止容器
docker stop sun-panel

# 备份当前数据库
cp -r /root/sun-panel-data/database /root/sun-panel-data/database_backup

# 删除损坏的数据库文件
rm -rf /root/sun-panel-data/database/*

# 重新启动容器（会自动创建新数据库）
docker start sun-panel
```

---

## 9. 安全建议

### 1. 修改默认密码

登录系统后，立即修改默认密码：

1. 登录后台管理
2. 进入"用户管理"页面
3. 修改 admin 用户的密码

### 2. 配置防火墙

```bash
# 允许 SSH 连接（如果需要远程管理）
sudo ufw allow 22/tcp

# 允许 Web 服务端口
sudo ufw allow 3030/tcp

# 启用防火墙
sudo ufw enable
```

### 3. 使用 HTTPS

推荐配置 Nginx 反向代理并启用 HTTPS：

```nginx
server {
    listen 80;
    server_name your-domain.com;
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl;
    server_name your-domain.com;

    ssl_certificate /path/to/ssl/certificate.crt;
    ssl_certificate_key /path/to/ssl/private.key;

    location / {
        proxy_pass http://localhost:3030;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

### 4. 定期备份数据

```bash
# 创建备份脚本
cat > /root/backup_sun-panel.sh << 'EOF'
#!/bin/bash
BACKUP_DIR="/root/backups"
DATE=$(date +%Y%m%d_%H%M%S)

mkdir -p $BACKUP_DIR
cd /root/sun-panel-data
tar -czf $BACKUP_DIR/sun-panel_backup_$DATE.tar.gz conf database

# 保留最近 7 天的备份
find $BACKUP_DIR -name "sun-panel_backup_*.tar.gz" -mtime +7 -delete

echo "Backup completed: $BACKUP_DIR/sun-panel_backup_$DATE.tar.gz"
EOF

chmod +x /root/backup_sun-panel.sh

# 添加到 cron 定时任务（每天凌晨 2 点执行）
echo "0 2 * * * /root/backup_sun-panel.sh" | crontab -
```

### 5. 更新 Docker 镜像

```bash
# 停止并删除旧容器
docker stop sun-panel
docker rm sun-panel

# 重新构建镜像（如果有代码更新）
docker build -t sun-panel:v1.0.0 .

# 启动新容器
docker run -d \
  --name sun-panel \
  -p 3030:3030 \
  -v /root/sun-panel-data/conf:/app/conf \
  -v /root/sun-panel-data/uploads:/app/uploads \
  -v /root/sun-panel-data/database:/app/database \
  --restart=always \
  sun-panel:v1.0.0
```

---

## 附录：常用命令

| 命令 | 说明 |
|------|------|
| `docker ps` | 查看运行中的容器 |
| `docker ps -a` | 查看所有容器 |
| `docker logs sun-panel` | 查看容器日志 |
| `docker stop sun-panel` | 停止容器 |
| `docker start sun-panel` | 启动容器 |
| `docker restart sun-panel` | 重启容器 |
| `docker rm sun-panel` | 删除容器 |
| `docker rmi sun-panel:v1.0.0` | 删除镜像 |
| `docker exec -it sun-panel bash` | 进入容器 |
| `docker-compose up -d` | 使用 Compose 启动 |
| `docker-compose down` | 使用 Compose 停止 |
| `docker-compose logs` | 查看 Compose 日志 |

---

## 技术支持

如果在部署过程中遇到问题，请收集以下信息并联系技术支持：

1. 服务器操作系统和版本
2. Docker 版本
3. 完整的错误日志
4. 问题发生的具体步骤

---

*文档版本: v1.0.0*  
*最后更新: 2026-07-15*