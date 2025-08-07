run:
	go mod tidy
	go run cmd/app/main.go

create-migrate: install-migrate
	migrate create -ext sql -dir db/migrations -seq create_$(name)_table && touch entity/$(name).go

install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

migrate-up:
	migrate -path db/migrations -database postgres://$(username):$(password)@$(host):$(port)/$(dbname)?sslmode=disable up

migrate-down:
	migrate -path db/migrations -database postgres://$(username):$(password)@$(host):$(port)/$(dbname)?sslmode=disable down

create-model:
	touch controller/$(name)_controller.go service/$(name)_service.go repository/$(name)_repository.go
