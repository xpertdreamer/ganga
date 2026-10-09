DEBUG ?= 0

ganga:
	@mkdir -p build
	go build -ldflags "-X ganga/util.DebugMode=$(DEBUG)" -mod=vendor -o build .

clean:
	rm -rf build
