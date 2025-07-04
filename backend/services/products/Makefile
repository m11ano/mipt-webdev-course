# app

.PHONY: run
run:
	go run cmd/app/main.go

# tests

PKGS=$(shell go run scripts/list_pkgs.go)

.PHONY: test
test:
	go test $(PKGS) -coverprofile=coverage.out && go tool cover -html=coverage.out -o coverage.html

.PHONY: test-integration
test-integration:
	go test ./tests -v

.PHONY: test-cov
test-cov:
	go test $(PKGS) -coverprofile=coverage.out
	go tool cover -func=coverage.out

# swagger

.PHONY: swagger
swagger:
	swag init --generalInfo cmd/app/main.go

# linters

.PHONY: lint
lint:
	golangci-lint run 
	
.PHONY: lint fix
lint-fix:
	golangci-lint run --fix

# migrations
# go install github.com/mikefarah/yq/v4@latest
# go install github.com/pressly/goose/v3/cmd/goose@latest

DBSTRING := $(shell yq e '.db.uri' config.yml)
GOOSE := goose -dir migrations postgres "$(DBSTRING)"

.PHONY: goose
goose:
	$(GOOSE) $(filter-out $@,$(MAKECMDGOALS))

%:
	@: