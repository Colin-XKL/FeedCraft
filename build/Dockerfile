FROM golang:1.19-alpine AS backend-builder
ENV GOCACHE=/root/.cache/go-build
ENV GOMODCACHE=/root/.cache/go-mod-cache

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod ./

# download Go modules and dependencies
RUN --mount=type=cache,id=gomod,target="/root/.cache/go-mod-cache" go mod download

# copy directory files i.e all files ending with .go
COPY . .
# compile application

RUN --mount=type=cache,id=gobuild,target="/root/.cache/go-build" GO111MODULE=on CGO_ENABLED=0 go build -o /feed-craft /app/cmd/main.go

FROM node:22-slim as frontend-builder
ENV PNPM_HOME="/pnpm"
ENV PATH="$PNPM_HOME:$PATH"

RUN npm install -g pnpm

WORKDIR /app

COPY web/admin/pnpm-lock.yaml web/admin/package.json /app/
RUN --mount=type=cache,id=pnpm,target=/pnpm/store pnpm install

COPY  web/admin /app/
RUN pnpm run build

## DEPLOY
FROM alpine

WORKDIR /app
ENV FC_DB_SQLITE_PATH=/usr/local/feed-craft/db/
RUN mkdir -p /usr/local/feed-craft/db/

COPY /web/start.html /web/favicon.ico /app/web/
COPY --from=backend-builder /feed-craft /app/feed-craft
COPY --from=frontend-builder /app/dist /app/web

ENV GIN_MODE=release
ENV LISTEN_ADDR=0.0.0.0:80
EXPOSE 80

ENTRYPOINT ["/app/feed-craft"]