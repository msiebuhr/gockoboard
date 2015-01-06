.PHONY: all test fmt benchmark git-add-hook

all: test fmt vet

vet:
	go vet ./...

clean:
	go clean ./...

test:
	go test ./...

fmt:
	go fmt ./...

benchmark:
	go test ./... -bench=".*"

git-pre-commit-hook:
	curl -s 'https://raw.githubusercontent.com/golang/go/master/misc/git/pre-commit' > .git/hooks/pre-commit
	chmod +x .git/hooks/pre-commit
