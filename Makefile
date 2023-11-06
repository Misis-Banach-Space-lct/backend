swag:
	~/go/bin/swag init -g ./cmd/server/main.go -o ./docs
	~/go/bin/swag fmt

debug: swag
	docker compose -f deployment/docker-compose.yaml up db -d 
	go run cmd/server/main.go

local:
	docker compose -f deployment/docker-compose.yaml up --build -d