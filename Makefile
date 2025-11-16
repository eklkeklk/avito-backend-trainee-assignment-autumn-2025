.PHONY: lint
lint:
	@golangci-lint run

.PHONY: lint-fix
lint-fix:
	@golangci-lint run --fix

.PHONY: lint-conf
lint-conf:
	@golangci-lint run --config=.golangci.yml

.PHONY: all
all: lint-conf