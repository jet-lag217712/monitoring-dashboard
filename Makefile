# Appliance development (Equate-Appliance VM) and production release targets.
.DEFAULT_GOAL := help

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo dev)
ARCH ?= $(shell uname -m | sed 's/x86_64/amd64/;s/aarch64/arm64/')
RUNTIME_SCRIPTS := deployments/production/appliance/scripts
RELEASE_SCRIPTS := appliance/scripts

.PHONY: help

help:
	@echo "Develop (Equate-Appliance VM from source):"
	@echo "  make appliance-bundle [ARCH=arm64|amd64] [VERSION=x.y.z]"
	@echo "  make appliance-stage HOST=vm.local [USER=$$USER] ARCH=... VERSION=..."
	@echo "  make appliance-configure BUNDLE=/tmp/equate-staging/bundle VERSION=..."
	@echo "  make appliance-upgrade BUNDLE=/tmp/equate-staging/bundle VERSION=... [CANARY=1]"
	@echo "  make appliance-verify"
	@echo "  make test [QUICK=1]"
	@echo
	@echo "Lab (GNS3 + appliance VM):"
	@echo "  make lab-validate"
	@echo "  make lab-appliance-setup"
	@echo "  make lab-appliance-upgrade OLD_VERSION=... NEW_VERSION=... BUNDLE=..."
	@echo "  make lab-gns3-bridge"
	@echo "  make lab-smoke"
	@echo "  make lab-mqtt-outage-drill"
	@echo
	@echo "Release:"
	@echo "  make appliance-bundle-all [VERSION=x.y.z]"
	@echo "  make appliance-package [ARCH=...] [VERSION=...]   # .eqa (+ sign if EQUATE_UPDATE_SIGNING_KEY set)"
	@echo "  make appliance-ova-amd64-ci [VERSION=x.y.z]"
	@echo "  make appliance-prepare-ova"
	@echo "  make appliance-publish-azure STORAGE_ACCOUNT=... VERSION=... ARCH=... [CHANNEL=stable]"
	@echo "  make appliance-generate-keys [KEYS_DIR=...]"
	@echo "  make appliance-setup-azure-channel STORAGE_ACCOUNT=..."
	@echo
	@echo "Engineer helpers:"
	@echo "  make db-migrate DATABASE_URL=..."
	@echo "  make db-bootstrap-roles DATABASE_URL=..."
	@echo "  make mqtt-dev-certs"

include make/develop.mk
include make/release.mk
