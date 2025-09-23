.DEFAULT_GOAL := install

# INSTALL_SCRIPT=./install.sh
BIN_FILE=little-txtreader

install:
	go build -o "${BIN_FILE}"

clean:
	go clean
	rm -f "cp.out"
	rm -f nohup.out
	rm -f "little-txtreader"

test:
	go test

check:
	go test

cover:
	go test -coverprofile cp.out
	go tool cover -html=cp.out

run:
	./little-txt-book-reader