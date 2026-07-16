# peep

A terminal UI for exploring JSON. Pipe in a payload, open a file, or paste from
the clipboard, then browse it as a collapsible tree, a Miller-column graph, or
raw pretty-printed text — with fuzzy search and live `jq` filtering built in.

Built with [Bubble Tea](https://github.com/charmbracelet/bubbletea) and
[gojq](https://github.com/itchyny/gojq).

## Features

- **Tree view** — collapsible JSON tree with `├──`/`└──` connector lines
- **Graph view** — Miller-column (depth-panel) layout for navigating wide/deep structures
- **Raw view** — scrollable, pretty-printed JSON
- **Search** — fuzzy-match keys/values and cycle through matches
- **Filter** — live `jq` expression evaluation against the loaded document
- **Clipboard integration** — copy a value, a subtree, or a JSON path with a keystroke
- **Flexible input** — read from a file argument, piped stdin, the system clipboard, or interactive paste

## Installation

Requires Go 1.26+.

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
| `q` / `ctrl+c` | Quit |

### Tree view

| Key | Action |
|---|---|
| `ctrl+d` / `ctrl+u` | Half-page down / up |

### Search (`/`)

| Key | Action |
|---|---|
| `n` / `N` | Next / previous match |
| `esc` | Exit back to tree |

### Filter (`:`)

Type a `jq` expression and see the result rendered live in the raw view.

| Key | Action |
|---|---|
| `esc` | Exit back to tree |

## Project layout

```
main.go        entry point: flag parsing, input loading, bubbletea startup
input/         source detection (file / stdin / clipboard / interactive paste)
model/         application state: App root model, Tree/Graph/Raw/Search/Filter sub-models, JSON node parsing
clipboard/     clipboard read/write helpers
ui/            terminal-adaptive lipgloss color styles
testdata/      sample JSON fixtures used by tests
```

## Development

```sh
make build     # go build -o peep .
make test      # go test ./...
make clean     # remove the built binary
```

## License

MIT — see [LICENSE](LICENSE).
