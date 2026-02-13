package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"

	"github.com/leechael/cio/internal/client"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// setupTestServer creates a test HTTP server and patches newClient to use it.
// Returns a cleanup function.
func setupTestServer(t *testing.T, handler http.HandlerFunc) func() {
	t.Helper()
	srv := httptest.NewServer(handler)

	orig := newClient
	newClient = func() (*client.Client, error) {
		return &client.Client{
			BaseURL:    srv.URL,
			Token:      "test-token",
			HTTPClient: srv.Client(),
		}, nil
	}

	return func() {
		newClient = orig
		srv.Close()
	}
}

func resetFlags(cmd *cobra.Command) {
	cmd.Flags().VisitAll(func(f *pflag.Flag) {
		_ = f.Value.Set(f.DefValue)
		f.Changed = false
	})
	for _, sub := range cmd.Commands() {
		resetFlags(sub)
	}
}

func executeCommand(args ...string) (string, error) {
	jqExpr = ""
	region = "us"
	jsonOutput = false
	plainOutput = false
	resetFlags(rootCmd)

	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	rootCmd.SetArgs(args)
	err := rootCmd.Execute()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	_, _ = buf.ReadFrom(r)
	return buf.String(), err
}

func TestSegmentsList(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/segments" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.Method != http.MethodGet {
			t.Errorf("method = %s", r.Method)
		}
		_, _ = w.Write([]byte(`{"segments":[{"id":1,"name":"VIP"}]}`))
	})
	defer cleanup()

	out, err := executeCommand("segments", "ls")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"name": "VIP"`) {
		t.Fatalf("got %s", out)
	}
}

func TestSegmentsListAlias(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"segments":[]}`))
	})
	defer cleanup()

	out, err := executeCommand("segments", "list")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "segments") {
		t.Fatalf("got %s", out)
	}
}

func TestSegmentsGet(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/segments/42" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"segment":{"id":42}}`))
	})
	defer cleanup()

	out, err := executeCommand("segments", "get", "42")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, `"id": 42`) {
		t.Fatalf("got %s", out)
	}
}

func TestSegmentsGetMissingArg(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	})
	defer cleanup()

	_, err := executeCommand("segments", "get")
	if err == nil {
		t.Fatal("expected error for missing arg")
	}
}

func TestSegmentsDelete(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/segments/99" {
			t.Errorf("path = %s", r.URL.Path)
		}
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		w.WriteHeader(204)
	})
	defer cleanup()

	_, err := executeCommand("segments", "rm", "99")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSegmentsDeleteAlias(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			t.Errorf("method = %s", r.Method)
		}
		w.WriteHeader(204)
	})
	defer cleanup()

	_, err := executeCommand("segments", "delete", "1")
	if err != nil {
		t.Fatal(err)
	}
}

func TestSegmentsCount(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/segments/5/customer_count" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"count":123}`))
	})
	defer cleanup()

	out, err := executeCommand("segments", "count", "5")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "123") {
		t.Fatalf("got %s", out)
	}
}

func TestCollectionsContentGetPutMerge(t *testing.T) {
	t.Run("GET when no body", func(t *testing.T) {
		cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/v1/collections/10/content" {
				t.Errorf("path = %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"data":"rows"}`))
		})
		defer cleanup()

		out, err := executeCommand("collections", "content", "10")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "rows") {
			t.Fatalf("got %s", out)
		}
	})

	t.Run("PUT when body provided", func(t *testing.T) {
		cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"ok":true}`))
		})
		defer cleanup()

		out, err := executeCommand("collections", "content", "10", "--body", `{"rows":[]}`)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, "true") {
			t.Fatalf("got %s", out)
		}
	})
}

func TestCampaignsTranslationMerge(t *testing.T) {
	t.Run("GET translation", func(t *testing.T) {
		cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/v1/campaigns/1/actions/2/language/en" {
				t.Errorf("path = %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"lang":"en"}`))
		})
		defer cleanup()

		out, err := executeCommand("campaigns", "translation", "1", "2", "en")
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(out, `"lang": "en"`) {
			t.Fatalf("got %s", out)
		}
	})

	t.Run("PUT translation", func(t *testing.T) {
		cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", r.Method)
			}
			_, _ = w.Write([]byte(`{}`))
		})
		defer cleanup()

		_, err := executeCommand("campaigns", "translation", "1", "2", "en", "--body", `{"subject":"hi"}`)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestCustomersSearchByEmail(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			t.Errorf("expected GET, got %s", r.Method)
		}
		if r.URL.Query().Get("email") != "user@test.com" {
			t.Errorf("email = %s", r.URL.Query().Get("email"))
		}
		_, _ = w.Write([]byte(`{"results":[]}`))
	})
	defer cleanup()

	_, err := executeCommand("customers", "search", "--email", "user@test.com")
	if err != nil {
		t.Fatal(err)
	}
}

