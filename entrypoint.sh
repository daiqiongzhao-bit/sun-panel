#!/bin/bash
set -e

# 创建持久化数据目录
mkdir -p /app/data/files/temp /app/data/database /app/data/backup /app/conf

# 确保自定义脚本/样式桩文件存在（避免挂载空卷导致 404 白屏）
mkdir -p /app/web/custom
[ -f /app/web/custom/index.js ] || printf '// Sun-Panel custom user script (empty)\n' > /app/web/custom/index.js
[ -f /app/web/custom/index.css ] || printf '/* Sun-Panel custom user style (empty) */\n' > /app/web/custom/index.css

if [ ! -f /app/conf/conf.ini ]; then
    echo "[init] Copying default config..."
    cp /app/conf-default/conf.ini /app/conf/
fi

# 使用软链接将文件存储和数据库指向持久化卷
# 如果旧路径已存在目录/文件（非软链接），先迁移
if [ -d /app/files ] && [ ! -L /app/files ]; then
    echo "[init] Migrating files to persistent volume..."
    cp -rn /app/files/* /app/data/files/ 2>/dev/null || true
    rm -rf /app/files
fi
if [ ! -d /app/files ] && [ ! -L /app/files ]; then
    ln -sf /app/data/files /app/files
fi

if [ -f /app/database.db ] && [ ! -L /app/database.db ]; then
    echo "[init] Migrating database to persistent volume..."
    cp /app/database.db /app/data/database/database.db
    rm /app/database.db
fi
if [ ! -f /app/database.db ] && [ ! -L /app/database.db ]; then
    ln -sf /app/data/database/database.db /app/database.db
fi

echo "[init] Starting Sun-Panel..."
exec ./sun-panel
