NAME = flood-it
OUTPUT = $(NAME)

.PHONY: build build-win

build:
	$(GOFLAGS) go build -o ./build/$(OUTPUT)
	
build-win: OUTPUT = $(NAME).exe
build-win: GOFLAGS += GOOS=windows
build-win: build
