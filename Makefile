DEBUG ?= 0

all:
ganga:
	@mkdir -p build
	go build -ldflags "-X ganga/util.DebugMode=$(DEBUG)" -mod=vendor -o build .

gena:
	@mkdir -p build
	go build -mod=vendor -ldflags "-X ganga/util.DebugMode=$(DEBUG)" -o build ./gena

clean:
	rm -rf build
