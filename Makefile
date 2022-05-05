BUILD=./cmd/app

all: clean build

build: clean
	go build -o $(BUILD)/bin $(BUILD)

test db:
	go test ./test

clean:
	if [ -f $(BUILD)/bin ]; then rm $(BUILD)/bin; fi