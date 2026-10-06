package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/fstest"
	"time"

	"github.com/gorilla/websocket"
)

// shortSocket returns a socket path short enough for a sockaddr_un (macOS
// temp dirs are too long).
func shortSocket(t *testing.T) string {
	t.Helper()
	dir, err := os.MkdirTemp("/tmp", "qc")
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = os.RemoveAll(dir) })
	return filepath.Join(dir, "c.sock")
}

func unixClient(sock string) *http.Client {
	return &http.Client{Transport: &http.Transport{
		DialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
			var d net.Dialer
			return d.DialContext(ctx, "unix", sock)
		},
	}}
}

func startServer(t *testing.T, sock string) context.CancelFunc {
	t.Helper()
	hub := newHub()
	go hub.run()
	ln, err := listenAddr(sock, "")
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	static := fstest.MapFS{"index.html": {Data: []byte("<html>chat</html>")}}
	go func() { done <- serve(ctx, ln, newHandler(hub, static)) }()
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("serve did not stop after cancel")
		}
	})
	return cancel
}

func TestListenAddr_UnixSocketIsPrivate(t *testing.T) {
	sock := shortSocket(t)
	ln, err := listenAddr(sock, "")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()

	info, err := os.Stat(sock)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm&0o077 != 0 {
		t.Fatalf("socket must not be accessible to group/other, got %o", perm)
	}
}

func TestListenAddr_ReplacesAStaleSocketFile(t *testing.T) {
	sock := shortSocket(t)
	if err := os.WriteFile(sock, nil, 0o600); err != nil {
		t.Fatal(err)
	}
	ln, err := listenAddr(sock, "")
	if err != nil {
		t.Fatalf("a stale file must not block listening: %v", err)
	}
	ln.Close()
}

func TestListenAddr_TCPWhenNoSocketGiven(t *testing.T) {
	ln, err := listenAddr("", "0")
	if err != nil {
		t.Fatal(err)
	}
	defer ln.Close()
	if ln.Addr().Network() != "tcp" {
		t.Fatalf("want tcp, got %s", ln.Addr().Network())
	}
}

func TestServe_ServesTheFrontendOverTheSocket(t *testing.T) {
	sock := shortSocket(t)
	startServer(t, sock)

	resp, err := unixClient(sock).Get("http://chat/")
	if err != nil {
		t.Fatal(err)
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "chat") {
		t.Fatalf("got %d %q", resp.StatusCode, body)
	}
}

func TestServe_StopsOnCancelAndFreesTheSocket(t *testing.T) {
	sock := shortSocket(t)
	cancel := startServer(t, sock)
	cancel()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		if _, err := net.DialTimeout("unix", sock, 100*time.Millisecond); err != nil {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("server still accepting connections after cancel")
}

func TestServe_WebSocketBroadcastBetweenTwoClients(t *testing.T) {
	sock := shortSocket(t)
	startServer(t, sock)

	dialer := websocket.Dialer{NetDialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", sock)
	}}
	a, _, err := dialer.Dial("ws://chat/ws?username=alice", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer a.Close()
	b, _, err := dialer.Dial("ws://chat/ws?username=bob", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer b.Close()

	// Give the hub a moment to register both clients, then send as alice.
	time.Sleep(100 * time.Millisecond)
	if err := a.WriteMessage(websocket.TextMessage, []byte(`{"type":"message","content":"hello bob"}`)); err != nil {
		t.Fatal(err)
	}

	_ = b.SetReadDeadline(time.Now().Add(3 * time.Second))
	for {
		_, msg, err := b.ReadMessage()
		if err != nil {
			t.Fatalf("bob never received the broadcast: %v", err)
		}
		if strings.Contains(string(msg), "hello bob") {
			return
		}
	}
}
