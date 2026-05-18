BINARY := ghi
CMD := ./cmd/ghi
PREFIX ?= $(HOME)/.local
BINDIR := $(PREFIX)/bin
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
LDFLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: build install uninstall clean release

build:
	go build $(LDFLAGS) -o $(BINARY) $(CMD)

install: build
	mkdir -p $(BINDIR)
	cp $(BINARY) $(BINDIR)/$(BINARY)
	@echo "installed to $(BINDIR)/$(BINARY)"
	@echo "make sure $(BINDIR) is in your PATH"

uninstall:
	rm -f $(BINDIR)/$(BINARY)

clean:
	rm -f $(BINARY)

release:
	@echo "building release binary..."
	CGO_ENABLED=0 go build $(LDFLAGS) -trimpath -o $(BINARY) $(CMD)
	@echo "done: $(BINARY)"
	@echo "run 'make install' to copy to $(BINDIR)"
