#build stage
FROM golang:alpine AS builder
RUN apk add --no-cache git
WORKDIR /go/src/bigdocument
COPY . .
RUN go get -d -v ./...
RUN GOBIN=/go/bin go install -v ./cmd/server

#final stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
COPY --from=builder /go/bin/server /server
COPY ./configs/test.json /config.json
ENTRYPOINT CONFIG_FILE_PATH=/config.json ./server
LABEL Name=bigdocument Version=0.0.1
EXPOSE 80
