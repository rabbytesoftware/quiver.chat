# Quiver Chat 🚀

A modern, real-time chat application built with **Next.js + shadcn/ui** frontend and **Go** backend. All frontend assets are embedded in the Go binary for zero-dependency deployment.

## Features

- 🎨 **Modern UI**: Built with Next.js, TypeScript, and shadcn/ui components
- 💬 **Real-time Chat**: WebSocket-based communication
- 🔄 **Zero Dependencies**: Single executable with embedded frontend
- 🌐 **Cross-platform**: Works on Windows, macOS, and Linux
- 📱 **Responsive Design**: Mobile-friendly interface
- 🎯 **Type-safe**: Full TypeScript support

## Architecture

- **Frontend**: Next.js with TypeScript and shadcn/ui components
- **Backend**: Go with WebSocket support using Gorilla WebSocket
- **Build**: Next.js static export embedded in Go binary using `embed.FS`

## Quick Start

1. **Download** the latest binary for your platform from the releases
2. **Run** the executable:
   ```bash
   ./quiver-chat-[platform]
   ```
3. **Open** your browser and navigate to `http://localhost:8080`
4. **Choose** a nickname and start chatting!

## Development

### Prerequisites

- Go 1.24.2 or higher
- Node.js 18+ and npm
- Git

### Setup

1. **Clone** the repository:

   ```bash
   git clone <repository-url>
   cd quiver.chat
   ```
2. **Install** frontend dependencies:

   ```bash
   cd frontend
   npm install
   ```
3. **Build** the application:

   ```bash
   cd ..
   make all
   ```

### Development Workflow

1. **Frontend Development**:

   ```bash
   cd frontend
   npm run dev
   ```
2. **Backend Development**:

   ```bash
   go run main.go                      # TCP, http://localhost:8080
   go run main.go -listen /tmp/chat.sock   # unix socket only, no TCP port
   ```
3. **Build for Production**:

   ```bash
   make all
   ```

## Project Structure

```
quiver.chat/
├── main.go                 # Go backend server
├── go.mod                  # Go dependencies
├── build.sh                # Build script
├── frontend/               # Next.js frontend
│   ├── src/
│   │   ├── components/     # React components
│   │   ├── hooks/          # Custom hooks
│   │   ├── types/          # TypeScript types
│   │   └── app/            # Next.js app directory
│   └── package.json        # Frontend dependencies
└── dist/                   # Built binaries
```

## How It Works

1. **Build Process**: The build script first compiles the Next.js frontend into static files
2. **Embedding**: Go's `embed.FS` embeds the static files into the binary
3. **Serving**: The Go server serves the embedded files and handles WebSocket connections
4. **Real-time**: WebSocket connections enable real-time message exchange

## API Endpoints

- `GET /`: Serves the Next.js application
- `WS /ws?username=<name>`: WebSocket endpoint for real-time chat

## Message Format

```json
{
  "type": "message|system",
  "username": "string",
  "content": "string",
  "time": "HH:MM"
}
```

## Build Targets

The build script creates executables for:

- Linux (amd64, arm64)
- macOS (amd64, arm64)
- Windows (amd64, arm64)

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Run tests and build
5. Submit a pull request

## License

This project is open source and available under the [MIT License](LICENSE).

---

**No external dependencies required - just run and go!** 🎉
