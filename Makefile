DEBUG ?= 0

all: util/util.go ganga.go
	@mkdir -p build
	go build -ldflags "-X ganga/util.DebugMode=$(DEBUG)" -mod=vendor -o build .

clean:
	rm -rf build
