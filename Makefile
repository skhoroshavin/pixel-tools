.PHONY: all build clean test fontpack recolor repack test-fontpack

BINDIR := bin

all: build

build: fontpack recolor repack

fontpack:
	@mkdir -p $(BINDIR)
	go build -o $(BINDIR)/fontpack ./cmd/fontpack

recolor:
	@mkdir -p $(BINDIR)
	go build -o $(BINDIR)/recolor ./cmd/recolor

repack:
	@mkdir -p $(BINDIR)
	go build -o $(BINDIR)/repack ./cmd/repack

clean:
	rm -rf $(BINDIR)
	rm -rf cmd/fontpack/example/output

test: test-fontpack test-repack recolor

test-fontpack: fontpack
	@mkdir -p cmd/fontpack/example/output
	cd cmd/fontpack/example && ../../../$(BINDIR)/fontpack fonts.json output
	@diff -r cmd/fontpack/example/output cmd/fontpack/example/expected_output || (echo "FAIL: fontpack output differs from expected" && exit 1)
	echo "SUCCESS: fontpack output is identical to expected"

test-repack: repack
	@mkdir -p cmd/repack/example/output
	cd cmd/repack/example && ../../../$(BINDIR)/repack areas output
	@diff -r cmd/repack/example/output cmd/repack/example/expected_output || (echo "FAIL: repack output differs from expected" && exit 1)
	echo "SUCCESS: repack output is identical to expected"
