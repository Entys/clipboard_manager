VERSION ?= $$(git describe --tags 2>/dev/null || git rev-parse --short HEAD 2>/dev/null || echo "latest")
LEVEL ?= release
OUTPUT ?= sshs$$(if [ "$${GOOS:-$$(go env GOOS)}" == "windows" ]; then echo '.exe'; else echo ''; fi)

GO_PACKAGE_PATH := github.com/Entys/clipboard_manager

ifeq ($(LEVEL),release)
GOLDFLAGS := -w -s
endif

build:
	go build -ldflags "$(GOLDFLAGS) -X '$(GO_PACKAGE_PATH)/cmd.Version=$(or $(strip $(VERSION)),latest)'" -o $(OUTPUT)

clean:
	rm -f clipboard

PREFIX ?= /usr/local

install: clipboard
	mkdir -p $(DESTDIR)$(PREFIX)/bin
	cp $< $(DESTDIR)$(PREFIX)/bin/clipboard

uninstall:
	rm -f $(DESTDIR)$(PREFIX)/bin/clipboard

default: build

.PHONY: clean install uninstall
