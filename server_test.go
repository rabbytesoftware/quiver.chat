package main

import (
	"context"
	"io"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
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
	base := "/tmp"
	if runtime.GOOS == "windows" {
		base = ""
	}
	dir, err := os.MkdirTemp(base, "qc")
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

	if runtime.GOOS == "windows" {
		t.Skip("windows has no unix permission bits; the socket sits in the user's profile")
	}
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

func TestListeners_SocketAloneOpensNoTCPPort(t *testing.T) {
	lns, err := listeners(shortSocket(t), "0", false)
	if err != nil {
		t.Fatal(err)
	}
	defer closeAll(lns)
	if len(lns) != 1 || lns[0].Addr().Network() != "unix" {
		t.Fatalf("want only the unix socket, got %d listeners", len(lns))
	}
}

func TestListeners_ExplicitPortAddsTCPNextToTheSocket(t *testing.T) {
	lns, err := listeners(shortSocket(t), "0", true)
	if err != nil {
		t.Fatal(err)
	}
	defer closeAll(lns)
	if len(lns) != 2 || lns[0].Addr().Network() != "unix" || lns[1].Addr().Network() != "tcp" {
		t.Fatalf("want unix socket then tcp, got %d listeners", len(lns))
	}
}

func TestListeners_TCPOnlyWithoutASocket(t *testing.T) {
	lns, err := listeners("", "0", false)
	if err != nil {
		t.Fatal(err)
	}
	defer closeAll(lns)
	if len(lns) != 1 || lns[0].Addr().Network() != "tcp" {
		t.Fatalf("want only tcp, got %d listeners", len(lns))
	}
}

func TestListeners_FailedTCPBindReleasesTheSocket(t *testing.T) {
	busy, err := net.Listen("tcp", ":0")
	if err != nil {
		t.Fatal(err)
	}
	defer busy.Close()
	_, port, _ := net.SplitHostPort(busy.Addr().String())

	sock := shortSocket(t)
	if _, err := listeners(sock, port, true); err == nil {
		t.Fatal("binding a busy port must fail")
	}
	// The socket opened first must not stay behind, or a retry would see a dead file.
	if c, err := net.DialTimeout("unix", sock, 200*time.Millisecond); err == nil {
		c.Close()
		t.Fatal("socket still accepting after a failed start")
	}
}

// startBoth serves one room on a unix socket and a TCP port, as the arrow does.
func startBoth(t *testing.T) (sock, tcpAddr string, cancel context.CancelFunc) {
	t.Helper()
	hub := newHub()
	go hub.run()
	sock = shortSocket(t)
	lns, err := listeners(sock, "0", true)
	if err != nil {
		t.Fatal(err)
	}
	tcpAddr = lns[1].Addr().String()
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	static := fstest.MapFS{"index.html": {Data: []byte("<html>chat</html>")}}
	go func() { done <- serveAll(ctx, lns, newHandler(hub, static)) }()
	t.Cleanup(func() {
		cancel()
		select {
		case <-done:
		case <-time.After(5 * time.Second):
			t.Error("serveAll did not stop after cancel")
		}
	})
	return sock, tcpAddr, cancel
}

func TestServeAll_OneRoomOverTheSocketAndTCP(t *testing.T) {
	sock, tcpAddr, _ := startBoth(t)

	// The page is served on both transports.
	resp, err := unixClient(sock).Get("http://chat/")
	if err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
	tcpResp, err := http.Get("http://" + tcpAddr + "/")
	if err != nil {
		t.Fatal(err)
	}
	tcpResp.Body.Close()
	if resp.StatusCode != http.StatusOK || tcpResp.StatusCode != http.StatusOK {
		t.Fatalf("socket %d, tcp %d", resp.StatusCode, tcpResp.StatusCode)
	}

	// alice sits in the Quiver window (socket), bob in a browser (TCP): one room.
	unixDialer := websocket.Dialer{NetDialContext: func(ctx context.Context, _, _ string) (net.Conn, error) {
		var d net.Dialer
		return d.DialContext(ctx, "unix", sock)
	}}
	alice, _, err := unixDialer.Dial("ws://chat/ws?username=alice", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer alice.Close()
	bob, _, err := websocket.DefaultDialer.Dial("ws://"+tcpAddr+"/ws?username=bob", nil)
	if err != nil {
		t.Fatal(err)
	}
	defer bob.Close()

	time.Sleep(100 * time.Millisecond)
	if err := alice.WriteMessage(websocket.TextMessage, []byte(`{"type":"message","content":"hello from the window"}`)); err != nil {
		t.Fatal(err)
	}
	_ = bob.SetReadDeadline(time.Now().Add(3 * time.Second))
	for {
		_, msg, err := bob.ReadMessage()
		if err != nil {
			t.Fatalf("the browser never saw the message sent from the socket client: %v", err)
		}
		if strings.Contains(string(msg), "hello from the window") {
			break
		}
	}

	// And the other way round.
	if err := bob.WriteMessage(websocket.TextMessage, []byte(`{"type":"message","content":"hello from the browser"}`)); err != nil {
		t.Fatal(err)
	}
	_ = alice.SetReadDeadline(time.Now().Add(3 * time.Second))
	for {
		_, msg, err := alice.ReadMessage()
		if err != nil {
			t.Fatalf("the Quiver window never saw the message sent from TCP: %v", err)
		}
		if strings.Contains(string(msg), "hello from the browser") {
			return
		}
	}
}

func TestServeAll_CancelStopsBothListeners(t *testing.T) {
	sock, tcpAddr, cancel := startBoth(t)
	cancel()

	deadline := time.Now().Add(3 * time.Second)
	for time.Now().Before(deadline) {
		_, unixErr := net.DialTimeout("unix", sock, 100*time.Millisecond)
		_, tcpErr := net.DialTimeout("tcp", tcpAddr, 100*time.Millisecond)
		if unixErr != nil && tcpErr != nil {
			return
		}
		time.Sleep(50 * time.Millisecond)
	}
	t.Fatal("a listener is still accepting after cancel")
}

func closeAll(lns []net.Listener) {
	for _, ln := range lns {
		_ = ln.Close()
	}
}
