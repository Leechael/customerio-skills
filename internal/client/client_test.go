package client

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("missing token", func(t *testing.T) {
		t.Setenv("CUSTOMERIO_API_TOKEN", "")
		_, err := New("us")
		if err == nil {
			t.Fatal("expected error for missing token")
		}
	})

	t.Run("us region", func(t *testing.T) {
		t.Setenv("CUSTOMERIO_API_TOKEN", "tok")
		c, err := New("us")
		if err != nil {
			t.Fatal(err)
		}
		if c.BaseURL != "https://api.customer.io" {
			t.Fatalf("got %s", c.BaseURL)
		}
	})

	t.Run("eu region", func(t *testing.T) {
		t.Setenv("CUSTOMERIO_API_TOKEN", "tok")
		c, err := New("eu")
		if err != nil {
			t.Fatal(err)
		}
		if c.BaseURL != "https://api-eu.customer.io" {
			t.Fatalf("got %s", c.BaseURL)
		}
	})

	t.Run("unknown region defaults to us", func(t *testing.T) {
		t.Setenv("CUSTOMERIO_API_TOKEN", "tok")
		c, err := New("asia")
		if err != nil {
			t.Fatal(err)
		}
		if c.BaseURL != "https://api.customer.io" {
			t.Fatalf("got %s", c.BaseURL)
		}
	})
}

func testServer(t *testing.T, handler http.HandlerFunc) *Client {
	t.Helper()
	srv := httptest.NewServer(handler)
	t.Cleanup(srv.Close)
	return &Client{BaseURL: srv.URL, Token: "test-token", HTTPClient: srv.Client()}
}

func TestGet(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		c := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("method = %s", r.Method)
			}
			if r.Header.Get("Authorization") != "Bearer test-token" {
				t.Errorf("auth = %s", r.Header.Get("Authorization"))
			}
			w.Write([]byte(`{"ok":true}`))
		})
		data, err := c.Get("/v1/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != `{"ok":true}` {
			t.Fatalf("got %s", data)
		}
	})

	t.Run("with query params", func(t *testing.T) {
		c := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Query().Get("email") != "a@b.com" {
				t.Errorf("query = %s", r.URL.RawQuery)
			}
			w.Write([]byte(`{}`))
		})
		q := url.Values{"email": {"a@b.com"}}
		_, err := c.Get("/v1/test", q)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("http error", func(t *testing.T) {
		c := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(401)
			w.Write([]byte(`{"error":"unauthorized"}`))
		})
		_, err := c.Get("/v1/test", nil)
		if err == nil {
			t.Fatal("expected error")
		}
		want := `HTTP 401: {"error":"unauthorized"}`
		if err.Error() != want {
			t.Fatalf("got %q, want %q", err.Error(), want)
		}
	})

	t.Run("http error empty body", func(t *testing.T) {
		c := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(500)
		})
		_, err := c.Get("/v1/test", nil)
		if err == nil {
			t.Fatal("expected error")
		}
	})

	t.Run("empty success body returns {}", func(t *testing.T) {
		c := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(204)
		})
		data, err := c.Get("/v1/test", nil)
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != "{}" {
			t.Fatalf("got %s", data)
		}
	})
}

func TestPost(t *testing.T) {
	t.Run("with raw json body", func(t *testing.T) {
		c := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPost {
				t.Errorf("method = %s", r.Method)
			}
			body, _ := io.ReadAll(r.Body)
			if string(body) != `{"name":"test"}` {
				t.Errorf("body = %s", body)
			}
			w.Write([]byte(`{"id":1}`))
		})
		data, err := c.Post("/v1/test", json.RawMessage(`{"name":"test"}`))
		if err != nil {
			t.Fatal(err)
		}
		if string(data) != `{"id":1}` {
			t.Fatalf("got %s", data)
		}
	})

	t.Run("with nil body", func(t *testing.T) {
		c := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			if len(body) != 0 {
				t.Errorf("expected empty body, got %s", body)
			}
			w.Write([]byte(`{}`))
		})
		_, err := c.Post("/v1/test", nil)
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("with struct body", func(t *testing.T) {
		type req struct {
			Name string `json:"name"`
		}
		c := testServer(t, func(w http.ResponseWriter, r *http.Request) {
			body, _ := io.ReadAll(r.Body)
			if string(body) != `{"name":"hello"}` {
				t.Errorf("body = %s", body)
			}
			w.Write([]byte(`{}`))
		})
		_, err := c.Post("/v1/test", req{Name: "hello"})
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestPut(t *testing.T) {
	c := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s", r.Method)
		}
		w.Write([]byte(`{"updated":true}`))
	})
	data, err := c.Put("/v1/test", json.RawMessage(`{}`))
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != `{"updated":true}` {
		t.Fatalf("got %s", data)
	}
}

func TestDelete(t *testing.T) {
	c := testServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		w.WriteHeader(204)
	})
	data, err := c.Delete("/v1/test", nil)
	if err != nil {
		t.Fatal(err)
	}
	if string(data) != "{}" {
		t.Fatalf("got %s", data)
	}
}

func TestEncodeBody(t *testing.T) {
	t.Run("nil", func(t *testing.T) {
		r, err := encodeBody(nil)
		if err != nil || r != nil {
			t.Fatalf("got %v, %v", r, err)
		}
	})

	t.Run("json.RawMessage", func(t *testing.T) {
		r, err := encodeBody(json.RawMessage(`{"a":1}`))
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(r)
		if string(data) != `{"a":1}` {
			t.Fatalf("got %s", data)
		}
	})

	t.Run("[]byte", func(t *testing.T) {
		r, err := encodeBody([]byte(`hello`))
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(r)
		if string(data) != "hello" {
			t.Fatalf("got %s", data)
		}
	})

	t.Run("struct", func(t *testing.T) {
		type s struct {
			X int `json:"x"`
		}
		r, err := encodeBody(s{X: 42})
		if err != nil {
			t.Fatal(err)
		}
		data, _ := io.ReadAll(r)
		if string(data) != `{"x":42}` {
			t.Fatalf("got %s", data)
		}
	})

	t.Run("unmarshalable", func(t *testing.T) {
		_, err := encodeBody(make(chan int))
		if err == nil {
			t.Fatal("expected error")
		}
	})
}
