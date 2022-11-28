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
	go test -v ./...


.PHONY: generate-mocks
generate-mocks:
	mockgen -source internal/db/models.go -destination internal/db/mock/mock_db.go
	mockgen -source internal/service/models.go -destination internal/servjce/mock/mock_service.go
