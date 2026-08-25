# BUG_REPRO

The following failures were observed while validating the initial project state.
Each section records what failed, how to reproduce it, and the complete command output.
They are preserved intentionally; only failing build gates are omitted from the generated Dockerfile.

## Failure 1: Go test (.)

- Observed problem: `Go test (.)` failed in the initial project state.
- Working directory: `.`
- Command: `cd /app && GOTOOLCHAIN=local GOPROXY=off GOSUMDB=off go test -count=1 ./...`
- Exit status: `1`

```text
ok  	example.com/graduation-showcase/cmd/showcase	0.012s
ok  	example.com/graduation-showcase/internal/api	0.017s
ok  	example.com/graduation-showcase/internal/domain	0.001s
ok  	example.com/graduation-showcase/internal/fixtures	0.001s
--- FAIL: Test2169BusinessRegression (0.01s)
    flow_test.go:56: 刷新后第二张单据标签应为新标签，实际为 旧标签
FAIL
FAIL	example.com/graduation-showcase/internal/flow014	0.015s
ok  	example.com/graduation-showcase/internal/importer	0.008s
ok  	example.com/graduation-showcase/internal/report	0.001s
ok  	example.com/graduation-showcase/internal/store	0.012s
ok  	example.com/graduation-showcase/internal/workflow	0.006s
FAIL
```

## Architecture reproduction

### linux/amd64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/showcase): exit `0`
- Frontend build (web): exit `0`
### linux/arm64
- Go toolchain version: exit `0`
- Node.js version: exit `0`
- Go build (.): exit `0`
- Go test (.): exit `1`
- Go run smoke (cmd/showcase): exit `0`
- Frontend build (web): exit `0`
