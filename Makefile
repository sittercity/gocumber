test: deps
	$(GOPATH)/bin/godep go test ./...

deps: .godep-install

.godep-install: $(GOPATH)/bin/godep $(GOPATH)/src/github.com/sittercity/gocumber/Godeps/Godeps.json
	$(GOPATH)/bin/godep restore
	touch .godep-install

