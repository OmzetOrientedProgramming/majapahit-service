include .env

migrate-up :
	migrate -source file:database/postgres/migrations -database postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable up $(N)
migrate-down :
	migrate -source file:database/postgres/migrations -database postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable down $(N)
migrate-force :
	migrate -source file:database/postgres/migrations -database postgres://$(DB_USERNAME):$(DB_PASSWORD)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable force $(N)
test :
	go test -v $$(go list ./... | grep -v ./main.go | grep -v /vendor/)
coverage :
	go test $$(go list ./... | grep -v ./main.go | grep -v /vendor/) -coverprofile=coverage.out && go tool cover -func coverage.out
lint:
	$$(go list -f {{.Target}} golang.org/x/lint/golint) -set_exit_status $$(go list ./... | grep -v /vendor/)
fmt:
	go fmt $$(go list ./... | grep -v /vendor/)
