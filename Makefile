MANAGER_CMD := manager
AGENT_CMD := agent

manager:
	go build -o $(MANAGER_CMD) ./cmd/manager/main.go

agent:
	go build -o $(AGENT_CMD) ./cmd/agent/main.go

ut:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

int: SHELL:=/bin/bash
int:
	set -xeuo pipefail;\
	docker-compose -f tests/docker-compose-test.yaml up --build -d;\
	test_status_code=0;\
	docker-compose -f tests/docker-compose-test.yaml logs > logs.txt;\
	docker-compose -f tests/docker-compose-test.yaml run integration_tests go test || test_status_code=$$?;\
	docker-compose -f tests/docker-compose-test.yaml down;\
	if grep -q fail logs.txt; then \
		cat logs.txt | grep fail;\
		test_status_code=1;\
	fi ;\
	exit $$test_status_code

.PHONY: manager agent test ut int
