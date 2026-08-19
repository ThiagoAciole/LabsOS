SHELL := /bin/bash
VERSION := $(shell cat VERSION 2>/dev/null || echo 0.1.0)

.PHONY: setup packages iso iso-dev test smoke release clean
setup:
	bash packaging/scripts/setup-build-env.sh
packages:
	bash packaging/scripts/build-packages.sh
iso: packages
	bash packaging/scripts/build.sh production
iso-dev: packages
	bash packaging/scripts/build.sh dev
test:
	bash packaging/scripts/test-qemu.sh
smoke:
	bash packaging/scripts/smoke-test.sh
release: iso
	bash packaging/scripts/release.sh
clean:
	rm -rf build dist
