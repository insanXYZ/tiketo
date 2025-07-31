run:
	go run cmd/app/main.go

create-migrate: install-migrate
	migrate create -ext sql -dir db/migrations -seq create_$(name)_table

install-migrate:
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest
