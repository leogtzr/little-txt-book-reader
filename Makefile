# INSTALL_SCRIPT=./install.sh
BIN_FILE=txtreader

install:
	go build -o "${BIN_FILE}"

clean:
	go clean

