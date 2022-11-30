.EXPORT_ALL_VARIABLES:
this_dir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))


.PHONY: clean
clean:
	docker-compose down

.PHONY: run
run:
	docker-compose up --build

.PHONY: support
support:
	docker-compose up db

.PHONY: complete_test
complete_test:
	@go test -v -timeout 30s -coverprofile=cov.out ./...
	@go tool cover -func=cov.out

.PHONY: test
test:
	@go test -tags=unit -v -timeout 30s -coverprofile=cov.out ./...
	@go tool cover -func=cov.out

.PHONY: integration_test
integration_test:
	@go test -tags=integration ./internal/integration_test -count=1 -v

.PHONY: generate-mocks
generate-mocks:
	mockgen -source internal/db/models.go -destination internal/db/mock/mock_db.go
	mockgen -source internal/service/models.go -destination internal/service/mock/mock_service.go
