package grid

import (
	"net/http/httptest"
	"testing"
)

func TestVersionHandler(t *testing.T) {
	g := New(
		WithAddress("127.0.0.1:"),
	)

	req := httptest.NewRequest("GET", "http://127.0.0.1/version", nil)
	w := httptest.NewRecorder()
	g.handleVersion(w, req)

	resp := w.Result()

	if resp.StatusCode != 200 {
		t.Error("Expected status code 200")
	}
}

// TODO: handler unit tests using httptest
