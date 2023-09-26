BUILD_OPTION=-ldflags="-s -w"
BUILD_NAME=EIP1559
SOURCE=./cmd/main.go
INSTALL_PATH=/data
ENV_PATH=env
BUILD_PATH=./bin

.PHONY: build clean build-darwin build-linux
default: build
all: clean build build-darwin build-linux
build:
	go build ${BUILD_OPTION} -o ${BUILD_PATH}/${BUILD_NAME} ${SOURCE}
build-darwin:
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build ${BUILD_OPTION} -o ${BUILD_PATH}/${BUILD_NAME}-dawin ${SOURCE}
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build ${BUILD_OPTION} -o ${BUILD_PATH}/${BUILD_NAME}-linux-amd64 ${SOURCE}
clean:
	rm -rf ${BUILD_PATH}
cleaninstall:
	rm -rf ${INSTALL_PATH}/${BUILD_NAME}/
install:
	cleaninstall
	build
	cp -r ${BUILD_PATH}/${BUILD_NAME} ${INSTALL_PATH}/${BUILD_NAME}
	cp -r ${ENV_PATH}/prod.yml ${INSTALL_PATH}/${BUILD_NAME}/${ENV_PATH}/

