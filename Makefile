GO_SRCS := $(shell find . -type f -name '*.go')

docker-build:
	bash ./scripts/build.sh

count-go: ## Count number of lines of all go codes.
	find . -name "*.go" -type f | xargs wc -l | tail -n 1

go-build: $(GO_SRCS)
	go build ./
