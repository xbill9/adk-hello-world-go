# Gemini Workspace for `adk-hello-world-go`

You are a Go developer working with Google Cloud and the Go Agent Development Kit (ADK).
This document describes the project layout, common commands, and the best practices to follow when modifying this repository.

## 1. Project Overview

Starter project for building AI agents with the Go ADK (`google.golang.org/adk` **v1.5.1**, `google.golang.org/genai` **v1.65.0**) and deploying to Cloud Run. Toolchain: Go 1.26.x with `go 1.25.0` module directives.

Five independent Go modules (each with its own `go.mod`, tests, and docs):

- `hello-agent/` — `hello_time_agent`: tells the time in a city via Google Search; main deployment target (`agent.go`).
- `a2a-server-go/` — A2A server exposing `check_prime_agent` with a `functiontool` prime checker (port 8086).
- `a2a-master-go/` — root agent served over A2A (port 8092) delegating to a local `roll_agent` and the remote `prime_agent`.
- `a2a-client-go/` — CLI runner for the same agent tree; prompt comes from `os.Args`.
- `a2a-gemini3-go/` — hello-agent variant on `gemini-3-pro-preview` with port-fallback logic (port 8095).

Root-level files: `cloudrun.go` (reference `adkgo deploy cloudrun` subcommand, built against ADK internals — not buildable standalone here), `Dockerfile` (multi-stage distroless build of hello-agent), `cloudbuild.yaml`, `Makefile`, and setup/run shell scripts (`init.sh`, `set_env.sh`, `cli.sh`, `web.sh`, `run-a2a-*.sh`, `cloudrun.sh`).

## 2. Common Commands

```bash
make build     # go build all five modules
make test      # go test -v all five modules — run after every change
make lint      # go vet all five modules
make format    # go fmt all five modules
make deps      # go get -u && go mod tidy per module
make deploy    # gcloud builds submit (cloudbuild.yaml)
```

Local run: `./cli.sh` (CLI), `./web.sh` (web UI on :8080). A2A demo: `./run-a2a-server-go.sh` then `./run-a2a-master-go.sh`.

## 3. Established Patterns in This Repo

Follow these existing conventions when adding or modifying code:

- **Structured logging with `log/slog`:** every entrypoint installs `slog.New(slog.NewJSONHandler(...))` as the default logger for Cloud Logging compatibility. Do not introduce `zap`, `logrus`, or bare `fmt.Println` logging.
- **Graceful shutdown:** entrypoints derive their context from `signal.NotifyContext(context.Background(), os.Interrupt)` and pass it through to models, runners, and launchers.
- **Auth fallback:** models prefer `GOOGLE_API_KEY` (Gemini API) and fall back to Vertex AI (`genai.BackendVertexAI`) with `GOOGLE_CLOUD_PROJECT`/`GOOGLE_CLOUD_LOCATION` and Application Default Credentials.
- **Tools via `functiontool.New`:** typed args structs with `json` + `jsonschema` tags; validate inputs and return errors rather than panicking (see `rollDieTool`'s `sides <= 0` check).
- **Agent loaders:** single-agent programs adapt their agent to the ADK `AgentLoader` interface with a small `SingleAgentLoader` type; `LoadAgent` should return an error (not `nil, nil`) for unknown names.
- **Launchers:** `full.NewLauncher()` for the dev experience (web/api/webui subcommands), `prod.NewLauncher()` for A2A servers. Port comes from `PORT` env var with a per-module default.

## 4. General Go Best Practices

### 4.1 Error Handling
- Always check errors explicitly; wrap with `fmt.Errorf("...: %w", err)` to preserve the chain.
- Add enough context to error messages to aid debugging (what was attempted, with which inputs).

### 4.2 Concurrency
- Manage goroutines with `sync.WaitGroup` or context cancellation — no naked goroutines.
- Use `context.Context` for cancellation and timeouts across call hierarchies.

### 4.3 Structure and Naming
- One clear responsibility per package; keep functions small.
- Exported identifiers use CamelCase and carry Godoc comments; short names (`i`, `err`, `tc`) are fine in small scopes.

### 4.4 Testing
- Use table-driven tests (see `a2a-server-go/main_test.go` for the house style).
- Test pure logic (tool functions, helpers) without requiring network or credentials.
- Run `make test` and `make lint` before committing; all must pass.

## 5. Google Cloud Deployment

- **Containerization:** multi-stage Docker builds ending on `gcr.io/distroless/static-debian11`; build with `CGO_ENABLED=0 GOOS=linux`.
- **Cloud Run:** deploy with `--no-allow-unauthenticated` and access via `gcloud run services proxy` (see `proxy.sh`); right-size instances and rely on autoscaling.
- **Secrets:** use Secret Manager for API keys and credentials; never commit keys or bake them into images. Grant service accounts least-privilege IAM roles.
- **Observability:** JSON logs flow to Cloud Logging automatically; use OpenTelemetry (already an ADK dependency) for tracing and custom metrics.
- **Builds:** automate with Cloud Build (`cloudbuild.yaml`); keep the substitution `_SERVICE_NAME` in sync with the Makefile's `SERVICE_NAME`.
