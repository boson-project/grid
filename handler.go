package grid

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/gorilla/mux"
)

func newHandler(g *Grid) http.Handler {
	h := mux.NewRouter()
	h.HandleFunc("/v1/version", g.handleVersion) // see handlers.go
	h.HandleFunc("/v1/events", g.handleEvents)
	h.HandleFunc("/v1/subscriptions", g.handleSubscriptions)
	return h
}

// handleVersion returns the version number provided, if any, JSON encoded.
func (g *Grid) handleVersion(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	if err := json.NewEncoder(w).Encode(g.version); err != nil {
		fmt.Fprintln(os.Stderr, fmt.Errorf("error encoding version: %v", err))
	}
}

// handleEvents returns the list of available events
func (g *Grid) handleEvents(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch r.Method {
	case "GET":
		var events []string
		if events, err = g.adapter.EventManager().List(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error while encoding error: %v", err))
			}
			return
		}
		if err := json.NewEncoder(w).Encode(events); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("error encoding version: %v", err))
		}
	case "POST":
		w.WriteHeader(http.StatusNotImplemented)
		err = json.NewEncoder(w).Encode("Creation of events not implemented.")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	case "DELETE":
		w.WriteHeader(http.StatusNotImplemented)
		err = json.NewEncoder(w).Encode("Deleteion of events not implemented.")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		err = json.NewEncoder(w).Encode("Method not supported.")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}

// handleSubscriptions returns the list of available subscriptions
func (g *Grid) handleSubscriptions(w http.ResponseWriter, r *http.Request) {
	var err error
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	switch r.Method {
	case "GET":
		var subscriptions []string
		if subscriptions, err = g.adapter.SubscriptionManager().List(); err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			err = json.NewEncoder(w).Encode(err)
			if err != nil {
				fmt.Fprintln(os.Stderr, fmt.Errorf("error while encoding error: %v", err))
			}
			return
		}
		if err := json.NewEncoder(w).Encode(subscriptions); err != nil {
			fmt.Fprintln(os.Stderr, fmt.Errorf("error encoding version: %v", err))
		}
	case "POST":
		w.WriteHeader(http.StatusNotImplemented)
		err = json.NewEncoder(w).Encode("Creation of subscriptions not implemented.")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	case "DELETE":
		w.WriteHeader(http.StatusNotImplemented)
		err = json.NewEncoder(w).Encode("Deleteion of subscriptions not implemented.")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
		err = json.NewEncoder(w).Encode("Method not supported.")
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
		}
	}
}
