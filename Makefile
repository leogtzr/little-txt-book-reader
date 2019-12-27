# INSTALL_SCRIPT=./install.sh
BIN_FILE=little-txt-book-reader

install:
	go build -o "${BIN_FILE}"

clean:
	go clean

test:
	go test

check:
	go test
