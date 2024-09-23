# install & build frontend
FROM node:22.6-alpine3.19 AS frontendBuilder
WORKDIR /go/zula/client
COPY ./frontend ./
RUN npm install
COPY . .
RUN npm run build

# build backend
FROM golang:alpine as builder
RUN apk update && apk add --no-cache git ca-certificates

WORKDIR /go/zula/service

COPY ./service ./
COPY --from=frontendBuilder /go/zula/client/dist ./static
RUN go get -d -v
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 ENV=production go build -ldflags="-w -s" -o bin .
RUN ls -la /go/zula/service/bin

# release
FROM golang:alpine
WORKDIR /root/
COPY --from=builder /go/zula/service/bin /root/bin
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

EXPOSE 8081

ENTRYPOINT ["/root/bin"]
