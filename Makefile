.PHONY: all build test clean atlaspack fontpack tilepack recolor test-atlaspack test-fontpack test-tilepack

all: build

build: atlaspack fontpack tilepack recolor

install: build
	cp -fv bin/* /usr/local/bin

test: test-atlaspack test-fontpack test-tilepack recolor

clean:
	rm -rf bin
	rm -rf cmd/fontpack/example/output
	rm -rf cmd/tilepack/example/output

atlaspack:
	@mkdir -p bin
	go build -ldflags="-s -w" -o bin/atlaspack$(BINARY_SUFFIX) ./cmd/atlaspack

fontpack:
	@mkdir -p bin
	go build -ldflags="-s -w" -o bin/fontpack$(BINARY_SUFFIX) ./cmd/fontpack

tilepack:
	@mkdir -p bin
	go build -ldflags="-s -w" -o bin/tilepack$(BINARY_SUFFIX) ./cmd/tilepack

recolor:
	@mkdir -p bin
	go build -ldflags="-s -w" -o bin/recolor$(BINARY_SUFFIX) ./cmd/recolor

test-atlaspack: atlaspack
	bin/atlaspack cmd/atlaspack/example/atlas.yaml cmd/atlaspack/example/output/game
	@diff -r cmd/atlaspack/example/output cmd/atlaspack/example/expected_output || (echo "FAIL: fontpack output differs from expected" && exit 1)
	echo "SUCCESS: atlaspack output is identical to expected"

test-fontpack: fontpack
	bin/fontpack cmd/fontpack/example/fonts.yaml cmd/fontpack/example/output
	@diff -r cmd/fontpack/example/output cmd/fontpack/example/expected_output || (echo "FAIL: fontpack output differs from expected" && exit 1)
	echo "SUCCESS: fontpack output is identical to expected"

test-tilepack: tilepack
	bin/tilepack cmd/tilepack/example/tilemaps cmd/tilepack/example/output
	@diff -r cmd/tilepack/example/output cmd/tilepack/example/expected_output || (echo "FAIL: tilepack output differs from expected" && exit 1)
	echo "SUCCESS: tilepack output is identical to expected"
