.PHONY: all build test clean fontpack recolor tilepack atlaspack test-fontpack test-tilepack

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
	go build -o bin/atlaspack ./cmd/atlaspack

fontpack:
	@mkdir -p bin
	go build -o bin/fontpack ./cmd/fontpack

tilepack:
	@mkdir -p bin
	go build -o bin/tilepack ./cmd/tilepack

recolor:
	@mkdir -p bin
	go build -o bin/recolor ./cmd/recolor

test-atlaspack: atlaspack
	@mkdir -p cmd/atlaspack/example/output
	bin/atlaspack cmd/atlaspack/example/atlas.yaml cmd/atlaspack/example/output/game
	@diff -r cmd/atlaspack/example/output cmd/atlaspack/example/expected_output || (echo "FAIL: fontpack output differs from expected" && exit 1)
	echo "SUCCESS: atlaspack output is identical to expected"

test-fontpack: fontpack
	@mkdir -p cmd/fontpack/example/output
	bin/fontpack cmd/fontpack/example/fonts.yaml cmd/fontpack/example/output
	@diff -r cmd/fontpack/example/output cmd/fontpack/example/expected_output || (echo "FAIL: fontpack output differs from expected" && exit 1)
	echo "SUCCESS: fontpack output is identical to expected"

test-tilepack: tilepack
	@mkdir -p cmd/tilepack/example/output
	bin/tilepack cmd/tilepack/example/tilemaps cmd/tilepack/example/output
	@diff -r cmd/tilepack/example/output cmd/tilepack/example/expected_output || (echo "FAIL: tilepack output differs from expected" && exit 1)
	echo "SUCCESS: tilepack output is identical to expected"
