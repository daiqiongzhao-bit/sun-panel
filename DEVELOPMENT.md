# CDF-Panel 开发与发布指南

> 本文档详细介绍 CDF-Panel 的开发环境搭建、代码管理、CI/CD 流水线以及版本发布流程。

---

## 目录

- [一、项目概览](#一项目概览)
- [二、开发环境搭建](#二开发环境搭建)
- [三、代码管理规范](#三代码管理规范)
- [四、CI/CD 自动化流水线](#四cicd-自动化流水线)
- [五、版本发布流程](#五版本发布流程)
- [六、服务器部署流程](#六服务器部署流程)
- [七、Docker Hub 管理](#七docker-hub-管理)
- [八、常见问题与故障排查](#八常见问题与故障排查)
- [九、版本历史](#九版本历史)

---

## 一、项目概览

### 1.1 仓库信息

| 项目 | 地址 |
|------|------|
| GitHub 仓库 | [daiqiongzhao-bit/sun-panel](https://github.com/daiqiongzhao-bit/sun-panel) |
| Docker Hub | [cdf3275/cdf-panel](https://hub.docker.com/r/cdf3275/cdf-panel) |
| 在线演示 | [https://960898.xyz](https://960898.xyz) |
| 服务器 | `49.51.200.165` (CentOS, root 用户) |
| Docker 容器名 | `cdf-panel` |
| 服务端口 | `3030` |

### 1.2 技术栈

| 层级 | 技术 | 版本 |
|------|------|------|
| 前端 | Vue 3 + TypeScript | - |
| UI 框架 | Naive UI | 2.x |
| 样式 | Tailwind CSS | 3.x |
| 构建 | Vite | 4.x |
| 后端 | Go + Gin | - |
| ORM | GORM | - |
| 数据库 | SQLite | - |
| 部署 | Docker 多阶段构建 | - |

### 1.3 项目结构

```
sun-panel/
├── src/                          # 前端源码
│   ├── api/                      # API 请求封装
│   ├── components/               # 组件
│   │   ├── apps/                 #   管理后台组件
│   │   ├── common/               #   通用组件
│   │   └── deskModule/           #   桌面模块（便签/时钟/搜索/监控）
│   ├── views/                    # 页面
│   │   ├── home/                 #   首页
│   │   ├── login/                #   登录页
│   │   └── paste/                #   粘贴板
│   ├── router/                   # 路由配置
│   ├── store/                    # Pinia 状态管理
│   ├── locales/                  # 国际化（中/英）
│   ├── utils/                    # 工具函数
│   └── main.ts                   # 入口文件
├── service/                      # 后端 Go 源码
│   ├── main.go                   # 入口
│   ├── api/api_v1/               # API 控制器
│   │   ├── common/               #   公共（返回格式、基础控制器）
│   │   ├── middleware/           #   中间件（认证、Gzip 压缩）
│   │   ├── panel/                #   面板业务
│   │   └── system/               #   系统管理业务
│   ├── models/                   # 数据模型
│   ├── router/                   # 路由注册
│   ├── initialize/               # 初始化（配置、数据库等）
│   ├── global/                   # 全局变量
│   ├── lib/                      # 工具库
│   └── assets/                   # 静态资源（版本号等）
├── public/                       # 静态资源
├── config/                       # 前端配置
├── doc/                          # 文档和截图
├── .github/workflows/            # GitHub Actions 工作流
├── Dockerfile                    # Docker 构建文件
├── docker-compose.yml            # Docker Compose 配置
├── build.sh                      # 本地构建脚本
├── entrypoint.sh                 # 容器启动脚本
├── README.md                     # 项目说明
├── UPDATELOG.md                  # 更新日志
├── DEPLOYMENT.md                 # 部署手册
└── DEVELOPMENT.md                # 本文档
```

---

## 二、开发环境搭建

### 2.1 前置要求

| 工具 | 最低版本 | 用途 |
|------|---------|------|
| Node.js | 16+ | 前端开发 |
| pnpm | 7+ | 前端包管理 |
| Go | 1.20+ | 后端开发 |
| Git | 2.0+ | 版本管理 |
| Docker | 20+ | 容器构建 |

### 2.2 一键安装开发环境（Windows）

```powershell
# 1. 安装 Node.js
winget install OpenJS.NodeJS.LTS

# 2. 安装 pnpm
npm install -g pnpm

# 3. 安装 Go
winget install GoLang.Go

# 4. 克隆项目
git clone https://github.com/daiqiongzhao-bit/sun-panel.git
cd sun-panel

# 5. 安装前端依赖
pnpm install

# 6. 安装后端依赖
cd service
go mod download
cd ..
```

### 2.3 启动开发模式

```bash
# 终端1：启动后端（端口 3030）
cd service
go run main.go

# 终端2：启动前端（端口 1002，代理到后端 3030）
pnpm dev
```

访问 `http://localhost:1002` 进入开发模式，修改代码自动热重载。

### 2.4 本地构建

```bash
# 前端构建
pnpm build
# 产物在 dist/ 目录

# 后端构建
cd service
go build -ldflags "-s -w" -o ../sun-panel ./
# 产物为 sun-panel 可执行文件

# Docker 本地构建
docker build -t cdf3275/cdf-panel:dev .
```

---

## 三、代码管理规范

### 3.1 Git 分支策略

```
main         ← 主分支，始终保持可部署状态
  └─ dev     ← 开发分支（可选，功能开发时使用）
```

当前采用 **单分支策略**（main），所有修改直接推送到 main。

### 3.2 提交信息格式

```
<type>: <简短描述>

示例：
  feat: 新增多语言切换功能
  fix: 修复登录页背景不显示问题
  docs: 更新 README 部署说明
  style: 调整右下角菜单位置
  perf: 优化首页加载速度
  refactor: 重构用户权限判断逻辑
  chore: 升级依赖版本
```

### 3.3 日常工作流

```bash
# === 每天开始工作前 ===
git pull origin main        # 拉取最新代码

# === 开发过程中 ===
git status                  # 查看修改了哪些文件
git add 文件名              # 添加要提交的文件
git commit -m "fix: 修复xxx问题"   # 提交

# === 推送改动 ===
git push origin main        # 推送到 GitHub

# === 查看历史 ===
git log --oneline -10       # 查看最近 10 条提交
git show 提交ID             # 查看某次提交详情
```

### 3.4 撤销操作

```bash
git reset HEAD 文件名       # 取消暂存某个文件
git checkout -- 文件名       # 丢弃某个文件的修改
git reset --soft HEAD~1     # 撤销最近一次 commit（保留修改）
```

---

## 四、CI/CD 自动化流水线

### 4.1 流水线架构

```
┌─────────────┐     git push tag     ┌──────────────────┐
│  本地开发机  │ ──────────────────→ │   GitHub Actions  │
└─────────────┘                      └────────┬─────────┘
                                              │
                                    ┌─────────▼─────────┐
                                    │  1. Checkout 代码  │
                                    │  2. 读取版本号     │
                                    │  3. 设置 QEMU      │
                                    │  4. 设置 Buildx    │
                                    │  5. 登录 Docker Hub │
                                    │  6. 多架构构建推送  │
                                    └─────────┬─────────┘
                                              │
                                    ┌─────────▼─────────┐
                                    │    Docker Hub       │
                                    │  cdf3275/cdf-panel  │
                                    │  ├── v1.2.0         │
                                    │  ├── latest         │
                                    │  └── 支持 amd64     │
                                    │       arm64  armv7  │
                                    └────────────────────┘
```

### 4.2 工作流文件详解

文件位置：`.github/workflows/docker-build-push.yml`

```yaml
name: docker-push                    # 工作流名称

on:
  workflow_dispatch:                 # 允许手动触发
  push:
    tags:
      - 'v*'                        # 推送 v 开头的 tag 时自动触发

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
      # 步骤1：拉取代码
      - name: Checkout
        uses: actions/checkout@v3

      # 步骤2：从 service/assets/version 读取版本号
      - name: Read version from file
        id: read_version
        run: echo "APP_VERSION=$(cut -d '|' -f 2 ./service/assets/version)" >> $GITHUB_ENV

      # 步骤3-4：多架构构建环境
      - name: Set up QEMU
        uses: docker/setup-qemu-action@v2
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v2

      # 步骤5：登录 Docker Hub（使用 GitHub Secrets）
      - name: Login to Docker Hub
        uses: docker/login-action@v2
        with:
          username: ${{ secrets.DOCKERHUB_USERNAME }}
          password: ${{ secrets.DOCKERHUB_TOKEN }}

      # 步骤6：构建并推送三个架构的镜像
      - name: Build and push
        uses: docker/build-push-action@v4
        with:
          context: .
          platforms: linux/amd64,linux/arm,linux/arm64
          push: true
          tags: |
            ${{ vars.DOCKER_IMAGE_NAME }}:${{ env.APP_VERSION }}
            ${{ vars.DOCKER_IMAGE_NAME }}:latest
```

### 4.3 GitHub Secrets 配置

> ⚠️ 这些 Secret 是自动构建推送 Docker Hub 的关键，切勿删除或泄露。

| Secret 名称 | 用途 | 值 |
|-------------|------|-----|
| `DOCKERHUB_USERNAME` | Docker Hub 用户名 | `cdf3275` |
| `DOCKERHUB_TOKEN` | Docker Hub 访问令牌 | 在 Docker Hub 生成 |

| Variable 名称 | 用途 | 值 |
|---------------|------|-----|
| `DOCKER_IMAGE_NAME` | Docker 镜像名称 | `cdf3275/cdf-panel` |

#### 如何更新 Docker Hub Token？

Token 过期后需重新生成：

1. 打开 https://hub.docker.com/settings/security
2. 点击 **New Access Token**
3. 描述填 `github-actions`，权限选 `Read & Write`
4. 复制生成的 token（格式：`dckr_pat_xxx...`）
5. 打开 https://github.com/daiqiongzhao-bit/sun-panel/settings/secrets/actions
6. 找到 `DOCKERHUB_TOKEN`，点击 **Update**，粘贴新 token，保存

### 4.4 如何查看构建状态

1. 打开 https://github.com/daiqiongzhao-bit/sun-panel/actions
2. 点击最新一次 `docker-push` 运行记录
3. 点击 `build` 任务可展开查看每个步骤的日志

---

## 五、版本发布流程

### 5.1 标准发布（推荐）

```bash
# === 第1步：修改代码 ===
# 正常开发，修改文件，提交代码...

# === 第2步：更新版本文档 ===
# 编辑 service/assets/version
# 将 v1.1.0 改为 v1.2.0（只改版本号部分）

# 编辑 UPDATELOG.md
# 在文件顶部添加新版本的更新内容，格式参考：

## v1.2.0
> 支持上个版本直接升级无需特殊处理

- [新增] xxx 功能
- [修复] xxx 问题
- [优化] xxx 性能

# 更新 README.md（如有必要，更新功能说明和截图）

# === 第3步：提交所有改动 ===
git add .
git commit -m "chore: 发布 v1.2.0"
git push origin main

# === 第4步：打标签（自动触发构建）===
git tag v1.2.0
git push origin v1.2.0
```

**打完标签后无需任何操作**，GitHub Actions 会自动：
1. 拉取代码
2. 构建前端 + 后端
3. 打包 Docker 镜像（amd64 / arm64 / armv7）
4. 推送到 Docker Hub `cdf3275/cdf-panel:v1.2.0` 和 `cdf3275/cdf-panel:latest`

### 5.2 从服务器推送（当前工作流）

由于当前开发环境在服务器（`49.51.200.165`）上，完整的发布流程如下：

```bash
# === 第1步：SSH 登录服务器 ===
ssh root@49.51.200.165

# === 第2步：进入项目并拉取最新 ===
cd /root/sun-panel/sun-panel
git pull origin main

# === 第3步：修改代码/文件 ===
# 编辑需要改动的文件...

# === 第4步：更新版本号 ===
vim service/assets/version
# 改为 1|v1.3.0

# === 第5步：更新更新日志 ===
vim UPDATELOG.md
# 在顶部添加新版本说明

# === 第6步：提交 ===
git add .
git commit -m "chore: 发布 v1.3.0"
git push origin main

# === 第7步：打标签 ===
git tag v1.3.0
git push origin v1.3.0

# === 第8步：等待 GitHub Actions 构建完成 ===
# 约 10-15 分钟后，Docker Hub 上会出现 v1.3.0 镜像

# === 第9步：拉取新镜像并重启容器 ===
docker pull cdf3275/cdf-panel:v1.3.0
docker rm -f cdf-panel
docker run -d --name cdf-panel \
  -p 3030:3030 \
  --restart unless-stopped \
  -v /root/docker_data/cdf-panel/conf:/app/conf \
  -v /root/docker_data/cdf-panel/uploads:/app/uploads \
  -v /root/docker_data/cdf-panel/database:/app/database \
  cdf3275/cdf-panel:v1.3.0

# === 第10步：验证 ===
docker ps | grep cdf-panel
curl -s -o /dev/null -w "%{http_code}" http://localhost:3030/
# 返回 200 即为正常
```

### 5.3 完整发布检查清单

每次发布前确认以下事项：

| 检查项 | ✓ |
|--------|---|
| `service/assets/version` 版本号已更新 | ☐ |
| `UPDATELOG.md` 已更新 | ☐ |
| `README.md` Docker 命令版本号已更新 | ☐ |
| 所有改动已提交到 GitHub | ☐ |
| Git tag 已推送（格式 `v1.x.x`） | ☐ |
| GitHub Actions 构建成功（绿色 ✓） | ☐ |
| Docker Hub 出现新版本镜像 | ☐ |
| 服务器容器已更新到新版本 | ☐ |

---

## 六、服务器部署流程

### 6.1 服务器信息

| 项目 | 值 |
|------|-----|
| IP 地址 | `49.51.200.165` |
| 系统 | CentOS (x86_64) |
| 用户 | `root` |
| Docker 容器名 | `cdf-panel` |
| 数据目录 | `/root/docker_data/cdf-panel/` |

### 6.2 数据目录结构

```
/root/docker_data/cdf-panel/
├── conf/          # 配置文件 (conf.ini)
├── uploads/       # 用户上传的文件（Logo、壁纸等）
├── database/      # SQLite 数据库文件
└── runtime/       # 运行日志
```

### 6.3 日常运维命令

```bash
# === 查看容器状态 ===
docker ps -a | grep cdf-panel

# === 查看日志 ===
docker logs cdf-panel           # 所有日志
docker logs cdf-panel --tail 50 # 最近 50 行
docker logs -f cdf-panel        # 实时滚动

# === 重启容器 ===
docker restart cdf-panel

# === 进入容器 ===
docker exec -it cdf-panel /bin/sh

# === 重置管理员密码 ===
docker exec cdf-panel ./sun-panel -password-reset
# 重置后密码为 12345678

# === 更新到新版本 ===
docker pull cdf3275/cdf-panel:v1.2.0
docker rm -f cdf-panel
docker run -d --name cdf-panel \
  -p 3030:3030 \
  --restart unless-stopped \
  -v /root/docker_data/cdf-panel/conf:/app/conf \
  -v /root/docker_data/cdf-panel/uploads:/app/uploads \
  -v /root/docker_data/cdf-panel/database:/app/database \
  cdf3275/cdf-panel:v1.2.0

# === 查看资源占用 ===
docker stats cdf-panel --no-stream
```

### 6.4 数据备份

```bash
# === 方法1：备份到本地 ===
cd /root/docker_data/cdf-panel
tar -czf /root/sun-panel/backup_$(date +%Y%m%d_%H%M%S).tar.gz \
  database/ conf/ uploads/

# === 方法2：使用管理后台 ===
# 登录 → 备份管理 → 创建备份

# === 方法3：定时自动备份（crontab）===
crontab -e
# 添加以下行（每天凌晨 3 点备份）
0 3 * * * tar -czf /root/sun-panel/auto_backup_$(date +\%Y\%m\%d).tar.gz \
  /root/docker_data/cdf-panel/database/ \
  /root/docker_data/cdf-panel/conf/ \
  /root/docker_data/cdf-panel/uploads/
```

### 6.5 数据恢复

```bash
# === 停止容器 ===
docker stop cdf-panel

# === 恢复数据 ===
cd /root/docker_data/cdf-panel
tar -xzf /root/sun-panel/backup_20260718_120000.tar.gz

# === 启动容器 ===
docker start cdf-panel
```

---

## 七、Docker Hub 管理

### 7.1 账号信息

| 项目 | 值 |
|------|-----|
| 用户名 | `cdf3275` |
| 仓库 | [cdf3275/cdf-panel](https://hub.docker.com/r/cdf3275/cdf-panel) |
| 管理地址 | https://hub.docker.com/repository/docker/cdf3275/cdf-panel |

### 7.2 Docker Hub 操作

```bash
# === 本地登录 Docker Hub ===
docker login -u cdf3275
# 输入 Access Token（不是密码）

# === 手动拉取镜像 ===
docker pull cdf3275/cdf-panel:v1.2.0
docker pull cdf3275/cdf-panel:latest

# === 手动推送镜像 ===
docker tag 本地镜像 cdf3275/cdf-panel:自定义tag
docker push cdf3275/cdf-panel:自定义tag

# === 查看镜像标签 ===
# 打开 https://hub.docker.com/r/cdf3275/cdf-panel/tags
```

### 7.3 镜像架构支持

通过 GitHub Actions 自动构建的镜像支持以下架构：

| 架构 | 适用设备 |
|------|---------|
| `linux/amd64` | 标准 x86_64 服务器、PC |
| `linux/arm64` | ARM64 设备（树莓派 4/5、Apple Silicon） |
| `linux/arm/v7` | ARMv7 设备（树莓派 2/3） |

---

## 八、常见问题与故障排查

### Q1：GitHub Actions 构建失败怎么办？

**步骤：**

1. 打开 Actions 页面：https://github.com/daiqiongzhao-bit/sun-panel/actions
2. 点击失败的运行记录
3. 查看 `build` 任务 → 展开失败步骤 → 查看错误日志
4. 常见失败原因：
   - Docker Hub 登录失败 → Token 过期，重新生成并更新 Secret
   - 前端构建失败 → `pnpm install` 或 `pnpm build` 报错，检查代码语法
   - 后端构建失败 → Go 编译报错，检查语法

### Q2：推送了 tag 但 Actions 没有触发？

```bash
# 确认 tag 格式正确（必须以 v 开头）
git tag                          # 查看所有本地 tag
git push origin --tags           # 确保所有 tag 都推送了

# 手动触发
# 在 GitHub Actions 页面点击 Run workflow → 选择分支 → Run workflow
```

### Q3：服务器上拉取镜像很慢？

```bash
# 从 Docker Hub 拉取时可能会慢，耐心等待或配置国内镜像加速
docker pull cdf3275/cdf-panel:latest
```

### Q4：容器启动后无法访问？

```bash
# 检查容器是否运行
docker ps -a | grep cdf-panel

# 检查日志
docker logs cdf-panel --tail 20

# 常见原因：
# - 端口被占用：用 netstat -tlnp | grep 3030 检查
# - 数据库文件损坏：检查 /root/docker_data/cdf-panel/database/
# - 配置文件错误：检查 /root/docker_data/cdf-panel/conf/conf.ini
```

### Q5：Git 推送时要求输入密码？

服务器上已配置 Token 认证，无需输入密码。如果出现提示：

```bash
# 检查远程地址（应包含 token）
git remote -v

# 如果地址中不含 token，重新设置：
git remote set-url origin https://daiqiongzhao-bit:你的GitHub_Token@github.com/daiqiongzhao-bit/sun-panel.git
```

### Q6：如何回滚到旧版本？

```bash
# 拉取旧版本镜像
docker pull cdf3275/cdf-panel:v1.1.0

# 停止当前容器
docker stop cdf-panel && docker rm cdf-panel

# 启动旧版本（注意数据可能不完全兼容）
docker run -d --name cdf-panel \
  -p 3030:3030 \
  --restart unless-stopped \
  -v /root/docker_data/cdf-panel/conf:/app/conf \
  -v /root/docker_data/cdf-panel/uploads:/app/uploads \
  -v /root/docker_data/cdf-panel/database:/app/database \
  cdf3275/cdf-panel:v1.1.0
```

---

## 九、版本历史

| 版本 | 日期 | 说明 |
|------|------|------|
| v1.2.0 | 2026-07-18 | 修复 6 个 Bug、12 个安全问题，多项性能优化 |
| v1.1.0 | 2026-07 | 企业功能增强（RBAC、部门、审计、通知、粘贴板密码保护） |
| v1.0.0 | 2026-07 | 基于 Sun-Panel 的首次二次开发版本 |

> 详细更新内容请查看 [UPDATELOG.md](./UPDATELOG.md)

---

## 附录：快速命令速查表

```bash
# ========== 开发 ==========
pnpm dev                              # 启动前端开发服务器
cd service && go run main.go         # 启动后端
pnpm build                            # 构建前端
docker build -t cdf3275/cdf-panel:dev . # 本地 Docker 构建

# ========== Git ==========
git status                            # 查看状态
git add .                             # 暂存所有改动
git commit -m "fix: xxx"             # 提交
git push origin main                  # 推送代码
git tag v1.2.0                        # 创建标签
git push origin v1.2.0               # 推送标签（触发构建）

# ========== 服务器 ==========
ssh root@49.51.200.165               # 登录服务器
docker logs cdf-panel --tail 50      # 查看日志
docker restart cdf-panel              # 重启容器
docker exec -it cdf-panel /bin/sh    # 进入容器

# ========== Docker Hub ==========
docker pull cdf3275/cdf-panel:v1.2.0 # 拉取镜像
docker login -u cdf3275               # 登录

# ========== 备份 ==========
tar -czf /root/sun-panel/backup.tar.gz /root/docker_data/cdf-panel/
```

---

> 最后更新：2026-07-18 · v1.2.0
