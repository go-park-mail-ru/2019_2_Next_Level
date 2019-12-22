all: install build

.PHONY: build
build:
	go build -o ./build/auth_service cmd/Auth/main.go
	go build -o ./build/mailpicker_service cmd/MailPicker/main.go
	go build -o ./build/http_service cmd/serverapi/main.go
	go build -o ./build/post_service cmd/post/main.go

install:
	    go get ./...