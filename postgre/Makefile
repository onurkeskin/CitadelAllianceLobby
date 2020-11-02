build:
	docker build -t citadelalliancedb .
run:
	docker run --name=caDB --mount source=citadelalliance-db-data,target=/var/lib/postgresql/10/main -d citadelalliancedb