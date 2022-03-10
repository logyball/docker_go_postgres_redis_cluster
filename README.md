# Golang, Redis Cluster Mode, Postgres

A very quick http server intended to demonstrate ability to spin up a local application stack:

1. Go backend using native `http` server
2. Redis implemented in [cluster mode](https://redis.io/topics/cluster-spec)
3. PostgreSQL

This is a very common application stack, but finding local solutions to mimicing Redis' cluster mode was something that I was finding difficult until I encountered [Grokzen's excellent repository linked here.](https://github.com/Grokzen/docker-redis-cluster/)

# Usage

### Startup

```bash
$ make start
docker-compose up -d
Creating network "docker_redis_cluster_default" with the default driver
Creating redis_cluster ... done
Creating pg            ... done
sleep 1  # allow redis container to start in cluster mode
go run main.go
INFO[0000] Postgres ping successful                     
INFO[0000] Redis Cluster Ping successful
```

### Shutdown
```bash
$ make down
docker-compose down --remove-orphans
Stopping redis ... done
Stopping pg    ... done
Removing redis ... done
Removing pg    ... done
Removing network docker_redis_cluster_default
```
