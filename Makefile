default: test

test: unit-test
	mkdir -p reports/cov reports/unit
	godep go test -coverprofile reports/gocumber.cov -v . > reports/gocumber.txt
	[ ! -f reports/gocumber.cov ] || gocov convert reports/gocumber.cov | gocov-xml > reports/cov/gocumber.xml
	[ ! -f reports/gocumber.cov ] || gocov convert reports/gocumber.cov | gocov-html > reports/cov/gocumber.html
	[ ! -f reports/gocumber.txt ] || go-junit-report < reports/gocumber.txt > reports/unit/gocumber.xml

unit-test: .godep-install
	godep go test -cover .

.godep-install: Godeps/Godeps.json
	command -v godep > /dev/null || go get github.com/tools/godep
	godep restore
	touch .godep-install

setup: .godep-install
	command -v gocov > /dev/null || go get github.com/axw/gocov/...
	command -v gocov-xml > /dev/null || go get github.com/AlekSi/gocov-xml
	command -v gocov-html > /dev/null || go get gopkg.in/matm/v1/gocov-html
	command -v go-junit-report > /dev/null || go get github.com/wancw/go-junit-report

clean:
	rm -rf reports
