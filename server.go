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

// listenAddr opens one listener. With a socket path it listens on that unix
// socket (owner-only) and opens no TCP port at all; otherwise it listens on
// every interface at port.
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

// listeners opens everything the flags ask for. A socket alone keeps the
// private, no-port mode that existing arrows rely on. A port that was given
// explicitly next to the socket opens TCP as well, so the Quiver window and
// browsers on the network reach the same chat room.
func listeners(socket, port string, portGiven bool) ([]net.Listener, error) {
	var lns []net.Listener
	closeAll := func() {
		for _, ln := range lns {
			_ = ln.Close()
		}
	}
	if socket != "" {
		ln, err := listenAddr(socket, "")
		if err != nil {
			return nil, err
		}
		lns = append(lns, ln)
	}
	if socket == "" || portGiven {
		ln, err := listenAddr("", port)
		if err != nil {
			closeAll()
			return nil, err
		}
		lns = append(lns, ln)
	}
	return lns, nil
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

// serveAll runs h on every listener until ctx is cancelled. If one of them
// fails the others are stopped too, and the first error is returned.
func serveAll(ctx context.Context, lns []net.Listener, h http.Handler) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	errc := make(chan error, len(lns))
	for _, ln := range lns {
		go func(ln net.Listener) { errc <- serve(ctx, ln, h) }(ln)
	}
	var first error
	for range lns {
		if err := <-errc; err != nil && first == nil {
			first = err
			cancel()
		}
	}
	return first
}
