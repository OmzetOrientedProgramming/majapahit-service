include .env

migrate-up :
	migrate -source file:database/postgres/migrations -database postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable up $(N)
migrate-down :
	migrate -source file:database/postgres/migrations -database postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable down $(N)
test :
	go test -v $$(go list ./... | grep -v ./main.go) 
coverage :
	go test $$(go list ./... | grep -v ./main.go) -coverprofile=coverage.out && go tool cover -func coverage.out
