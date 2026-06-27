# Requirement: `my` CLI — OpenClaw data-dir registry & Podman launcher

## Summary

Build `my`, a Go CLI (`github.com/xhd2015/my`) with an `openclaw` subcommand group:

```
my openclaw add data-dir <path> [--note "..."]
my openclaw list
my openclaw run-in-podman --data-dir <path> [--rebuild] [--container-name NAME] [--port PORT]
```

No auto-discovery. Users explicitly register data directories. `run-in-podman` works
standalone (registration not required) but `add`/`list` provide bookkeeping.

## Approved decisions

| Topic | Decision |
|-------|----------|
| `--data-dir` | Points to the `.openclaw` directory itself (not parent) |
| OpenClaw install | `my` ships a minimal Containerfile; `npm install -g openclaw@latest` inside image |
| Local openclaw repo | **Not required** |
| Registration for run | **Not required** — `--data-dir` works without prior `add` |
| Container name | Default `openclaw-gateway`; override `--container-name` |
| Port | Default `18789`; override `--port` |
| Rebuild | Only when image missing, container spec hash changed, or `--rebuild` |
| macOS Podman | Auto-start podman machine if stopped |
| Gateway token | Read from `openclaw.json` (`gateway.auth.token`); fallback `.env` `OPENCLAW_GATEWAY_TOKEN` |
| Duplicate `add` | Warn on stderr, update note if changed, still succeed (exit 0) |

## Data models & storage

### Registry: `~/.config/my/openclaw.json`

Override in tests via `MY_CONFIG_DIR` env.

```json
{
  "data_dirs": [
    {
      "path": "/abs/path/to/.openclaw",
      "note": "optional note",
      "added_at": "2026-06-27T10:00:00+08:00"
    }
  ],
  "image": {
    "spec_hash": "sha256:...",
    "built_at": "2026-06-27T10:05:00+08:00"
  }
}
```

### Container spec (embedded in `my` binary or `internal/openclaw/container/Containerfile`)

```dockerfile
FROM node:22-bookworm-slim
RUN npm install -g openclaw@latest
USER node
WORKDIR /home/node
ENTRYPOINT ["openclaw"]
```

Image name: `my-openclaw:local`

Spec hash: SHA256 of Containerfile contents. Stored in registry `image.spec_hash`.
Rebuild when: image missing OR stored hash ≠ current hash OR `--rebuild`.

### OpenClaw config (read-only, user's data dir)

- `<data-dir>/openclaw.json` — **required** for `run-in-podman`
- `<data-dir>/.env` — optional; `OPENCLAW_GATEWAY_TOKEN` used if not in json
- `<data-dir>/workspace/` — mounted to `/home/node/.openclaw/workspace`

Token resolution order:
1. `gateway.auth.token` in `openclaw.json`
2. `OPENCLAW_GATEWAY_TOKEN` in `<data-dir>/.env`
3. Error if neither found

## CLI behavior

### `my openclaw add data-dir <path> [--note "..."]`

- Resolve `~` and relative paths to absolute
- Verify path exists and is a directory
- Append to registry (or update note if path already registered)
- On duplicate path: print warning to stderr (`data dir already registered: <path>`)
- Exit 0

### `my openclaw list`

- Print tab-separated table: `path`, `note`, `added_at`
- Header row when non-empty
- Empty registry: print `(no data dirs registered)` to stdout, exit 0

### `my openclaw run-in-podman --data-dir <path> [flags]`

Flow:
1. Ensure podman machine running (macOS only; no-op on Linux in tests)
2. Resolve `--data-dir` to absolute; must exist
3. Require `<data-dir>/openclaw.json` exists
4. Resolve gateway token (json then .env)
5. Rebuild image if needed (see rebuild rules)
6. `podman stop <name>` + `podman rm <name>` (ignore errors if absent)
7. `podman run -d --replace` with:
   - `--name <container-name>` (default `openclaw-gateway`)
   - `-v <data-dir>:/home/node/.openclaw`
   - `-v <data-dir>/workspace:/home/node/.openclaw/workspace`
   - `-e OPENCLAW_GATEWAY_TOKEN=<token>`
   - `-p <port>:18789`
   - `my-openclaw:local gateway --bind lan`
8. Print to stdout: dashboard URL `http://127.0.0.1:<port>/`
9. Print hint: `podman logs -f <container-name>`
10. Exit 0

Flags:
- `--rebuild` — force `podman build`
- `--container-name` — default `openclaw-gateway`
- `--port` — default `18789`

### Top-level `my`

- Dispatcher with subcommands; `my openclaw` delegates to openclaw group
- Use `github.com/xhd2015/less-flags` subcommand pattern

## Podman abstraction (testability)

Inject podman execution via interface or `PATH` stub script in tests.
Record commands issued for assertion.

macOS machine check: `podman machine info`; if not running, `podman machine start`.

## Test tree (doctest)

Location: `my/tests/openclaw/`

```
tests/openclaw/
├── DOCTEST.md
├── SETUP.md
├── registry/
│   ├── add-new-dir/
│   ├── add-with-note/
│   ├── add-duplicate-warns/
│   ├── add-missing-dir/
│   └── list-empty/
│   └── list-populated/
├── run-in-podman/
│   ├── missing-data-dir/
│   ├── missing-openclaw-json/
│   ├── happy-path/
│   ├── token-from-json/
│   ├── token-from-env/
│   ├── custom-container-name/
│   ├── custom-port/
│   ├── rebuild-flag/
│   ├── rebuild-on-spec-hash-change/
│   ├── skip-rebuild-when-image-exists/
│   └── podman-machine-auto-start/
```

## Test scenarios & expected output

| Leaf | Expected |
|------|----------|
| `add-new-dir` | Registry created; path stored absolute; exit 0 |
| `add-with-note` | Note persisted |
| `add-duplicate-warns` | stderr contains warning; note updated; exit 0 |
| `add-missing-dir` | exit 1, stderr mentions path not found |
| `list-empty` | stdout `(no data dirs registered)` |
| `list-populated` | header + rows with path, note, added_at |
| `missing-data-dir` | exit 1 |
| `missing-openclaw-json` | exit 1, stderr mentions openclaw.json |
| `happy-path` | podman run with correct mounts, token env, `--bind lan`; dashboard URL printed |
| `token-from-json` | OPENCLAW_GATEWAY_TOKEN matches json value |
| `token-from-env` | token from .env when json has no token |
| `custom-container-name` | `--name` matches flag |
| `custom-port` | `-p` matches flag |
| `rebuild-flag` | `podman build` invoked |
| `rebuild-on-spec-hash-change` | build when stored hash stale |
| `skip-rebuild-when-image-exists` | no build when image present and hash matches |
| `podman-machine-auto-start` | machine start called when stopped (darwin stub) |

## How to test

- Doc-style doctests under `my/tests/openclaw/`
- Version `0.0.2`, full DSN in root DOCTEST.md
- Stub `podman` script on PATH recording invocations
- Temp config dir via `MY_CONFIG_DIR`
- Fixture data dirs under `testdata/` with minimal `openclaw.json`
- Run: `doctest vet ./tests/openclaw && doctest test ./tests/openclaw`

## Module layout (implementer guidance)

```
cmd/my/main.go
internal/openclaw/
  registry.go
  list.go
  add.go
  run.go
  token.go
  image.go
  podman.go
internal/openclaw/container/Containerfile
```

## Out of scope (v1)

- Session/activity listing from jsonl
- `my` subcommands beyond `openclaw`
- Multiple concurrent instances documentation
- OpenClaw version pinning flag