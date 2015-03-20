default: test

test: unit-test
	@godep go test -coverprofile ./test.cov -v .
	@go tool cover -func ./test.cov | awk '$$3 !~ /^100/ { print; gaps++ } END { exit gaps }'

unit-test: .godep-install
	godep go test -cover .

.godep-install: Godeps/Godeps.json
	command -v godep > /dev/null || go get github.com/tools/godep
	godep restore
	touch .godep-install

setup: .godep-install

clean:
	rm -f ./test.cov

.PHONY: clean default setup test unit-test
