.PHONY: all build test clean fontpack recolor tilepack test-fontpack test-tilepack

all: build

build: fontpack recolor tilepack

test: test-fontpack recolor test-tilepack

clean:
	rm -rf bin
	rm -rf cmd/fontpack/example/output
	rm -rf cmd/tilepack/example/output

fontpack:
	@mkdir -p bin
	go build -o bin/fontpack ./cmd/fontpack

recolor:
	@mkdir -p bin
	go build -o bin/recolor ./cmd/recolor

tilepack:
	@mkdir -p bin
	go build -o bin/tilepack ./cmd/tilepack

test-fontpack: fontpack
	@mkdir -p cmd/fontpack/example/output
	cd cmd/fontpack/example && ../../../bin/fontpack fonts.yaml output
	@diff -r cmd/fontpack/example/output cmd/fontpack/example/expected_output || (echo "FAIL: fontpack output differs from expected" && exit 1)
	echo "SUCCESS: fontpack output is identical to expected"

test-tilepack: tilepack
	@mkdir -p cmd/tilepack/example/output
	cd cmd/tilepack/example && ../../../bin/tilepack tilemaps output
	@diff -r cmd/tilepack/example/output cmd/tilepack/example/expected_output || (echo "FAIL: tilepack output differs from expected" && exit 1)
	echo "SUCCESS: tilepack output is identical to expected"
