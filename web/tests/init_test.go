package web_test

import (
	"authosaurous/pkg/testutils"
	"context"
	"os"
	"testing"
)

// Global test setup
func TestMain(m *testing.M) {

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	// Initialize test container for PostgreSQL
	server := testutils.StartTestServer(ctx)
	// Run all tests
	code := m.Run()
  // cleanups
	server.StopTestServer(ctx)
	os.Exit(code)
}

