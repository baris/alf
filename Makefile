NAME        := alf
VERSION     := 0.0.1
SRCDIR      := src
PKGS        := $(NAME)
SOURCES     := $(foreach pkg, $(PKGS), $(wildcard $(SRCDIR)/$(pkg)/*.go))

# symlinks confuse go tools, let's not mess with it and use -L
GOPATH  := $(shell pwd -L)
export GOPATH

PATH := bin:$(PATH)
export PATH

all: clean $(NAME)

.PHONY: clean
clean:
	@echo Cleaning $(NAME)...
	@rm -f $(NAME) bin/$(NAME)

deps:
	@echo Getting dependencies...
	@$(foreach pkg, $(PKGS), go get -t $(pkg);)

$(NAME): $(SOURCES) deps
	@echo Building $(NAME)...
	@go build -o bin/$(NAME) $@
