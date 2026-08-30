# peep

A terminal UI for exploring JSON. Pipe in a payload, open a file, or paste from
the clipboard, then browse it as a collapsible tree, a Miller-column graph, or
raw pretty-printed text — with English-friendly search and live `jq` filtering built in.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and
[gojq](https://github.com/itchyny/gojq).

## Features

- **Tree view** — collapsible JSON tree with `├──`/`└──` connector lines
- **Graph view** — Miller-column (depth-panel) layout for navigating wide/deep structures
- **Order- and precision-preserving** — object keys keep the source document's original order (never re-sorted), and large integers (e.g. IDs beyond 2^53) keep full precision wherever a value is displayed, copied, or exported
- **Raw view** — scrollable, pretty-printed JSON
- **Search** — English-friendly search over keys/values (`id:2`, `which have id 2?`, `value greater than 5`) with `n`/`N` to cycle hits
- **Filter** — live `jq` expression evaluation against the loaded document
- **Clipboard integration** — copy a value, a subtree, a JSON path, a schema, or an LLM-friendly export with a keystroke
- **LLM export** — copy or print a compact schema (types + example values) or a dense YAML-like serialization, sized for pasting into a prompt
- **Headless mode** — run a `jq` query, a schema dump, or an LLM export from a script, with no TUI
- **Flexible input** — read from a file argument, piped stdin, the system clipboard, or interactive paste

## Installation

### Homebrew

```sh
brew tap guisolski/peep https://github.com/guisolski/peep
brew trust guisolski/peep   # recent Homebrew versions require trusting non-core taps
brew install peep
```

### Go

Requires Go 1.27+.

```sh
go install github.com/guisolski/peep@latest   # installs to $(go env GOPATH)/bin
```

Or from a clone:

```sh
git clone git@github.com:guisolski/peep.git
cd peep
make install   # builds and installs to $(PREFIX)/bin
```

The default prefix is `/usr/local`, falling back to `/opt/homebrew` when
`/usr/local/bin` doesn't exist (the usual case on Apple Silicon macOS). If the
chosen prefix isn't writable by your user, run `sudo make install`. Any other
location works too: `make install PREFIX=/custom/path`.

Or just build the binary locally:

```sh
make build      # produces ./peep
```

## Usage

```sh
peep file.json          # open a file
cat file.json | peep    # read piped stdin
peep --paste             # read JSON from the system clipboard
peep                      # no args, no pipe: prompts for interactive paste (Ctrl+D to submit)
```

### Headless flags

Skip the TUI entirely and print to stdout — useful in scripts or when piping into another tool (including another `peep`):

```sh
peep --query '.items[]' file.json   # run a jq expression, print each result
peep --schema file.json              # print a compact schema (types + example values)
peep --llm file.json                 # print a dense, YAML-like export sized for an LLM prompt
peep --version                       # print the version and exit
```

`--query`, `--schema`, and `--llm` are mutually exclusive. Combine them by piping: `peep --query '.items[]' file.json | peep --schema` gives you the schema of a filtered subset.

## Keybindings

### Global

| Key | Action |
|---|---|
| `j` / `k` or `↓` / `↑` | Move down / up |
| `l` / `→` or `enter` | Expand node / drill in |
| `h` / `←` | Collapse node / go back |
| `gg` | Jump to top |
| `g` | Toggle graph view |
| `r` | Toggle raw view |
| `/` | Search |
| `:` | Filter (jq) |
| `y` (double-tap) | Copy current subtree as JSON |
| `Y` | Copy current node's JSON path |
| `S` | Copy current subtree's schema (types + example values) |
| `L` | Copy current subtree as a compact LLM-friendly export |
| `q` / `ctrl+c` | Quit |

### Tree view

| Key | Action |
|---|---|
| `ctrl+d` / `ctrl+u` | Half-page down / up |

### Search (`/`)

Type plain English or simple patterns. Hits are highlighted in the tree; collapsed
ancestors expand automatically.

| Example | Meaning |
|---|---|
| `alice` | key or value contains `alice` |
| `id 2` | both terms must match the same node |
| `id:2` / `id: 2` | field `id` whose value contains `2` |
| `which have id 2?` | same idea after stripping English filler words |
| `value greater than 5` / `id > 5` | numeric comparison on a field |
| `"via place"` | literal phrase |

| Key | Action |
|---|---|
| `n` / `N` | Next / previous match |
| `esc` | Exit back to tree |

### Filter (`:`)

Type a `jq` expression and see the result rendered live in the raw view. Use this
for complex queries; `/` is for quick human search.

| Key | Action |
|---|---|
| `esc` | Exit back to tree |

## Project layout

```
main.go        entry point: flag parsing, input loading, bubbletea startup
cli.go         headless dispatch for --query/--schema/--llm
input/         source detection (file / stdin / clipboard / interactive paste)
model/         application state: App root model, Tree/Graph/Raw/Search/Filter sub-models, JSON node parsing, schema/LLM export, jq evaluation
clipboard/     clipboard write helper
ui/            terminal-adaptive lipgloss color styles
testdata/      sample JSON fixtures used by tests
Formula/       Homebrew formula, auto-generated/updated by the release workflow — do not edit by hand
```

See [docs/](docs/index.md) for architecture and contributor guides.

## Development

```sh
make build     # go build -o peep .
make vet       # go vet ./...
make lint      # golangci-lint run ./...
make test      # go test ./... -race
make clean     # remove the built binary
```

## License

MIT — see [LICENSE](LICENSE).
