# 🚀 Quiver Chat

A simple, lightweight chatroom application written in Go that runs as a single executable with no external dependencies.

## ✨ Features

- **Cross-platform**: Runs on Windows, Linux, and macOS
- **Zero dependencies**: Single executable file that includes everything
- **Web-based client**: Access via any modern web browser
- **Real-time messaging**: WebSocket-based instant messaging
- **No user management**: Simple nickname-based chat
- **Modern UI**: Clean, responsive design that works on desktop and mobile
- **Auto-reconnection**: Client automatically reconnects if connection is lost

## 🏗️ Architecture

- **Backend**: Go server with embedded WebSocket support
- **Frontend**: HTML/CSS/JavaScript client embedded in the binary
- **Protocol**: WebSocket for real-time bidirectional communication
- **Deployment**: Single static binary with embedded assets

## 🚀 Quick Start

### Running a Pre-built Binary

1. Download the appropriate binary for your platform from the releases
2. Run the executable:
   ```bash
   # Linux/macOS
   ./quiver-chat-[platform]
   
   # Windows
   quiver-chat-windows-amd64.exe
   ```
3. Open your browser and go to `http://localhost:8080`
4. Choose a nickname and start chatting!

### Building from Source

#### Prerequisites
- Go 1.21 or later

#### Build Steps

1. Clone or download the source code
2. Run the build script:
   ```bash
   # Linux/macOS
   ./build.sh
   
   # Windows
   build.bat
   ```
3. Find your binaries in the `dist/` directory

#### Manual Build
```bash
# For your current platform
go build -o quiver-chat .

# For a specific platform
GOOS=linux GOARCH=amd64 go build -o quiver-chat-linux .
```

## 🖥️ Supported Platforms

- **Linux**: amd64, arm64
- **macOS**: amd64 (Intel), arm64 (Apple Silicon)
- **Windows**: amd64, arm64

## 🎮 Usage

1. **Start the server**: Run the executable on the machine that will host the chat
2. **Connect clients**: Users open their browsers to `http://[server-ip]:8080`
3. **Choose nickname**: Each user picks a display name when joining
4. **Chat away**: Send messages in real-time to all connected users

## 🔧 Configuration

By default, the server runs on port 8080. You can modify the port in the source code if needed.

## 📁 Project Structure

```
quiver-chat/
├── main.go           # Main server application
├── static/
│   └── index.html    # Web client (embedded in binary)
├── build.sh          # Unix build script
├── build.bat         # Windows build script
├── go.mod            # Go module file
├── go.sum            # Go dependencies
└── README.md         # This file
```

## 🛡️ Security Considerations

This is a simple chatroom application designed for trusted networks. Consider these security aspects:

- No authentication or authorization
- Messages are not encrypted in transit (uses plain WebSocket)
- No rate limiting or abuse protection
- No message persistence or logging
- Suitable for local networks, demos, or development

For production use in untrusted environments, consider adding:
- HTTPS/WSS encryption
- Authentication
- Rate limiting
- Input validation and sanitization
- Message moderation

## 🤝 Contributing

This is a simple demonstration project. Feel free to fork and enhance it with additional features like:

- User authentication
- Message history persistence
- File sharing
- Private messaging
- Room management
- SSL/TLS support

## 📄 License

This project is provided as-is for educational and demonstration purposes.

## 🔍 Technical Details

- **Language**: Go 1.21+
- **WebSocket Library**: gorilla/websocket
- **Static Files**: Embedded using Go's embed package
- **Binary Size**: ~8-12MB (varies by platform)
- **Memory Usage**: Minimal, scales with number of connected users
- **Concurrency**: Go routines handle each WebSocket connection

## 🆘 Troubleshooting

**Port already in use**
- Another application is using port 8080
- Kill the process or modify the port in the source code

**Can't connect from other machines**
- Make sure firewall allows incoming connections on port 8080
- Use the server's IP address instead of localhost when connecting from other machines

**Build fails**
- Ensure Go 1.21+ is installed
- Run `go mod tidy` to download dependencies

---

Made with ❤️ in Go 