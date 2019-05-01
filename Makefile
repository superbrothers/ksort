.PHONY: build
build:
		@OSARCH="$(shell uname -s | tr "[:upper:]" "[:lower:]")/amd64" hack/make-binaries.sh

.PHONY: build-cross
build-cross:
		@hack/make-binaries.sh

.PHONY: archives
archives:
		@hack/make-all.sh

.PHONY: test
test:
		go test -v ./...

.PHONY: test-archive
test-archive:
		@hack/test-archive.sh

.PHONY: clean
clean:
		rm -rf out
