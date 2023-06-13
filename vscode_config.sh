#!/usr/bin/env bash

ROOT_DIR="$(dirname "$(realpath "$0")")"
mkdir -p "${ROOT_DIR}/.vscode"

cat << EOF >"${ROOT_DIR}/.vscode/tasks.json"
{
    "version": "2.0.0",
    "tasks": [
        {
            "label": "Cleanup debug build artifacts",
            "type": "shell",
            "command": "rm",
            "args": [
                "-Rf",
                "\${workspaceFolder}/build/debug"
            ],
            "group": "build",
        },
        {
            "label": "Create macOS .app skeleton",
            "type": "shell",
            "command": "mkdir",
            "args": [
                "-p",
                "\${workspaceFolder}/build/debug/${APP_NAME}.app/Contents/MacOS",
            ],
            "group": "build",
            "dependsOn": "Cleanup debug build artifacts"
        },
        {
            "label": "Embed provisioning profile",
            "type": "shell",
            "command": "cp",
            "args": [
                "${DEBUG_PROVISION_PROFILE_PATH}",
                "\${workspaceFolder}/build/debug/${APP_NAME}.app/Contents/embedded.provisionprofile",
            ],
            "group": "build",
            "dependsOn": "Create macOS .app skeleton"
        },
        {
            "label": "Build debug app",
            "type": "shell",
            "command": "go",
            "args": [
                "build",
                "-gcflags=all=-N -l",
                "-o", "\${workspaceFolder}/build/debug/${APP_NAME}.app/Contents/MacOS/${APP_NAME}",
                "\${workspaceFolder}/cmd"
            ],
            "problemMatcher": ["\$go"],
            "group": {
                "kind": "build",
                "isDefault": true
            },
            "dependsOn": "Embed provisioning profile",
        },
        {
            "label": "Codesign Debug",
            "type": "shell",
            "command": "codesign",
            "args": [
                "--deep",
                "-s",
                "${DEBUG_CERTIFICATE_NAME}",
                "--entitlements",
                "\${workspaceFolder}/vz.entitlements",
                "\${workspaceFolder}/build/debug/${APP_NAME}.app/Contents/MacOS/${APP_NAME}",
            ],
            "group": "build",
            "dependsOn": "Build debug app",
        },
    ],
}
EOF

cat << EOF >"${ROOT_DIR}/.vscode/launch.json"
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "List NICs",
            "type": "go",
            "request": "launch",
            "mode": "exec",
            "program": "\${workspaceFolder}/build/debug/${APP_NAME}.app/Contents/MacOS/${APP_NAME}",
            "args": [
                "netif",
            ],
            "preLaunchTask": "Codesign Debug",
        },
        {
            "name": "Install macOS",
            "type": "go",
            "request": "launch",
            "mode": "exec",
            "program": "\${workspaceFolder}/build/debug/${APP_NAME}.app/Contents/MacOS/${APP_NAME}",
            "args": [
                "install",
            ],
            "preLaunchTask": "Codesign Debug",
        },
        {
            "name": "Run macOS with NIC",
            "type": "go",
            "request": "launch",
            "mode": "exec",
            "program": "\${workspaceFolder}/build/debug/${APP_NAME}.app/Contents/MacOS/${APP_NAME}",
            "args": [
                "run",
                "--nic", "en0",
            ],
            "preLaunchTask": "Codesign Debug",
        },
    ]
}
EOF
