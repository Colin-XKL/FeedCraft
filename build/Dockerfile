FROM golang:1.19-alpine AS backend-builder

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY . .
# compile application
RUN GO111MODULE=on CGO_ENABLED=0 go build -o /feed-craft /app/cmd/main.go

FROM node:22 as frontend-builder
WORKDIR /app
RUN npm install -g pnpm
COPY web/admin/pnpm-lock.yaml web/admin/package.json /app/
RUN pnpm install

COPY  web/admin /app/
RUN pnpm run build

## DEPLOY
FROM alpine

WORKDIR /app
ENV FC_DB_SQLITE_PATH=/usr/local/feed-craft/db/
RUN mkdir -p /usr/local/feed-craft/db/

COPY --from=backend-builder /feed-craft /app/feed-craft
COPY /web/start.html /web/favicon.ico /app/web/
COPY --from=frontend-builder /app/dist /app/web
ENV GIN_MODE=release
ENV LISTEN_ADDR=0.0.0.0:80
EXPOSE 80

ENTRYPOINT ["/app/feed-craft"]