SHELL := /bin/bash
VERSION := $(shell cat VERSION 2>/dev/null || echo 0.1.0)

.PHONY: setup packages iso iso-dev test smoke release clean
setup:
	bash iso/scripts/setup-build-env.sh
packages:
	bash iso/scripts/build-packages.sh
iso: packages
	bash iso/scripts/build.sh production
iso-dev: packages
	bash iso/scripts/build.sh dev
test:
	bash iso/scripts/test-qemu.sh
smoke:
	bash iso/scripts/smoke-test.sh
release: iso
	bash iso/scripts/release.sh
clean:
	rm -rf build dist
