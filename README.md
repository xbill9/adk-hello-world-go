# ADK Hello World (Go)

Starter project for building AI agents with the [Go Agent Development Kit (ADK 2.0)](https://github.com/google/adk-go) (`google.golang.org/adk/v2`) and deploying them across multiple environments (local CLI, local web UI, Cloud Run), including an Agent-to-Agent (A2A) multi-agent demo.

**Versions:** Go ADK 2.0 (`google.golang.org/adk/v2` `v2.1.0`) · Google GenAI SDK `v1.65.0` · Go `1.26.5` module directives (built with Go `1.26.3`)

See [Building AI Agents with the GO ADK (v2)](<Building_AI_Agents_with_the_GO_Agent_Development_Kit_(ADK)_v2.md>) for the full tutorial.

## Resources & Links

* **Go ADK 2.0 GA Portal:** [https://adk.dev/2.0/](https://adk.dev/2.0/) (General Availability: June 30, 2026)
* **Official Go ADK Repository:** [https://github.com/google/adk-go](https://github.com/google/adk-go)
* **Go ADK 1.x Compatibility Guide:** [https://adk.dev/2.0/#adk-go-1x-compatibility](https://adk.dev/2.0/#adk-go-1x-compatibility)
* **Go ADK 2.0 Package Documentation:** [pkg.go.dev/google.golang.org/adk/v2](https://pkg.go.dev/google.golang.org/adk/v2)
* **Google GenAI SDK:** [github.com/googleapis/go-genai](https://github.com/googleapis/go-genai)

## Project Layout

Five independent Go modules:

| Module | Description |
|---|---|
| `hello-agent/` | Single agent (`hello_time_agent`) that tells the time in a city using Google Search. Main deployment target. |
| `a2a-server-go/` | A2A server exposing a `check_prime_agent` with a prime-checking function tool (default port `8086`). |
| `a2a-master-go/` | Root agent served over A2A (default port `8092`) that delegates to a local `roll_agent` and the remote `prime_agent`. |
| `a2a-client-go/` | CLI client running the same root/roll/prime agent tree; takes the user prompt from command-line args. |
| `a2a-gemini3-go/` | Variant of hello-agent on `gemini-3-pro-preview` with web UI + A2A and automatic port fallback (default port `8095`). |

`cloudrun.go` is a reference `adkgo deploy cloudrun` subcommand (built against ADK internals; not part of the local modules).

## Setup

```bash
./init.sh           # One-time: saves project ID (~/project_id.txt) and optional Gemini API key (~/gemini.key)
source set_env.sh   # Per-shell: exports PROJECT_ID, GOOGLE_CLOUD_PROJECT, REGION, etc.
```

Authentication: set `GOOGLE_API_KEY` to use the Gemini API directly, otherwise the agents fall back to Vertex AI with Application Default Credentials (`gcloud auth application-default login`).

## Running Locally

```bash
./cli.sh                  # hello-agent in interactive CLI mode
./web.sh                  # hello-agent with web UI + REST API (http://localhost:8080)

# A2A multi-agent demo (two terminals):
./run-a2a-server-go.sh    # 1. start the prime-checker A2A server on :8086
./run-a2a-master-go.sh    # 2. start the root agent A2A server on :8092
# or run the one-shot client instead of the master:
cd a2a-client-go && go run main.go "Roll a 6-sided die and check if it's prime."
```

## Make Targets

The root `Makefile` operates on all five modules:

```bash
make build    # go build each module
make test     # go test -v each module
make lint     # go vet each module
make format   # go fmt each module
make deps     # go get -u && go mod tidy each module
```

## Deployment (Cloud Run)

```bash
./cloudrun.sh        # adkgo deploy cloudrun (compiles locally, generates Dockerfile, deploys, starts auth proxy)
./quickrun.sh        # gcloud run deploy --source . using the root Dockerfile
make deploy          # gcloud builds submit with cloudbuild.yaml
./proxy.sh           # gcloud run services proxy -> http://127.0.0.1:8081/ui/
```

The root `Dockerfile` is a multi-stage build of `hello-agent` onto a distroless base image listening on port `8080`.

## Environment Variables

| Variable | Default | Used by |
|---|---|---|
| `GOOGLE_API_KEY` | – | Gemini API auth (falls back to Vertex AI if unset) |
| `GOOGLE_CLOUD_PROJECT` / `PROJECT_ID` | – | Vertex AI project |
| `GOOGLE_CLOUD_LOCATION` / `REGION` | `us-central1` | Vertex AI location |
| `MODEL_NAME` | `gemini-2.5-flash` (`gemini-3-pro-preview` in a2a-gemini3-go) | hello-agent, a2a-gemini3-go |
| `PORT` | `8086` / `8092` / `8095` | A2A servers |
| `ADK_PRIME_AGENT_URL` | `http://localhost:8086` | a2a-client-go, a2a-master-go |
| `ADK_MODEL_NAME` | `gemini-2.5-flash` | a2a-client-go |
