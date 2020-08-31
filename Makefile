MANAGER_CMD := manager
AGENT_CMD := agent

manager:
	go build -o $(MANAGER_CMD) ./cmd/manager/main.go

agent:
	go build -o $(AGENT_CMD) ./cmd/agent/main.go

test:
	go test -v -count=1 -race -gcflags=-l -timeout=30s ./...

.PHONY: manager agent test
