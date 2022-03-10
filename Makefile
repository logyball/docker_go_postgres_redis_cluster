up:
	docker-compose up -d

start: up
	sleep 1  # allow redis container to start in cluster mode
	go run main.go

down:
	docker-compose down --remove-orphans
