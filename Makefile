.PHONY: release vscode

APP_NAME := "MacOSVM"

DEBUG_CERTIFICATE_NAME := ""
DEBUG_PROVISION_PROFILE_PATH := ""

RELEASE_CERTIFICATE_NAME := ""
RELEASE_PROVISION_PROFILE_PATH := ""

ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))

release:
	@echo "Building Release app..."

	# Prepare app skeleton
	rm -Rf "$(ROOT_DIR)/build/release"
	mkdir -p "$(ROOT_DIR)/build/release/$(APP_NAME).app/Contents/MacOS"

	# Embed provision profile
	cp $(RELEASE_PROVISION_PROFILE_PATH) "$(ROOT_DIR)/build/release/$(APP_NAME).app/Contents/embedded.provisionprofile"

	# Build app
	go build -o "$(ROOT_DIR)/build/release/$(APP_NAME).app/Contents/MacOS/$(APP_NAME)" "$(ROOT_DIR)/cmd"

	# Codesign final binary with proper entitlements
	codesign --deep -s $(RELEASE_CERTIFICATE_NAME) --entitlements "$(ROOT_DIR)/vz.entitlements" "$(ROOT_DIR)/build/release/$(APP_NAME).app/Contents/MacOS/$(APP_NAME)"

vscode:
	APP_NAME=$(APP_NAME) \
		DEBUG_CERTIFICATE_NAME=$(DEBUG_CERTIFICATE_NAME) \
		DEBUG_PROVISION_PROFILE_PATH=$(DEBUG_PROVISION_PROFILE_PATH) \
		/bin/bash "$(ROOT_DIR)/vscode_config.sh"