func TestCustomersSearchByBody(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("expected POST, got %s", r.Method)
		}
		if r.URL.Path != "/v1/customers" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"results":[]}`))
	})
	defer cleanup()

	_, err := executeCommand("customers", "search", "--body", `{"filter":{}}`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestSendEmail(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/send/email" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"delivery_id":"abc"}`))
	})
	defer cleanup()

	out, err := executeCommand("send", "email", "--body", `{"to":"a@b.com"}`)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "abc") {
		t.Fatalf("got %s", out)
	}
}

func TestJQFilter(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"segments":[{"name":"VIP"},{"name":"Free"}]}`))
	})
	defer cleanup()

	out, err := executeCommand("segments", "ls", "--jq", ".segments[0].name")
	if err != nil {
		t.Fatal(err)
	}
	if strings.TrimSpace(out) != "VIP" {
		t.Fatalf("got %q", out)
	}
}

func TestJQRequiresJSONMode(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{"segments":[]}`))
	})
	defer cleanup()

	_, err := executeCommand("segments", "ls", "--plain", "--jq", ".segments[0]")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "--jq requires JSON output mode") {
		t.Fatalf("got %q", err.Error())
	}
}

func TestHTTPError(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(403)
		_, _ = w.Write([]byte(`{"error":"forbidden"}`))
	})
	defer cleanup()

	_, err := executeCommand("segments", "ls")
	if err == nil {
		t.Fatal("expected error")
	}
	if !strings.Contains(err.Error(), "403") {
		t.Fatalf("got %q", err.Error())
	}
}

func TestBroadcastsTrigger(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/campaigns/7/triggers" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"id":"trig-1"}`))
	})
	defer cleanup()

	_, err := executeCommand("broadcasts", "trigger", "7", "--body", `{}`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestEspSuppressionSuppress(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/esp_suppression/bad@test.com" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{}`))
	})
	defer cleanup()

	_, err := executeCommand("esp-suppression", "suppress", "bad@test.com", "--body", `{"reason":"bounce"}`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestWebhooksGetPutMerge(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			_, _ = w.Write([]byte(`{"webhook":{}}`))
		})
		defer cleanup()

		_, err := executeCommand("webhooks", "get", "3")
		if err != nil {
			t.Fatal(err)
		}
	})

	t.Run("PUT", func(t *testing.T) {
		cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodPut {
				t.Errorf("expected PUT, got %s", r.Method)
			}
			_, _ = w.Write([]byte(`{}`))
		})
		defer cleanup()

		_, err := executeCommand("webhooks", "get", "3", "--body", `{"url":"https://example.com"}`)
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestSnippetsUpsert(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPut {
			t.Errorf("method = %s", r.Method)
		}
		if r.URL.Path != "/v1/snippets" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{}`))
	})
	defer cleanup()

	_, err := executeCommand("snippets", "upsert", "--body", `{"snippets":{"footer":"<p>hi</p>"}}`)
	if err != nil {
		t.Fatal(err)
	}
}

func TestObjectsGet(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/objects/companies/acme/attributes" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"object":{}}`))
	})
	defer cleanup()

	_, err := executeCommand("objects", "get", "companies", "acme")
	if err != nil {
		t.Fatal(err)
	}
}

