.EXPORT_ALL_VARIABLES:
this_dir := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))


.PHONY: setup
setup:
	podman-compose -f docker-compose.setup.yml up

.PHONY: clean
clean:
	podman-compose down

.PHONY: run
run:
	podman-compose up -d --build

.PHONY: support
support:
	podman-compose up -d db elasticsearch filebeat localstack prometheus grafana

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
	@go test -tags=service ./internal/integration_test -count=1 -v

.PHONY: generate-mocks
generate-mocks:
	mockgen -source internal/db/models.go -destination internal/db/mock/mock_db.go
	mockgen -source internal/service/models.go -destination internal/service/mock/mock_service.go

.PHONY: build
build:
	docker build -t aborgesrodrigues/todo -f dev/docker/Dockerfile .
	docker push aborgesrodrigues/todo

.PHONY: deploy
deploy:
	kubectl apply -f dev/kubernetes --recursive

.PHONY: undeploy
undeploy:
	kubectl delete -f dev/kubernetes --recursive

.PHONY: generate_data
generate_data:
	go run cmd/generate_data/main.go

.PHONY: generate_trace
generate_trace:
	curl -o trace.out "http://localhost:8080/debug/pprof/trace?seconds=10" && go tool trace trace.out  

.PHONY: pprof_heap
pprof_heap:
	go tool pprof -http=:9091 http://localhost:8080/debug/pprof/heap

.PHONY: pprof_allocs
pprof_allocs:
	go tool pprof -http=:9092 http://localhost:8080/debug/pprof/allocs

.PHONY: pprof_cpu
pprof_cpu:
	go tool pprof -http=:9091 http://localhost:8080/debug/pprof/profile?seconds=300
# 	curl -o profile.out "http://localhost:8080/debug/pprof/profile?seconds=300" && go tool pprof profile.out  
