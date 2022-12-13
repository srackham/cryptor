# cryptor Makefile

# Set defaults (see http://clarkgrubb.com/makefile-style-guide#prologue)
MAKEFLAGS += --warn-undefined-variables
SHELL := bash
.SHELLFLAGS := -eu -o pipefail -c
.DEFAULT_GOAL := test
.DELETE_ON_ERROR:
.SUFFIXES:
.ONESHELL:
.SILENT:

GOFLAGS ?=
PACKAGES = ./...
XFLAG_PATH = github.com/srackham/cryptor/internal/cli

.PHONY: install
install:
	LDFLAGS="-X $(XFLAG_PATH).BUILT=$$(date +%Y-%m-%dT%H:%M:%S%:z)"
	# The version number is set to the tag of latest commit
	VERS="$$(git tag --points-at HEAD)"
	if [ -n "$$VERS" ]; then
		[[ ! $$VERS =~ ^v[0-9]+\.[0-9]+\.[0-9]+$$ ]] && echo "illegal VERS=$$VERS " && exit 1
		LDFLAGS="$$LDFLAGS -X $(XFLAG_PATH).VERS=$$VERS"
	fi
	LDFLAGS="$$LDFLAGS -X $(XFLAG_PATH).OS=$$(go env GOOS)/$$(go env GOARCH)"
	go install -ldflags "$$LDFLAGS"

.PHONY: test
test: install
	go vet $(PACKAGES)
	go test -cover $(PACKAGES)

.PHONY: clean
clean: fmt
	go mod verify
	go mod tidy
	go clean -i $(PACKAGES)

.PHONY: fmt
fmt:
	gofmt -w -s $$(find . -name '*.go')

.PHONY: tag
# Tag the latest commit with the VERS environment variable e.g. make tag VERS=v1.0.0
tag:
	[[ ! $$VERS =~ ^v[0-9]+\.[0-9]+\.[0-9]+$$ ]] && echo "error: illegal VERS=$$VERS " && exit 1
	git tag -a -m "$$VERS" $$VERS

.PHONY: push
push: test
	git push -u --tags origin master

DIST_DIR := ./dist

.PHONY: build-dist
# Build executable distributions and compress them to Zip files.
# The VERS environment variable sets version number.
# If VERS is not set the version number defaults to v0.0.0 and version tag
# checks are skipped (v0.0.0 is reserved for testing only).
#
# Normally you want to build from a version-tagged commit, if it
# is not the current head then: stash current changes; temporarily checkout the
# tagged commit; run the build-dist task; revert to previous commit; pop the
# stash e.g. to make a distribution for version v.1.4.0:
#
#   git stash       	# Stash working directory changes
#   git checkout v1.4.0
#   make build-dist
#   git checkout master	# Restore previous commit
#   git stash pop   	# Restore previous working directory changes

# build-dist: clean test validate-docs
build-dist:
	VERS=$${VERS:-v0.0.0}	# v0.0.0 is for testing.
	[[ ! $$VERS =~ ^v[0-9]+\.[0-9]+\.[0-9]+$$ ]] && echo "error: illegal version tag: $$VERS " && exit 1
	if [[ $$VERS != "v0.0.0" ]]; then
		[[ $$(ls $(DIST_DIR)/cryptor-$$VERS* 2>/dev/null | wc -w) -gt 0 ]] && echo "error: built version $$VERS already exists" && exit 1
		headtag="$$(git tag --points-at HEAD)"
		[[ -z "$$headtag" ]] && echo "error: the latest commit has not been tagged" && exit 1
		[[ $$headtag != $$VERS ]] && echo "error: the latest commit tag does not equal $$VERS" && exit 1
		[[ -n "$$(git status --porcelain)" ]] && echo "error: changes in the working directory" && exit 1
	else
		echo "WARNING: no VERS env variable specified, defaulting to v0.0.0 test build"
	fi
	mkdir -p $(DIST_DIR)
	BUILT=$$(date +%Y-%m-%dT%H:%M:%S%:z)
	COMMIT=$$(git rev-parse HEAD)
	BUILD_FLAGS="-X $(XFLAG_PATH).BUILT=$$BUILT -X $(XFLAG_PATH).COMMIT=$$COMMIT -X $(XFLAG_PATH).VERS=$$VERS"
	build () {
		export GOOS=$$1
		export GOARCH=$$2
		LDFLAGS="$$BUILD_FLAGS -X $(XFLAG_PATH).OS=$$GOOS/$$GOARCH"
		LDFLAGS="$$LDFLAGS -s -w"	# Strip symbols to decrease executable size
		NAME=cryptor-$$VERS-$$GOOS-$$GOARCH
		EXE=$$NAME/cryptor
		if [ "$$1" = "windows" ]; then
			EXE=$$EXE.exe
		fi
		ZIP=$$NAME.zip
		rm -f $$ZIP
		rm -rf $$NAME
		mkdir $$NAME
		cp ../LICENSE $$NAME
		cp ../README.md $$NAME/README.txt
		go build -ldflags "$$LDFLAGS" -o $$EXE ..
		zip $$ZIP $$NAME/*
	}
	cd $(DIST_DIR)
	build linux amd64
	build darwin amd64
	build windows amd64
	build windows 386
	sha1sum cryptor-$$VERS*.zip > cryptor-$$VERS-checksums-sha1.txt

.PHONY: release
# Upload release binary distributions for the version assigned to the VERS environment variable e.g. make release VERS=v1.0.0
release:
	REPO=srackham/cryptor
	[[ ! $$VERS =~ ^v[0-9]+\.[0-9]+\.[0-9]+$$ ]] && echo "error: illegal VERS=$$VERS " && exit 1
	upload () {
		export GOOS=$$1
		export GOARCH=$$2
		FILE=cryptor-$$VERS-$$GOOS-$$GOARCH.zip
		gh release upload $$VERS --repo $$REPO $$FILE
	}
	gh release create $$VERS --repo $$REPO --draft --title "Cryptor $$VERS" --notes "Cryptor reports crypto currency portfolio statistics."
	sleep 5	# Wait to avoid "release not found" error.
	cd $(DIST_DIR)
	upload linux amd64
	upload darwin amd64
	upload windows amd64
	upload windows 386
	SUMS=cryptor-$$VERS-checksums-sha1.txt
	gh release upload $$VERS --repo $$REPO $$SUMS