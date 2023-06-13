# go-macos-virtualization

A simple application to create and launch macOS VMs written in golang.

Relying heavily on [Code-Hex/vz](https://github.com/Code-Hex/vz) with the usage of VMNet.

## Getting started

1. Open `Makefile` and populate empty variables:

   - `APP_NAME` - name of the MacOS app bundle. Doesn't really matter

   - `DEBUG_CERTIFICATE_NAME` - Mac Developer certificate to sign debug binary with. Used by vscode launch configurations

   - `DEBUG_PROVISION_PROFILE_PATH` - location of the provision profile generated in the dev portal with VMNet and VM Networking capabilities associated with the certificate above. Used by vscode launch configurations

   - `RELEASE_CERTIFICATE_NAME` (optional) - Developer ID Application certificate to build release binary with. Used by `make release` command

   - `RELEASE_PROVISION_PROFILE_PATH` (optional) - location of the provision profile generated in the dev portal with VMNet and VM Networking capabilities associated with the certificate above. Used by `make release` command

2. Execute `make vscode` to generate vscode launch configurations and tasks for development purposes. Or `make release` to generate release binary with prod certificate

## VMNet and VM Networking capabilities

These 2 capabilities require Apple approval. In order to get it. Follow [this Apple Forum thread](https://developer.apple.com/forums/thread/656411) on how to request it.

Afterwards, simply generate yourself Mac Development certificate, App ID with those capabilities, and provision profile. Input those in Makefile and enjoy.
