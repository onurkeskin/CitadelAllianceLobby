#!/bin/sh

docker-machine create \
  --driver virtualbox \
  --virtualbox-cpu-count 2 \
  --virtualbox-memory 1024 \
  --virtualbox-disk-size 20000 \
  swarm-manager-1

for i in 1 2; do
docker-machine create \
    --driver virtualbox \
    --virtualbox-cpu-count 2 \
    --virtualbox-memory 1024 \
    --virtualbox-disk-size 20000 \
    swarm-node-$i
done

eval $(docker-machine env swarm-manager-1)

docker swarm init --advertise-addr $(docker-machine ip swarm-manager-1)

TOKEN=$(docker swarm join-token -q worker)

for i in 1 2; do
  eval $(docker-machine env swarm-node-$i)
  docker swarm join --token $TOKEN $(docker-machine ip swarm-manager-1):2377
done

echo "Swarm cluster has been successfuly created !";

eval $(docker-machine env swarm-manager-1)

docker node ls