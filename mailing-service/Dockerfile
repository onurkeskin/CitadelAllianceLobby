FROM golang:latest as builder
LABEL maintainer="onurkeskindev@gmail.com"

WORKDIR /go/src/keon.com/CitadelAllianceLobbyServer

COPY mailing-service mailing-service
COPY user-service/proto user-service/proto

WORKDIR /go/src/keon.com/CitadelAllianceLobbyServer/user-service/proto/user
RUN go get -d


WORKDIR /go/src/keon.com/CitadelAllianceLobbyServer/mailing-service

RUN go get
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo .


FROM alpine:latest

RUN apk --no-cache add ca-certificates

RUN mkdir /app
WORKDIR /app
COPY --from=builder /go/src/keon.com/CitadelAllianceLobbyServer/mailing-service/mailing-service .

CMD ["./mailing-service"]