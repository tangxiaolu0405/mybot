package mcp

import (
	"context"
	"testing"
)

func TestTryCall_unknownToolNotPanic(t *testing.T) {
	mgr := &Manager{
		clients: make(map[string]*stdioClient),
		routes:  make(map[string]*toolRoute),
	}
	out, err, ok := mgr.TryCall(context.Background(), "read_file", `{}`)
	if ok || err != nil || out != "" {
		t.Fatalf("expected noop, got ok=%v err=%v out=%q", ok, err, out)
	}
}

func TestTryCall_missingRouteNotPanic(t *testing.T) {
	mgr := &Manager{
		clients: map[string]*stdioClient{"browser": nil},
		routes:  map[string]*toolRoute{"browser_navigate": {serverName: "browser", toolName: "browser_navigate"}},
	}
	// route exists but client nil -> false, no panic
	_, _, ok := mgr.TryCall(context.Background(), "browser_navigate", `{}`)
	if ok {
		t.Fatal("expected false when client nil")
	}
}
