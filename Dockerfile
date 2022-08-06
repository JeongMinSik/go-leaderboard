## Build
FROM golang:1.18 AS build

ENV CGO_ENABLED=0
COPY ./app/ /go/src/app/
WORKDIR /go/src/app
RUN go build

## Deploy
FROM alpine

WORKDIR /
COPY --from=build /go/src/app/go-leaderboard /
ENTRYPOINT ["/go-leaderboard"]