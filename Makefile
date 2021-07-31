NAME     = kf2-antiddos
VERSION := $(shell git describe)
GOCMD    = go
LDFLAGS := "$(LDFLAGS) -s -w -X 'main.Version=$(VERSION)'"
GOBUILD  = $(GOCMD) build -ldflags=$(LDFLAGS)
SRCMAIN  = ./cmd/$(NAME)
SRCDOC   = ./doc
BINDIR   = bin
BIN      = $(BINDIR)/$(NAME)
README   = $(SRCDOC)/README
LICENSE  = LICENSE
PREFIX   = /usr

.PHONY: all prep doc build check-build linux-amd64 windows-amd64 compile install check-install uninstall clean

all: build

prep: clean
	go mod init $(NAME); go mod tidy
	mkdir $(BINDIR)
	
doc: check-build
	test -d $(SRCDOC) || mkdir $(SRCDOC)
	$(BIN) --help > $(README)
	
build: prep
	$(GOBUILD) -o $(BIN) $(SRCMAIN)

check-build:
	test -e $(BIN)

linux-amd64: prep
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BIN)-linux-amd64 $(SRCMAIN)

windows-amd64: prep
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BIN)-windows-amd64.exe $(SRCMAIN)

compile: linux-386 windows-386 linux-amd64 windows-amd64
	
install: check-build doc
	install -m 755 -d         $(PREFIX)/bin/
	install -m 755 $(BIN)     $(PREFIX)/bin/
	install -m 755 -d         $(PREFIX)/share/licenses/$(NAME)/
	install -m 644 $(LICENSE) $(PREFIX)/share/licenses/$(NAME)/
	install -m 755 -d         $(PREFIX)/share/doc/$(NAME)/
	install -m 644 $(README)  $(PREFIX)/share/doc/$(NAME)/

check-install:
	test -e $(PREFIX)/bin/$(NAME) || \
	test -d $(PREFIX)/share/licenses/$(NAME) || \
	test -d $(PREFIX)/share/doc/$(NAME)

uninstall: check-install
	rm -f  $(PREFIX)/bin/$(NAME)
	rm -rf $(PREFIX)/share/licenses/$(NAME)
	rm -rf $(PREFIX)/share/doc/$(NAME)

clean:
	rm -rf $(BINDIR)
