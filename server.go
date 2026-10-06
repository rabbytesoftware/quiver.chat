package main

import (
	"context"
	"errors"
	"io/fs"
	"net"
	"net/http"
	"os"
	"time"
)

// newHandler serves the embedded frontend and the chat WebSocket.
func newHandler(hub *Hub, static fs.FS) http.Handler {
	mux := http.NewServeMux()
	mux.Handle("/", http.FileServer(http.FS(static)))
	mux.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		handleWebSocket(hub, w, r)
	})
	return mux
}

// listenAddr opens the listener. With a socket path it listens on that unix
// socket (owner-only) and opens no TCP port at all; otherwise it listens on
// every interface at port, as before.
func listenAddr(socket, port string) (net.Listener, error) {
	if socket == "" {
		return net.Listen("tcp", ":"+port)
	}
	// A leftover file from a crashed run would make bind fail.
	if err := os.Remove(socket); err != nil && !errors.Is(err, os.ErrNotExist) {
		return nil, err
	}
	ln, err := net.Listen("unix", socket)
	if err != nil {
		return nil, err
	}
	if err := os.Chmod(socket, 0o600); err != nil {
		_ = ln.Close()
		return nil, err
	}
	return ln, nil
}

// serve runs h on ln until ctx is cancelled, then shuts down. Hijacked
// WebSocket connections are not waited for: Shutdown does not track them.
func serve(ctx context.Context, ln net.Listener, h http.Handler) error {
	srv := &http.Server{Handler: h}
	errc := make(chan error, 1)
	go func() { errc <- srv.Serve(ln) }()

	select {
	case err := <-errc:
		return err
	case <-ctx.Done():
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
		<-errc
		return nil
	}
}
