VERSION:= 0.1.0

PROJECT_NAME:=salary

OUTPUT_FILE:=
	ifeq ($(OS), Windows_NT)
		OUTPUT_FILE+=./build/$(PROJECT_NAME).exe
	else
		UNAME_S:=$(shell uname -s)
		ifeq ($(UNAME_S), Linux)
			OUTPUT_FILE+=./build/$(PROJECT_NAME)
		endif
		ifeq ($(UNAME_S), Darwin)
			OUTPUT_FILE+=./build/$(PROJECT_NAME)_osx
		endif
	endif

.PHONY: build
build:
	go build -o $(OUTPUT_FILE) -v ./main.go

.PHONY: buildprod
buildprod:
	env GOOS=linux GOARCH=amd64 go build -o  $(OUTPUT_FILE)  -v ./main.go

.PHONY: run
run:
	$(OUTPUT_FILE) start

.PHONY: install
install:
	$(OUTPUT_FILE) install

DEFAULT_COAL := build