.PHONY: all gena ganga clean

DEBUG ?= 0

all: ganga gena

ganga:
	@mkdir -p build
	go build -ldflags "-X ganga/util.DebugMode=$(DEBUG)" -mod=vendor -o build .

gena:
	@mkdir -p build
	cd gena && go build -ldflags "-X ganga/util.DebugMode=$(DEBUG)" -mod=vendor -o ../build .

clean:
	rm -rf build
