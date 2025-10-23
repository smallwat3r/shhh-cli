APP := shhh
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo 1.3.3)

PLATFORMS := linux darwin
ARCHES := amd64 arm64

export CGO_ENABLED=0

.PHONY: clean build pack release

clean:
	rm -rf dist

build:
	@echo ">> Building $(APP) $(VERSION)"
	gox \
	  -os="$(PLATFORMS)" \
	  -arch="$(ARCHES)" \
	  -output="dist/$(APP)_{{.OS}}_{{.Arch}}/$(APP)" \
	  -ldflags="-s -w -X main.shhhCLIVersion=$(VERSION) -buildid=" \
	  -gcflags=all=-trimpath=$(PWD)

pack: build
	@echo ">> Packaging artifacts"
	@cd dist && \
	for d in $(APP)_*; do \
	  if echo $$d | grep -qi windows; then \
	    cp "$$d/$(APP)" "$$d/$(APP).exe"; \
	    zip -q "$${d}.zip" "$$d/$(APP).exe"; \
	  else \
	    tar czf "$${d}.tar.gz" -C "$$d" $(APP); \
	  fi; \
	done

release: clean pack
	@echo "Artifacts in ./dist ready."
