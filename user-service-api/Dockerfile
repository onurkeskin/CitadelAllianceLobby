# user-service/Dockerfile
FROM golang:latest as builder
LABEL maintainer="onurkeskindev@gmail.com"

WORKDIR /go/src/keon.com/CitadelAllianceLobbyServer/user-service-api

COPY . .

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -o user-service-api-linux-amd64 -a -installsuffix cgo .

FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/keon.com/CitadelAllianceLobbyServer/user-service-api/user-service-api-linux-amd64 .
COPY --from=builder /go/src/keon.com/CitadelAllianceLobbyServer/user-service-api/certs ./certs

EXPOSE 6767

CMD ["./user-service-api-linux-amd64"]