package grid_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/boson-project/grid"
)

func StartGrid(t *testing.T) (g *grid.Grid, err error) {
	t.Helper()

	listening := make(chan bool)
	errCh := make(chan error)

	g = grid.New(
		grid.WithAddress("127.0.0.1:"),                  // OS-chosen port
		grid.WithOnListen(func() { listening <- true }), // signal start
	)

	go func() {
		if err := g.Serve(context.Background()); err != nil {
			errCh <- err
		}
	}()

	select {
	case err = <-errCh:
	case <-listening:
	}
	return g, err
}

// TestCancel ensures the service starts and stops without error with all defaults using
// a cancelable context.
func TestStart(t *testing.T) {
	// A context which, when canceled, triggers a graceful shutdown of the server.
	ctx, cancel := context.WithCancel(context.Background())

	// Grid instance which immediately triggers a shutdown when listening.
	g := grid.New(grid.WithOnListen(cancel))

	// Serve, which should return without error on graceful shutdown.
	if err := g.Serve(ctx); err != nil {
		t.Fatal(err)
	}
}

// TestVersion ensures that the /v1/version endpoint returns the version structure.
func TestVersion(t *testing.T) {
	g, err := StartGrid(t)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get("http://" + g.Addr().String() + "/v1/version")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected HTTP 200, got %v", res.StatusCode)
	}
}

func TestEventsEndpointAvailable(t *testing.T) {
	g, err := StartGrid(t)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get("http://" + g.Addr().String() + "/v1/events")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected HTTP 200, got %v", res.StatusCode)
	}
}

func TestSubscriptionsEndpointAvailable(t *testing.T) {
	g, err := StartGrid(t)
	if err != nil {
		t.Fatal(err)
	}

	res, err := http.Get("http://" + g.Addr().String() + "/v1/subscriptions")
	if err != nil {
		t.Fatal(err)
	}
	if res.StatusCode != 200 {
		t.Fatalf("Expected HTTP 200, got %v", res.StatusCode)
	}
}

// See handlers_test.go for individual endpoint unit tests.