func TestExportsDownload(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/exports/55/download" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"url":"https://dl.example.com/file.csv"}`))
	})
	defer cleanup()

	out, err := executeCommand("exports", "download", "55")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "dl.example.com") {
		t.Fatalf("got %s", out)
	}
}

func TestNewsletterTranslationMerge(t *testing.T) {
	t.Run("GET", func(t *testing.T) {
		cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
			if r.Method != http.MethodGet {
				t.Errorf("expected GET, got %s", r.Method)
			}
			if r.URL.Path != "/v1/newsletters/8/language/fr" {
				t.Errorf("path = %s", r.URL.Path)
			}
			_, _ = w.Write([]byte(`{"lang":"fr"}`))
		})
		defer cleanup()

		_, err := executeCommand("newsletters", "translation", "8", "fr")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestTransactionalDeliveries(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/transactional/12/deliveries" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"deliveries":[]}`))
	})
	defer cleanup()

	_, err := executeCommand("transactional", "deliveries", "12")
	if err != nil {
		t.Fatal(err)
	}
}

func TestInfoIPAddresses(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/info/ip_addresses" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"addresses":["1.2.3.4"]}`))
	})
	defer cleanup()

	out, err := executeCommand("info", "ip-addresses")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "1.2.3.4") {
		t.Fatalf("got %s", out)
	}
}

func TestIndexAttributes(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/v1/index/attributes" {
			t.Errorf("path = %s", r.URL.Path)
		}
		_, _ = w.Write([]byte(`{"attributes":[]}`))
	})
	defer cleanup()

	_, err := executeCommand("index", "attributes")
	if err != nil {
		t.Fatal(err)
	}
}

func TestRequireBodyMissing(t *testing.T) {
	cleanup := setupTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte(`{}`))
	})
	defer cleanup()

	_, err := executeCommand("segments", "create")
	if err == nil {
		t.Fatal("expected error for missing body")
	}
	if !strings.Contains(err.Error(), "body is required") {
		t.Fatalf("got %q", err.Error())
	}
}

// Ensure all top-level groups exist
func TestAllGroupsRegistered(t *testing.T) {
	groups := []string{
		"activities", "broadcasts", "campaigns", "collections",
		"customers", "esp-suppression", "exports", "imports",
		"index", "info", "messages", "newsletters", "objects",
		"segments", "send", "sender-identities", "snippets",
		"subscription-topics", "transactional", "webhooks", "workspaces",
	}

	cmds := make(map[string]bool)
	for _, c := range rootCmd.Commands() {
		cmds[c.Name()] = true
	}

	for _, g := range groups {
		if !cmds[g] {
			t.Errorf("missing command group: %s", g)
		}
	}
}

// Verify subcommand counts for key groups
func TestSubcommandCounts(t *testing.T) {
	cases := []struct {
		group string
		min   int
	}{
		{"customers", 8},
		{"segments", 7},
		{"campaigns", 11},
		{"broadcasts", 14},
		{"newsletters", 13},
		{"transactional", 7},
		{"collections", 6},
		{"messages", 3},
		{"send", 3},
		{"exports", 5},
		{"webhooks", 4},
		{"esp-suppression", 4},
		{"sender-identities", 3},
		{"snippets", 3},
		{"objects", 4},
		{"imports", 2},
		{"index", 2},
	}

	for _, tc := range cases {
		t.Run(tc.group, func(t *testing.T) {
			var target *cobra.Command
			for _, c := range rootCmd.Commands() {
				if c.Name() == tc.group {
					target = c
					break
				}
			}
			if target == nil {
				t.Fatalf("group %s not found", tc.group)
			}
			got := len(target.Commands())
			if got < tc.min {
				names := make([]string, 0, len(target.Commands()))
				for _, sub := range target.Commands() {
					names = append(names, sub.Name())
				}
				t.Errorf("%s: got %d subcommands (%v), want >= %d", tc.group, got, names, tc.min)
			}
		})
	}
}
