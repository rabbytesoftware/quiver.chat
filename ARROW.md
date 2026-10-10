# Quiver Chat

A small real-time chat. Its interface opens inside the Quiver app, and the same
chat room is also served on a TCP port (8080 by default), so anyone who can reach
the machine can join from a browser at `http://<machine>:8080`. Change the port
with the `CHAT_PORT` variable.

```arrow
schema: "arrow@v0"
metadata:
  name: rabbytesoftware.quiver-chat
  description: Real-time chat that opens inside Quiver and is also reachable on a TCP port.
  license: MIT
  url: https://github.com/rabbytesoftware/quiver.chat
  quiver: github.com/rabbytesoftware/quiver.chat
  maintainers:
    - name: rabbytesoftware
      url: https://github.com/rabbytesoftware
  media:
    icon: "https://raw.githubusercontent.com/rabbytesoftware/quiver.core/develop/docs/quiver-icon.svg"
    banner: "https://raw.githubusercontent.com/rabbytesoftware/quiver.core/develop/docs/quiver-banner.svg"
  tags: [chat, arrow-app]

netbridge:
  - name: CHAT_PORT
    default: 8080
    protocol: tcp
    required: true

targets:
  "*":
    lifecycle:
      install:
        - type: fetch
          url:
            default: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-linux-amd64.tar.gz
            linux/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-linux-arm64.tar.gz
            darwin/amd64: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-macos-amd64.tar.gz
            darwin/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-macos-arm64.tar.gz
            windows/amd64: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-windows-amd64.zip
            windows/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-windows-arm64.zip
          to: ${INSTALL_PATH}/quiver-chat.archive
          title: Downloading Quiver Chat
          timeout: 5m
        - type: extract
          from: ${INSTALL_PATH}/quiver-chat.archive
          to: ${INSTALL_PATH}
          title: Unpacking Quiver Chat
          timeout: 1m
      update:
        - type: fetch
          url:
            default: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-linux-amd64.tar.gz
            linux/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-linux-arm64.tar.gz
            darwin/amd64: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-macos-amd64.tar.gz
            darwin/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-macos-arm64.tar.gz
            windows/amd64: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-windows-amd64.zip
            windows/arm64: https://github.com/rabbytesoftware/quiver.chat/releases/download/${REF}/quiver-chat-windows-arm64.zip
          to: ${INSTALL_PATH}/quiver-chat.archive
          title: Downloading Quiver Chat
          timeout: 5m
        - type: extract
          from: ${INSTALL_PATH}/quiver-chat.archive
          to: ${INSTALL_PATH}
          title: Unpacking Quiver Chat
          timeout: 1m
      execute:
        - type: run
          command:
            default: './quiver-chat-linux-amd64 -listen "${ARROW_UI_LISTEN}" -port "${CHAT_PORT}"'
            linux/arm64: './quiver-chat-linux-arm64 -listen "${ARROW_UI_LISTEN}" -port "${CHAT_PORT}"'
            darwin/amd64: './quiver-chat-macos-amd64 -listen "${ARROW_UI_LISTEN}" -port "${CHAT_PORT}"'
            darwin/arm64: './quiver-chat-macos-arm64 -listen "${ARROW_UI_LISTEN}" -port "${CHAT_PORT}"'
            windows/amd64: '.\quiver-chat-windows-amd64.exe -listen "${ARROW_UI_LISTEN}" -port "${CHAT_PORT}"'
            windows/arm64: '.\quiver-chat-windows-arm64.exe -listen "${ARROW_UI_LISTEN}" -port "${CHAT_PORT}"'
          title: Starting Quiver Chat
          ui:
            title: Quiver Chat
            path: /
            listen: [unix]
      stop:
        - type: signal
          signal: graceful
          timeout: 10s
          exit_on_failure: false
      uninstall: []
```
