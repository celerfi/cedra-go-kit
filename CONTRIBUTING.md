# Contributing to cedra-go-kit

## Getting Started

1. Fork the repository and clone your fork
2. Create a feature branch: `git checkout -b feat/your-feature`
3. Make your changes
4. Run tests and vet: `go test ./... && go vet ./...`
5. Push and open a pull request against `master`

## Guidelines

- Every new function must have a corresponding test
- Run `gofmt -w .` before committing
- Keep PRs focused — one feature or fix per PR
- Do not break existing tests
- Match the existing code style (no docstrings, no unnecessary comments)

## Running Tests

```bash
go test ./... -v -race
```

## Reporting Bugs

Open an issue with:
- Go version (`go version`)
- OS and architecture
- Minimal reproduction case
- Expected vs actual behavior

## Security Issues

Do not open public issues for security vulnerabilities. See [SECURITY.md](SECURITY.md).
