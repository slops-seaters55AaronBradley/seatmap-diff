# seatmap-diff

> CLI tool for diffing and auditing changes in YAML/JSON infrastructure configs across environments

---

## Installation

```bash
go install github.com/yourorg/seatmap-diff@latest
```

Or build from source:

```bash
git clone https://github.com/yourorg/seatmap-diff.git
cd seatmap-diff
go build -o seatmap-diff .
```

---

## Usage

Compare two infrastructure config files across environments:

```bash
seatmap-diff --base configs/staging.yaml --target configs/production.yaml
```

Output a structured audit report in JSON:

```bash
seatmap-diff --base staging.json --target production.json --format json --output report.json
```

Ignore specific keys during comparison:

```bash
seatmap-diff --base staging.yaml --target production.yaml --ignore metadata.timestamp,metadata.version
```

### Flags

| Flag | Description | Default |
|------|-------------|---------|
| `--base` | Path to the base config file | required |
| `--target` | Path to the target config file | required |
| `--format` | Output format: `text`, `json`, `yaml` | `text` |
| `--output` | Write output to a file instead of stdout | — |
| `--ignore` | Comma-separated list of keys to ignore | — |
| `--strict` | Exit with non-zero code if differences are found | `false` |

---

## Example Output

```
[CHANGED]  server.replicas          3 → 5
[ADDED]    server.autoscaling       true
[REMOVED]  database.legacy_mode
```

---

## License

MIT © yourorg
