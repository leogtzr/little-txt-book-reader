.DEFAULT_GOAL := install

# INSTALL_SCRIPT=./install.sh
BIN_FILE=little-txt-book-reader

install:
	go build -o "${BIN_FILE}"

clean:
	go clean
	rm --force "cp.out"

test:
	go test

check:
	go test

cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out
