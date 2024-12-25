dev:
	@nodemon --watch './**/*.go' --signal SIGTERM --exec 'go' run main.go

build:
	@go build -o instashop

run:
	@./instashop

build_run: build run
