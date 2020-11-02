build:
	docker build -t mailing-service .

run:
	docker run --net="host" \
		-p 50054 \
		-e MICRO_SERVER_ADDRESS=:50054 \
		-e MICRO_REGISTRY=mdns \
		-e MICRO_BROKER=nats \
		mailing-service