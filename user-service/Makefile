GOPATH = /home/onur/GoWorkspace

SRC_DIR = $(GOPATH)/src/keon.com/CitadelAllianceLobbyServer/user-service
DST = $(SRC_DIR)/proto/user/user.proto

build:
	protoc -I=$(SRC_DIR) \
	 -I. \
	--go_out=plugins=grpc:$(SRC_DIR) $(DST)
	docker build -t user-service .
run:
	docker run -p 50052:50051
