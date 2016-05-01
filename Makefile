PACKAGE=zbx_mongo

#.DEFAULT_GOAL: ${PACKAGE}

all:
	go build -ldflags "-s -w" -buildmode=c-shared -o ${PACKAGE}.so ./${PACKAGE}

build-debug:
	go build -buildmode=c-shared -o ${PACKAGE}.so ./${PACKAGE}

clean:
	go clean


