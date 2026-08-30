# peep Documentation

peep is a terminal UI for exploring JSON: pipe in a payload, open a file, or paste
from the clipboard, then browse it as a collapsible tree, a Miller-column graph, or
raw pretty-printed text — with English-friendly search and live `jq` filtering built in. It also
runs headless (no TUI) for scripting.

## Why peep exists

Understanding an unfamiliar JSON payload in a terminal usually means picking one of
a few unsatisfying options: `cat`/`less` a file and scroll through an undifferentiated
wall of text with no sense of shape or depth; reach for `jq`, which is precise but
needs you to already know the path you're after; or paste it into a browser-based
JSON viewer, which means leaving the terminal — and often the machine — you're
actually working on (an SSH session, a CI log, a response piped straight from
`curl`).

peep is built for the moment *before* you know what you're looking for: a large or
deeply-nested API response, a config file, a log line, a database export. The goal
is to make the **shape** of the data visible first — via a collapsible tree or a
Miller-column graph — so you can navigate to what matters, then drop into `jq`
filtering or English-friendly search once you know what you're after, all without leaving the
shell or losing your place in the surrounding workflow.

## Design goals

- **Structure before query** — tree and graph views exist so you can see what's
  there before you have to write a `jq` expression against it.
- **A good pipeline citizen** — headless mode (`--query`/`--schema`/`--llm`) makes
  peep composable with the rest of the shell instead of a dead-end viewer; output
  can even be piped back into another `peep` invocation.
- **LLM-prompt-friendly by design** — `--schema` and `--llm` exist specifically to
  turn a large payload into a compact, token-cheap shape or value summary for
  pasting into an LLM prompt, not just for human reading.
- **No setup ceremony** — file argument, piped stdin, system clipboard, or
  interactive paste all work with zero configuration.

For install steps, keybindings, and CLI flags, see the
[README](https://github.com/guisolski/peep/blob/main/README.md). This
`docs/` tree covers the internals: how the pieces fit together and where to make a
given kind of change.

## Getting Started

- [Project Structure](PROJECT_STRUCTURE.md) — package layout and where to find code

## Architecture

- [Architecture](architecture/architecture.md) — Bubble Tea model tree, mode switching, the headless dispatch path, and the JSON node/export pipeline

## Guides

- [Adding a View](guides/adding-a-view.md) — how to add a new `SubModel` (view mode) to the TUI
