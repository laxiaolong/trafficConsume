# Binary name
BINARY=trafficConsume
DATE=$(shell date -Iseconds)
# Builds the project
build:
		GO111MODULE=on go build -trimpath -o ${BINARY} -ldflags "-X main.date=${DATE}"
release:
		# Clean
		rm -rf *.gz

		# Build for mac
		go clean
		GO111MODULE=on go build -trimpath -ldflags "-s -w -X main.version=${VERSION} -X main.date=${DATE}"
		tar czvf ${BINARY}-mac64-${VERSION}.tar.gz ./${BINARY}
		# Build for arm
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=arm64 GO111MODULE=on go build -trimpath -ldflags "-s -w -X main.version=${VERSION} -X main.date=${DATE}"
		tar czvf ${BINARY}-arm64-${VERSION}.tar.gz ./${BINARY}
		# Build for linux
		go clean
		CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GO111MODULE=on go build -trimpath -ldflags "-s -w -X main.version=${VERSION} -X main.date=${DATE}"
		tar czvf ${BINARY}-linux64-${VERSION}.tar.gz ./${BINARY}
		# Build for win
		go clean
		CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GO111MODULE=on go build -trimpath -ldflags "-s -w -X main.version=${VERSION} -X main.date=${DATE}"
		zip ${BINARY}-win64-${VERSION}.zip ./${BINARY}.exe

		go clean
# Cleans our projects: deletes binaries
clean:
		go clean
		rm -rf *.gz *.zip

.PHONY:  clean build