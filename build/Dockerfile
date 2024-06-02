FROM golang:1.19-alpine AS build

# create a working directory inside the image
WORKDIR /app

# copy Go modules and dependencies to image
COPY go.mod ./

# download Go modules and dependencies
RUN go mod download

# copy directory files i.e all files ending with .go
COPY . .
# compile application
RUN go build -o /feed-craft /app/cmd/main.go

##
## STEP 2 - DEPLOY
##
FROM alpine

WORKDIR /app
ENV FC_DB_SQLITE_PATH=/usr/local/feed-craft/db/
RUN mkdir -p /usr/local/feed-craft/db/

COPY --from=build /feed-craft /app/feed-craft
COPY /web/index.html /web/favicon.ico /app/web/
ENV GIN_MODE=release
EXPOSE 8080

ENTRYPOINT ["/app/feed-craft"]