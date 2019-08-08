# INSTALL_SCRIPT=./install.sh
BIN_FILE=txtreader

install:
	go build -o "${BIN_FILE}"
	# ${INSTALL_SCRIPT}
	# cp ${BIN_FILE} ~/bin

clean:
	go clean

