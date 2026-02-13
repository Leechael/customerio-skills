package test

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"

	"github.com/cucumber/godog"
)

type bddContext struct {
	mockServer *httptest.Server
	mux        *http.ServeMux
	output     string
	exitCode   int
}

func (b *bddContext) aMockAPIServerIsRunning() error {
	b.mux = http.NewServeMux()
	b.mockServer = httptest.NewServer(b.mux)
	return nil
}

func (b *bddContext) theMockServerRespondsToWith(route string, body *godog.DocString) error {
	parts := strings.SplitN(route, " ", 2)
	if len(parts) != 2 {
		return fmt.Errorf("invalid route format %q, expected \"METHOD /path\"", route)
	}
	method := parts[0]
	path := parts[1]
	responseBody := body.Content

	b.mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != method {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, responseBody)
	})
	return nil
}

func (b *bddContext) commandFor(command string) (*exec.Cmd, error) {
	args := strings.Fields(command)
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	if bin := os.Getenv("CIO_BINARY"); bin != "" {
		return exec.Command(bin, args[1:]...), nil
	}

	goArgs := append([]string{"run", ".."}, args[1:]...)
	return exec.Command("go", goArgs...), nil
}

func (b *bddContext) iRun(command string) error {
	cmd, err := b.commandFor(command)
	if err != nil {
		return err
	}
	out, err := cmd.CombinedOutput()
	b.output = string(out)
	b.exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			b.exitCode = exitErr.ExitCode()
		} else {
			return err
		}
	}
	return nil
}

func (b *bddContext) iRunAgainstTheMockServer(command string) error {
	cmd, err := b.commandFor(command)
	if err != nil {
		return err
	}
	cmd.Env = append(os.Environ(),
		"CIO_BASE_URL="+b.mockServer.URL,
		"CUSTOMERIO_API_TOKEN=test-token",
	)
	out, err := cmd.CombinedOutput()
	b.output = string(out)
	b.exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			b.exitCode = exitErr.ExitCode()
		} else {
			return err
		}
	}
	return nil
}

func (b *bddContext) iRunWithoutAnAPIToken(command string) error {
	cmd, err := b.commandFor(command)
	if err != nil {
		return err
	}
	var env []string
	for _, e := range os.Environ() {
		if !strings.HasPrefix(e, "CUSTOMERIO_API_TOKEN=") {
			env = append(env, e)
		}
	}
	if b.mockServer != nil {
		env = append(env, "CIO_BASE_URL="+b.mockServer.URL)
	}
	cmd.Env = env

	out, err := cmd.CombinedOutput()
	b.output = string(out)
	b.exitCode = 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			b.exitCode = exitErr.ExitCode()
		} else {
			return err
		}
	}
	return nil
}

func (b *bddContext) theExitCodeShouldBe(expected int) error {
	if b.exitCode != expected {
		return fmt.Errorf("expected exit code %d but got %d\nOutput:\n%s", expected, b.exitCode, b.output)
	}
	return nil
}

func (b *bddContext) theOutputShouldContain(expected string) error {
	if !strings.Contains(b.output, expected) {
		return fmt.Errorf("expected output to contain %q but got:\n%s", expected, b.output)
	}
	return nil
}

func (b *bddContext) cleanup() {
	if b.mockServer != nil {
		b.mockServer.Close()
		b.mockServer = nil
	}
}

func InitializeScenario(ctx *godog.ScenarioContext) {
	b := &bddContext{}

	ctx.After(func(ctx context.Context, sc *godog.Scenario, err error) (context.Context, error) {
		b.cleanup()
		return ctx, nil
	})

	ctx.Step(`^a mock API server is running$`, b.aMockAPIServerIsRunning)
	ctx.Step(`^the mock server responds to "([^"]*)" with:$`, b.theMockServerRespondsToWith)
	ctx.Step(`^I run "([^"]*)"$`, b.iRun)
	ctx.Step(`^I run "([^"]*)" against the mock server$`, b.iRunAgainstTheMockServer)
	ctx.Step(`^I run "([^"]*)" without an API token$`, b.iRunWithoutAnAPIToken)
	ctx.Step(`^the exit code should be (\d+)$`, b.theExitCodeShouldBe)
	ctx.Step(`^the output should contain "([^"]*)"$`, b.theOutputShouldContain)
}

func TestFeatures(t *testing.T) {
	suite := godog.TestSuite{
		ScenarioInitializer: InitializeScenario,
		Options: &godog.Options{
			Format:   "pretty",
			Paths:    []string{"features"},
			TestingT: t,
		},
	}

	if suite.Run() != 0 {
		t.Fatal("non-zero status returned, failed to run feature tests")
	}
}
