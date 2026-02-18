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

test: test-fontpack recolor repack

test-fontpack: fontpack
	@mkdir -p cmd/fontpack/example/output
	cd cmd/fontpack/example && ../../../$(BINDIR)/fontpack fonts.json output
	@echo "Verifying output files..."
	@test -f cmd/fontpack/example/output/mana_branches.xml || (echo "Missing mana_branches.xml" && exit 1)
	@test -f cmd/fontpack/example/output/mana_roots.xml || (echo "Missing mana_roots.xml" && exit 1)
	@test -f cmd/fontpack/example/output/mana_trunk.xml || (echo "Missing mana_trunk.xml" && exit 1)
	@echo "All font files generated successfully!"
