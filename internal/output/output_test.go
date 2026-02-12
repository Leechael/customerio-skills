package output

import (
	"bytes"
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func captureStdout(t *testing.T, fn func()) string {
	t.Helper()
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	fn()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	buf.ReadFrom(r)
	return buf.String()
}

func TestPrintPretty(t *testing.T) {
	t.Run("pretty prints json", func(t *testing.T) {
		data := json.RawMessage(`{"name":"test","count":1}`)
		out := captureStdout(t, func() {
			if err := Print(data, ""); err != nil {
				t.Fatal(err)
			}
		})
		if !strings.Contains(out, "  \"name\": \"test\"") {
			t.Fatalf("not indented:\n%s", out)
		}
		if !strings.HasSuffix(out, "}\n") {
			t.Fatalf("no trailing newline:\n%q", out)
		}
	})

	t.Run("empty object", func(t *testing.T) {
		out := captureStdout(t, func() {
			Print(json.RawMessage(`{}`), "")
		})
		if strings.TrimSpace(out) != "{}" {
			t.Fatalf("got %q", out)
		}
	})

	t.Run("array", func(t *testing.T) {
		out := captureStdout(t, func() {
			Print(json.RawMessage(`[1,2,3]`), "")
		})
		if !strings.Contains(out, "[\n") {
			t.Fatalf("not indented:\n%s", out)
		}
	})
}

func TestPrintJQ(t *testing.T) {
	t.Run("simple field access", func(t *testing.T) {
		data := json.RawMessage(`{"name":"hello","count":5}`)
		out := captureStdout(t, func() {
			if err := Print(data, ".name"); err != nil {
				t.Fatal(err)
			}
		})
		if strings.TrimSpace(out) != "hello" {
			t.Fatalf("got %q", out)
		}
	})

	t.Run("numeric field", func(t *testing.T) {
		data := json.RawMessage(`{"count":42}`)
		out := captureStdout(t, func() {
			if err := Print(data, ".count"); err != nil {
				t.Fatal(err)
			}
		})
		if strings.TrimSpace(out) != "42" {
			t.Fatalf("got %q", out)
		}
	})

	t.Run("array index", func(t *testing.T) {
		data := json.RawMessage(`{"items":["a","b","c"]}`)
		out := captureStdout(t, func() {
			if err := Print(data, ".items[0]"); err != nil {
				t.Fatal(err)
			}
		})
		if strings.TrimSpace(out) != "a" {
			t.Fatalf("got %q", out)
		}
	})

	t.Run("nested access", func(t *testing.T) {
		data := json.RawMessage(`{"segments":[{"name":"vip"}]}`)
		out := captureStdout(t, func() {
			if err := Print(data, ".segments[0].name"); err != nil {
				t.Fatal(err)
			}
		})
		if strings.TrimSpace(out) != "vip" {
			t.Fatalf("got %q", out)
		}
	})

	t.Run("object result", func(t *testing.T) {
		data := json.RawMessage(`{"a":{"b":1}}`)
		out := captureStdout(t, func() {
			if err := Print(data, ".a"); err != nil {
				t.Fatal(err)
			}
		})
		if !strings.Contains(out, `"b": 1`) {
			t.Fatalf("got %q", out)
		}
	})

	t.Run("multiple results", func(t *testing.T) {
		data := json.RawMessage(`{"items":["x","y"]}`)
		out := captureStdout(t, func() {
			if err := Print(data, ".items[]"); err != nil {
				t.Fatal(err)
			}
		})
		lines := strings.Split(strings.TrimSpace(out), "\n")
		if len(lines) != 2 || lines[0] != "x" || lines[1] != "y" {
			t.Fatalf("got %q", out)
		}
	})

	t.Run("invalid jq expression", func(t *testing.T) {
		err := Print(json.RawMessage(`{}`), "..invalid[")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "jq parse error") {
			t.Fatalf("got %q", err.Error())
		}
	})

	t.Run("invalid json input", func(t *testing.T) {
		err := Print(json.RawMessage(`not json`), ".foo")
		if err == nil {
			t.Fatal("expected error")
		}
		if !strings.Contains(err.Error(), "json unmarshal error") {
			t.Fatalf("got %q", err.Error())
		}
	})

	t.Run("null result", func(t *testing.T) {
		data := json.RawMessage(`{"a":1}`)
		out := captureStdout(t, func() {
			if err := Print(data, ".b"); err != nil {
				t.Fatal(err)
			}
		})
		if strings.TrimSpace(out) != "null" {
			t.Fatalf("got %q", out)
		}
	})

	t.Run("pipe and select", func(t *testing.T) {
		data := json.RawMessage(`{"items":[{"n":1},{"n":2},{"n":3}]}`)
		out := captureStdout(t, func() {
			if err := Print(data, `[.items[] | select(.n > 1)]`); err != nil {
				t.Fatal(err)
			}
		})
		if !strings.Contains(out, `"n": 2`) || !strings.Contains(out, `"n": 3`) {
			t.Fatalf("got %q", out)
		}
	})
}
