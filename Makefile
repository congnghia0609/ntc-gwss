NAME=ntc-gwss
VERSION=0.0.1

.PHONY: build
build:
	@go build -o $(NAME)

.PHONY: run
run: build
	@./$(NAME) -e development

.PHONY: run-test
run-test: build
	@./$(NAME) -e test

.PHONY: run-stag
run-stag: build
	@./$(NAME) -e staging

.PHONY: run-prod
run-prod: build
	@./$(NAME) -e production

.PHONY: clean
clean:
	@rm -f $(NAME)

.PHONY: deps-save
deps-save:
	@godep save

.PHONY: deps
deps:
	@godep restore

.PHONY: test
test:
	@go test -v ./test/*
