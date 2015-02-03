test: deps
	$(GOPATH)/bin/godep go test ./...

deps: .godep-install

.godep-install: $(GOPATH)/bin/godep $(GOPATH)/src/github.com/sittercity/gocumber/Godeps/Godeps.json
	$(GOPATH)/bin/godep restore
	touch .godep-install


test-cov: deps test-deps
	mkdir -p reports/cov
	$(GOPATH)/bin/gocov test github.com/sittercity/gocumber | $(GOPATH)/bin/gocov-xml > reports/cov/gocumber.xml

test-cov-html: deps test-deps
	mkdir -p reports/cov
	$(GOPATH)/bin/gocov test github.com/sittercity/gocumber | $(GOPATH)/bin/gocov-html > reports/cov/gocumber.html

test-deps: $(GOPATH)/bin/gocov $(GOPATH)/bin/gocov-xml $(GOPATH)/bin/gocov-html $(GOPATH)/bin/go-junit-report

$(GOPATH)/bin/gocov:
	go get github.com/axw/gocov/...

$(GOPATH)/bin/gocov-xml:
	go get github.com/AlekSi/gocov-xml

$(GOPATH)/bin/gocov-html:
	go get gopkg.in/matm/v1/gocov-html

$(GOPATH)/bin/go-junit-report:
	go get github.com/wancw/go-junit-report
