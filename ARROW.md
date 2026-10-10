# Quiver Chat

A small real-time chat. The arrow serves the chat on a TCP port (8080 by default),
so anyone who can reach the machine can open it in a browser at
`http://<machine>:8080`. Change the port with the `CHAT_PORT` variable.

```arrow
schema: "arrow@v0"
metadata:
  name: rabbytesoftware.quiver-chat
  description: Real-time chat served on a TCP port, open it in any browser.
  license: MIT
  url: https://github.com/rabbytesoftware/quiver.chat
  quiver: github.com/rabbytesoftware/quiver.chat
  maintainers:
    - name: rabbytesoftware
      url: https://github.com/rabbytesoftware
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
            default: './quiver-chat-linux-amd64 -port "${CHAT_PORT}"'
            linux/arm64: './quiver-chat-linux-arm64 -port "${CHAT_PORT}"'
            darwin/amd64: './quiver-chat-macos-amd64 -port "${CHAT_PORT}"'
            darwin/arm64: './quiver-chat-macos-arm64 -port "${CHAT_PORT}"'
            windows/amd64: '.\quiver-chat-windows-amd64.exe -port "${CHAT_PORT}"'
            windows/arm64: '.\quiver-chat-windows-arm64.exe -port "${CHAT_PORT}"'
          title: Starting Quiver Chat
      stop:
        - type: signal
          signal: graceful
          timeout: 10s
          exit_on_failure: false
      uninstall: []
```
