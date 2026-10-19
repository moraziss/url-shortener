package urlshort

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

// fallback answers 418 so tests can tell it apart from a redirect.
func fallback() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusTeapot)
	})
}

func serve(h http.Handler, path string) *httptest.ResponseRecorder {
	rec := httptest.NewRecorder()
	h.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, path, nil))
	return rec
}

func TestMapHandlerRedirectsKnownPath(t *testing.T) {
	h := MapHandler(map[string]string{"/go": "https://go.dev"}, fallback())

	rec := serve(h, "/go")

	if rec.Code != http.StatusFound {
		t.Fatalf("status = %d, want %d", rec.Code, http.StatusFound)
	}
	if got := rec.Header().Get("Location"); got != "https://go.dev" {
		t.Errorf("Location = %q, want %q", got, "https://go.dev")
	}
}

func TestMapHandlerFallsBackForUnknownPath(t *testing.T) {
	h := MapHandler(map[string]string{"/go": "https://go.dev"}, fallback())

	rec := serve(h, "/missing")

	if rec.Code != http.StatusTeapot {
		t.Errorf("status = %d, want fallback status %d", rec.Code, http.StatusTeapot)
	}
}

func TestYAMLHandlerRedirects(t *testing.T) {
	data := []byte("- path: /docs\n  url: https://example.com/docs\n- path: /blog\n  url: https://example.com/blog\n")

	h, err := YAMLHandler(data, fallback())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	for path, want := range map[string]string{
		"/docs": "https://example.com/docs",
		"/blog": "https://example.com/blog",
	} {
		rec := serve(h, path)
		if rec.Code != http.StatusFound {
			t.Errorf("%s: status = %d, want %d", path, rec.Code, http.StatusFound)
		}
		if got := rec.Header().Get("Location"); got != want {
			t.Errorf("%s: Location = %q, want %q", path, got, want)
		}
	}
}

func TestYAMLHandlerFallsBackForUnknownPath(t *testing.T) {
	h, err := YAMLHandler([]byte("- path: /docs\n  url: https://example.com/docs\n"), fallback())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if rec := serve(h, "/nope"); rec.Code != http.StatusTeapot {
		t.Errorf("status = %d, want fallback status %d", rec.Code, http.StatusTeapot)
	}
}

func TestYAMLHandlerRejectsInvalidYAML(t *testing.T) {
	if _, err := YAMLHandler([]byte("path: [unclosed"), fallback()); err == nil {
		t.Fatal("expected an error for invalid YAML, got nil")
	}
}
