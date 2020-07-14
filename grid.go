package grid

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"
)

// DefaultAddress of the server is to listen only on the loopback interface.
const DefaultAddress = "127.0.0.1:1111"

// Grid service.
type Grid struct {
	adapter      Adapter // Adapter to the underlying serverless runtime
	address      string  // address upon which to listen
	verbose      bool    // Verbose logging enabled
	version      string  // an externally-provided version, if any.
	onListen     func()  // optional function to run on listen
	httpListener net.Listener
	httpServer   *http.Server
}

type Adapter interface {
	Instances() (int, error)
	SubscriptionManager() SubscriptionManager
	EventManager() EventManager
}

// Subscription manager
type SubscriptionManager interface {
	// Create a new subscription
	Create(string) error
	// Delete a subscription
	Delete(string) error
	// List all active subscriptions
	List() ([]string, error)
}

// EventManager for registering emitted events
type EventManager interface {
	// Create an event registration
	Create(string) error
	// Delete an event registration (must have been created by this service).
	Delete(string) error
	// List all available events
	List() ([]string, error)
}

// New grid service instance.
func New(options ...Option) *Grid {
	g := &Grid{
		address:  DefaultAddress,
		adapter:  noopAdapter{},
		onListen: func() {},
	}
	for _, option := range options {
		option(g)
	}
	return g
}

type Option func(*Grid)

func WithAdapter(a Adapter) Option {
	return func(g *Grid) {
		g.adapter = a
	}
}

func WithVerbose(v bool) Option {
	return func(g *Grid) {
		g.verbose = v
	}
}

func WithVersion(v string) Option {
	return func(g *Grid) {
		g.version = v
	}
}

func WithAddress(a string) Option {
	return func(g *Grid) {
		g.address = a
	}
}

func WithOnListen(f func()) Option {
	return func(g *Grid) {
		g.onListen = f
	}
}

func (g *Grid) Serve(ctx context.Context) (err error) {
	// Listen
	g.httpListener, err = net.Listen("tcp", g.address)
	if err != nil {
		return
	}
	g.onListen()
	if g.verbose {
		fmt.Printf("listening on %v\n", g.httpListener.Addr())
	}

	// Serve
	g.httpServer = &http.Server{
		Handler:        newHandler(g), // see handler.go
		ReadTimeout:    5 * time.Second,
		WriteTimeout:   10 * time.Second,
		IdleTimeout:    30 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	errCh := make(chan error)
	go func() {
		errCh <- g.httpServer.Serve(g.httpListener)
	}()

	// Wait
	// block awaiting os signals or internal errors from serve.
	select {
	case err = <-errCh:
	case <-ctx.Done():
	}

	// Shutdown
	// If there was an error reported from Serve, that takes precidence over a
	// shutdown error.  Print the shutdown error if it exists to log and return
	// the runtime error.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	shutdownErr := g.httpServer.Shutdown(shutdownCtx)
	if shutdownErr != nil {
		if err != nil && err != http.ErrServerClosed {
			fmt.Printf("error shutting down: %v\n", err)
			return err
		}
		if shutdownErr != http.ErrServerClosed {
			return shutdownErr
		}
	}
	return nil
}

// Addr is the current listening address, or nil if not yet listening
func (g *Grid) Addr() net.Addr {
	if g.httpListener == nil {
		return nil
	}
	return g.httpListener.Addr()
}

// noop implementations for zero-value struct in testing

type noopAdapter struct{}

func (n noopAdapter) Instances() (int, error)                  { return 0, nil }
func (n noopAdapter) SubscriptionManager() SubscriptionManager { return noopSubscriptionManager{} }
func (n noopAdapter) EventManager() EventManager               { return noopEventManager{} }

type noopSubscriptionManager struct{}

func (n noopSubscriptionManager) Create(string) error     { return nil }
func (n noopSubscriptionManager) Delete(string) error     { return nil }
func (n noopSubscriptionManager) List() ([]string, error) { return []string{}, nil }

type noopEventManager struct{}

func (n noopEventManager) Create(string) error     { return nil }
func (n noopEventManager) Delete(string) error     { return nil }
func (n noopEventManager) List() ([]string, error) { return []string{}, nil }
