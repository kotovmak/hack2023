git_tag = $(shell git describe --abbrev=0 --tags)

.PHONY: build
build:
	make test
	GOOS=linux GOARCH=amd64 go build -o ./dist/hack.app ./cmd/hack2023
	ssh -f 91 'screen -S hack -X quit'
	scp ./dist/hack.app 91:/home/bitrix/www/hack2023/
	$(eval REVISION = $(shell git log -1 --pretty=format:"%H"))
	ssh -f 91 'export VERSION="${git_tag}" && export REVISION="$(REVISION)" && export APP_TIER="prod" && screen -dmS hack /home/bitrix/www/hack2023/hack.app /dev/null 2>&1 &'

.PHONY: test
test:
	go test -cover -race -v -timeout 30s ./internal/app/...

.PHONY: migrate-create
migrate-create:
	migrate create -ext sql -dir migrations/ -seq $(name)

.PHONY: up
up:
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/hack2023" up

.PHONY: down
down:
	migrate -path migrations -database "mysql://root:password@tcp(localhost:3306)/hack2023" down

.PHONY: run
run: 
	go fmt ./internal/app/...
	go vet -composites=false ./internal/app/...
	go test -cover -race -v -timeout 30s ./internal/app/...
	go run ./cmd/hack2023

.PHONY: swagger
swagger:
	~/go/bin/swag init -g cmd/hack2023/main.go --exclude dist/

.DEFAULT_GOAL := run
