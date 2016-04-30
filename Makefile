PACKAGE=zbx_mongo

#.DEFAULT_GOAL: ${PACKAGE}

all:
	go build -buildmode=c-shared -o ${PACKAGE}.so ./${PACKAGE}


clean:
	go clean


