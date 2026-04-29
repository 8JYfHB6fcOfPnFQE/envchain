# envchain

A tool for chaining and validating environment variable sets across multiple deployment contexts.

---

## Installation

```bash
go install github.com/yourname/envchain@latest
```

Or build from source:

```bash
git clone https://github.com/yourname/envchain.git && cd envchain && go build ./...
```

---

## Usage

Define your environment variable sets in a `.envchain.yaml` file:

```yaml
chains:
  production:
    requires:
      - DATABASE_URL
      - API_KEY
      - SECRET_TOKEN
    inherits: base

  base:
    requires:
      - APP_ENV
      - LOG_LEVEL
```

Then validate and apply a chain before running your application:

```bash
# Validate that all required variables are set
envchain validate --chain production

# Run a command within a validated environment context
envchain run --chain production -- ./myapp serve
```

If any required variables are missing, `envchain` will report them and exit with a non-zero status:

```
✗ Missing variables for chain "production":
  - API_KEY
  - SECRET_TOKEN
```

---

## License

MIT © [yourname](https://github.com/yourname)