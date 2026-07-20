# build frontend
FROM node:22 AS web_image

RUN npm config set registry https://repo.huaweicloud.com/repository/npm/

RUN npm install pnpm@8.15.6 -g

WORKDIR /build

COPY ./package.json /build

COPY ./pnpm-lock.yaml /build

RUN pnpm install --unsafe-perm

COPY . /build

# 追加版本号到 .env，不要覆盖原有的 VITE_GLOB_API_URL 等配置
RUN echo "VITE_APP_VERSION=1.0.0" >> /build/.env

RUN pnpm run build

# build backend
FROM golang:1.22-bookworm as server_image

WORKDIR /build

COPY ./service .

RUN go env -w GOPROXY=https://goproxy.cn,direct

RUN apt-get update && apt-get install -y --no-install-recommends bash curl gcc git

RUN go env -w GO111MODULE=on \
    && export PATH=$PATH:/go/bin \
    && go install -v github.com/go-bindata/go-bindata/...@latest \
    && go install -v github.com/elazarl/go-bindata-assetfs/...@latest \
    && go-bindata-assetfs -o=assets/bindata.go -pkg=assets assets/... \
    && go mod tidy \
    && go build -o sun-panel --ldflags="-X sun-panel/global.RUNCODE=release -X sun-panel/global.ISDOCKER=docker" main.go

# run_image
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=web_image /build/dist /app/web

# 预置自定义脚本/样式空桩，避免未挂载 custom 卷时 404 白屏
RUN mkdir -p /app/web/custom && \
    printf '// Sun-Panel custom user script (empty)\n' > /app/web/custom/index.js && \
    printf '/* Sun-Panel custom user style (empty) */\n' > /app/web/custom/index.css

COPY --from=server_image /build/sun-panel /app/sun-panel

COPY --from=server_image /build/assets/conf.example.ini /app/conf-default/conf.ini

COPY --from=server_image /build/assets/lang /app/assets/lang

COPY entrypoint.sh /app/entrypoint.sh

RUN apt-get update && apt-get install -y --no-install-recommends ca-certificates tzdata \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && chmod +x /app/entrypoint.sh /app/sun-panel \
    && cd /app/web && find . -type f \( -name "*.html" -o -name "*.css" -o -name "*.js" -o -name "*.svg" -o -name "*.json" \) -exec gzip -k9 {} \; && cd /app

EXPOSE 3030

ENTRYPOINT ["/app/entrypoint.sh"]
