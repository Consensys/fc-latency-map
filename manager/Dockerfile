# Copyright (C) 2020 ConsenSys Software Inc

FROM golang:1.17-alpine as builder
WORKDIR /go/src/github.com/ConsenSys/fc-latency-map
RUN apk update && apk add --no-cache make gcc musl-dev linux-headers git
COPY ./extern/ ./extern/
COPY ./manager/ ./manager/
WORKDIR /go/src/github.com/ConsenSys/fc-latency-map/manager
# Get all dependancies, but don't install.
RUN go get -d -v github.com/ConsenSys/fc-latency-map/manager/cmd/manager-server
# Do a full compile of app and dependancies, forcing static linking.
RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o /go/bin/manager-server ./cmd/manager-server

FROM alpine:latest
WORKDIR /
COPY --from=builder /go/bin/manager-server /manager-server
CMD ["/manager-server", "--host", "0.0.0.0", "--port", "3000"]
EXPOSE 3000
