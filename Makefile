MANAGER_CMD := manager
AGENT_CMD := agent

manager:
	go build -o $(MANAGER_CMD) ./cmd/manager/main.go

agent:
	go build -o $(AGENT_CMD) ./cmd/agent/main.go

ut:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

int:
	./tests/test.sh
	
.PHONY: manager agent test ut int
