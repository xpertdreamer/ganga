all: util/util.go ganga.go
	@mkdir -p build
	go build -ldflags "-X ganga/util.DebugMode=$(DEBUG)" -o build .

clean:
	rm -rf build
