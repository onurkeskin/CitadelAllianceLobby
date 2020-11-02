FROM golang:latest as builder

WORKDIR /go/src/keon.com/CitadelAllianceLobbyServer

COPY user-cli user-cli
COPY user-service/proto user-service/proto


WORKDIR /go/src/keon.com/CitadelAllianceLobbyServer/user-service/proto/user
RUN go get -d

WORKDIR /go/src/keon.com/CitadelAllianceLobbyServer/user-cli

RUN go get -d
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/keon.com/CitadelAllianceLobbyServer/user-cli/user-cli .

ENTRYPOINT ["./user-cli"]

CMD ["./user-cli"]