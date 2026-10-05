# Quiver Chat

A small real-time chat. Its interface opens inside the Quiver app: the arrow
serves on a socket that Quiver provides, so no network port is opened for it.

```arrow
schema: "arrow@v0"
metadata:
  name: rabbytesoftware.quiver-chat
  description: Real-time chat whose interface runs inside Quiver, with no port opened.
  version: 27.7.1
  license: MIT
  url: https://github.com/rabbytesoftware/quiver.chat
  quiver: github.com/rabbytesoftware/quiver.chat
  maintainers:
    - name: rabbytesoftware
      url: https://github.com/rabbytesoftware
  tags: [chat, arrow-app]
targets:
  "*":
    lifecycle:
      install:
        - type: fetch
          url:
            default: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-linux-amd64.tar.gz
            linux/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-linux-arm64.tar.gz
            darwin/amd64: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-macos-amd64.tar.gz
            darwin/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-macos-arm64.tar.gz
            windows/amd64: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-windows-amd64.zip
            windows/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-windows-arm64.zip
          to: ./quiver-chat.archive
          title: Downloading Quiver Chat
          timeout: 5m
        - type: extract
          from: ./quiver-chat.archive
          to: ./
          title: Unpacking Quiver Chat
          timeout: 1m
      update:
        - type: fetch
          url:
            default: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-linux-amd64.tar.gz
            linux/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-linux-arm64.tar.gz
            darwin/amd64: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-macos-amd64.tar.gz
            darwin/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-macos-arm64.tar.gz
            windows/amd64: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-windows-amd64.zip
            windows/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/nightly/quiver-chat-windows-arm64.zip
          to: ./quiver-chat.archive
          title: Downloading Quiver Chat
          timeout: 5m
        - type: extract
          from: ./quiver-chat.archive
          to: ./
          title: Unpacking Quiver Chat
          timeout: 1m
      execute:
        - type: ui
          title: Quiver Chat
          listen: [unix]
          path: /
        - type: run
          command:
            default: './quiver-chat-linux-amd64 -listen "${ARROW_UI_LISTEN}"'
            linux/arm64: './quiver-chat-linux-arm64 -listen "${ARROW_UI_LISTEN}"'
            darwin/amd64: './quiver-chat-macos-amd64 -listen "${ARROW_UI_LISTEN}"'
            darwin/arm64: './quiver-chat-macos-arm64 -listen "${ARROW_UI_LISTEN}"'
            windows/amd64: '.\quiver-chat-windows-amd64.exe -listen "${ARROW_UI_LISTEN}"'
            windows/arm64: '.\quiver-chat-windows-arm64.exe -listen "${ARROW_UI_LISTEN}"'
          title: Starting Quiver Chat
      stop:
        - type: signal
          signal: graceful
          timeout: 10s
          exit_on_failure: false
      uninstall: []
```
