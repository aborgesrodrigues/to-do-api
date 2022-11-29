.EXPORT_ALL_VARIABLES:
this_dir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))


.PHONY: clean
clean:
	docker-compose down

.PHONY: run
run:
	docker-compose up --build

.PHONY: test
test:
	@go test -v -timeout 30s -coverprofile=cov.out ./...
	@go tool cover -func=cov.out


.PHONY: generate-mocks
generate-mocks:
	mockgen -source internal/db/models.go -destination internal/db/mock/mock_db.go
	mockgen -source internal/service/models.go -destination internal/service/mock/mock_service.go
